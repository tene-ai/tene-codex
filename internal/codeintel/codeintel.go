// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package codeintel

import (
	"bufio"
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const MaxFileSize = 2 << 20

type Capability struct {
	Provider   string   `json:"provider"`
	Available  bool     `json:"available"`
	Languages  []string `json:"languages,omitempty"`
	Calls      bool     `json:"calls"`
	Imports    bool     `json:"imports"`
	DataFlow   bool     `json:"data_flow"`
	Runtime    bool     `json:"runtime"`
	Confidence float64  `json:"confidence"`
	Reason     string   `json:"reason,omitempty"`
}

type Component struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Kind            string   `json:"kind"`
	Locator         string   `json:"locator"`
	File            string   `json:"file"`
	PrimaryLayer    string   `json:"primary_layer"`
	SecondaryLayers []string `json:"secondary_layers,omitempty"`
	LayerReason     string   `json:"layer_reason"`
	Imports         []string `json:"imports"`
	References      []string `json:"references"`
	Calls           []string `json:"calls"`
	Inputs          []string `json:"inputs"`
	Outputs         []string `json:"outputs"`
	Effects         []string `json:"effects"`
	Unknown         []string `json:"unknown"`
	Provider        string   `json:"provider"`
	Confidence      float64  `json:"confidence"`
}

type File struct {
	ID          string  `json:"id"`
	Path        string  `json:"path"`
	Layer       string  `json:"layer"`
	LayerReason string  `json:"layer_reason"`
	Confidence  float64 `json:"confidence"`
}

type Edge struct {
	From, To, Kind, Locator, Provider string
	Confidence                        float64
}
type Report struct {
	Providers       []Capability `json:"providers"`
	Files           []File       `json:"files"`
	Components      []Component  `json:"components"`
	Edges           []Edge       `json:"edges"`
	Diagnostics     []string     `json:"diagnostics"`
	SemanticContext string       `json:"semantic_context,omitempty"`
}

func ExploreCodeGraph(ctx context.Context, root, query string) (string, error) {
	available := false
	for _, c := range Discover(root) {
		if c.Provider == "codegraph" {
			available = c.Available
		}
	}
	if !available {
		return "", fmt.Errorf("CODEGRAPH_UNAVAILABLE: existing .codegraph index and executable are required")
	}
	if strings.TrimSpace(query) == "" {
		query = "Summarize declarations, references, call paths, inputs, outputs, and side effects for changed code"
	}
	cmd := exec.CommandContext(ctx, "codegraph", "explore", query)
	cmd.Dir = root
	b, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("CODEGRAPH_QUERY_FAILED: %w", err)
	}
	if len(b) > 1<<20 {
		b = b[:1<<20]
	}
	return string(b), nil
}

func Discover(root string) []Capability {
	_, cgDir := os.Stat(filepath.Join(root, ".codegraph"))
	_, cgBin := exec.LookPath("codegraph")
	codegraph := Capability{Provider: "codegraph", Available: cgDir == nil && cgBin == nil, Calls: true, Imports: true, DataFlow: true, Confidence: .9}
	if cgDir != nil {
		codegraph.Reason = ".codegraph directory is absent; indexing is never automatic"
	} else if cgBin != nil {
		codegraph.Reason = "codegraph executable is unavailable"
	}
	_, goErr := exec.LookPath("go")
	return []Capability{codegraph, {Provider: "go-ast", Available: goErr == nil, Languages: []string{"go"}, Calls: true, Imports: true, DataFlow: false, Confidence: .8, Reason: func() string {
		if goErr != nil {
			return "go executable is unavailable"
		}
		return ""
	}()}, {Provider: "filesystem", Available: true, Confidence: .4}}
}

func Analyze(ctx context.Context, root string, requested []string, changed bool) (Report, error) {
	report := Report{Providers: Discover(root), Files: []File{}, Components: []Component{}, Edges: []Edge{}, Diagnostics: []string{}}
	paths, err := sourcePaths(ctx, root, requested, changed)
	if err != nil {
		return report, err
	}
	for _, path := range paths {
		if filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
			continue
		}
		info, err := os.Stat(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil || info.Size() > MaxFileSize {
			report.Diagnostics = append(report.Diagnostics, "skipped "+path+": unreadable or oversized")
			continue
		}
		layer, reason, confidence := classify(path)
		report.Files = append(report.Files, File{ID: "file:" + path, Path: path, Layer: layer, LayerReason: reason, Confidence: confidence})
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, filepath.Join(root, filepath.FromSlash(path)), nil, parser.SkipObjectResolution)
		if err != nil {
			report.Diagnostics = append(report.Diagnostics, fmt.Sprintf("parse %s: %v", path, err))
			continue
		}
		imports := make([]string, 0, len(file.Imports))
		for _, im := range file.Imports {
			imports = append(imports, strings.Trim(im.Path.Value, `"`))
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			component := Component{ID: "symbol:" + path + ":" + fn.Name.Name, Name: fn.Name.Name, Kind: "function", Locator: fmt.Sprintf("%s:%d", path, fset.Position(fn.Pos()).Line), File: path, PrimaryLayer: layer, LayerReason: reason, Imports: append([]string(nil), imports...), References: []string{}, Calls: []string{}, Inputs: fieldShapes(fn.Type.Params), Outputs: fieldShapes(fn.Type.Results), Effects: []string{}, Unknown: []string{"incoming references require a semantic index"}, Provider: "go-ast", Confidence: confidence}
			if fn.Recv != nil {
				component.Kind = "method"
				component.Inputs = append(fieldShapes(fn.Recv), component.Inputs...)
			}
			if fn.Body != nil {
				ast.Inspect(fn.Body, func(n ast.Node) bool {
					switch x := n.(type) {
					case *ast.CallExpr:
						component.Calls = append(component.Calls, exprName(x.Fun))
					case *ast.AssignStmt:
						component.Effects = append(component.Effects, "writes local or selected values")
					case *ast.SendStmt:
						component.Effects = append(component.Effects, "sends channel value")
					case *ast.GoStmt:
						component.Effects = append(component.Effects, "starts goroutine")
					}
					return true
				})
			}
			component.Calls = clean(component.Calls)
			component.Effects = clean(component.Effects)
			if len(component.Outputs) == 0 {
				component.Outputs = []string{"none"}
			}
			if len(component.Effects) == 0 {
				component.Effects = []string{"no statically detected side effect"}
			}
			report.Components = append(report.Components, component)
			report.Edges = append(report.Edges, Edge{From: "file:" + path, To: component.ID, Kind: "declares", Locator: component.Locator, Provider: "go-ast", Confidence: .95})
			for _, call := range component.Calls {
				report.Edges = append(report.Edges, Edge{From: component.ID, To: "call:" + call, Kind: "calls", Locator: component.Locator, Provider: "go-ast", Confidence: .65})
			}
		}
	}
	sort.Slice(report.Files, func(i, j int) bool { return report.Files[i].Path < report.Files[j].Path })
	sort.Slice(report.Components, func(i, j int) bool { return report.Components[i].ID < report.Components[j].ID })
	return report, nil
}

func sourcePaths(ctx context.Context, root string, requested []string, changed bool) ([]string, error) {
	if len(requested) > 0 {
		var out []string
		for _, p := range requested {
			p = filepath.ToSlash(filepath.Clean(p))
			if p == ".." || strings.HasPrefix(p, "../") || filepath.IsAbs(p) {
				return nil, fmt.Errorf("CODE_PATH_OUTSIDE_ROOT: %s", p)
			}
			out = append(out, p)
		}
		return clean(out), nil
	}
	if changed {
		cmd := exec.CommandContext(ctx, "git", "status", "--porcelain=v1")
		cmd.Dir = root
		b, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("CODE_GIT_STATUS: %w", err)
		}
		var out []string
		s := bufio.NewScanner(strings.NewReader(string(b)))
		for s.Scan() {
			line := s.Text()
			if len(line) > 3 {
				p := strings.TrimSpace(line[3:])
				if i := strings.Index(p, " -> "); i >= 0 {
					p = p[i+4:]
				}
				out = append(out, filepath.ToSlash(p))
			}
		}
		return clean(out), s.Err()
	}
	var out []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			n := d.Name()
			if n == ".git" || n == ".tene-workflow" || n == "vendor" || n == "node_modules" || n == "dist" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) == ".go" {
			rel, _ := filepath.Rel(root, path)
			out = append(out, filepath.ToSlash(rel))
		}
		return nil
	})
	return clean(out), err
}

func classify(path string) (string, string, float64) {
	p := strings.ToLower(filepath.ToSlash(path))
	rules := []struct {
		keys          []string
		layer, reason string
	}{{[]string{"cmd/", "controller", "route", "handler", "ui/"}, "Interface", "entry-point path/name"}, {[]string{"repository", "storage", "state/", "database", "cache", "queue", "client"}, "Persistence", "data boundary path/name"}, {[]string{".github/", "infra/", "deploy", "config", "server", "auth"}, "Infrastructure", "runtime path/name"}, {[]string{"domain/", "service", "usecase", "workflow/", "internal/"}, "Business Logic", "processing-rule path/name"}}
	for _, r := range rules {
		for _, k := range r.keys {
			if strings.Contains(p, k) {
				return r.layer, r.reason, .55
			}
		}
	}
	return "Unknown", "no deterministic layer rule matched", .3
}
func fieldShapes(f *ast.FieldList) []string {
	if f == nil {
		return []string{}
	}
	var out []string
	for _, x := range f.List {
		shape := exprName(x.Type)
		if len(x.Names) == 0 {
			out = append(out, shape)
		} else {
			for _, n := range x.Names {
				out = append(out, n.Name+" "+shape)
			}
		}
	}
	return out
}
func exprName(e ast.Expr) string {
	switch x := e.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.SelectorExpr:
		return exprName(x.X) + "." + x.Sel.Name
	case *ast.StarExpr:
		return "*" + exprName(x.X)
	case *ast.ArrayType:
		return "[]" + exprName(x.Elt)
	case *ast.MapType:
		return "map[" + exprName(x.Key) + "]" + exprName(x.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.FuncType:
		return "func"
	case *ast.IndexExpr:
		return exprName(x.X) + "[" + exprName(x.Index) + "]"
	case *ast.IndexListExpr:
		return exprName(x.X) + "[...]"
	case *ast.ChanType:
		return "chan " + exprName(x.Value)
	default:
		return fmt.Sprintf("%T", e)
	}
}
func clean(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, v := range in {
		v = strings.TrimSpace(v)
		if v != "" && !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	sort.Strings(out)
	return out
}
