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
	"regexp"
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
	}()}, {Provider: "typescript-static", Available: true, Languages: []string{"javascript", "typescript", "jsx", "tsx"}, Calls: true, Imports: true, DataFlow: true, Confidence: .7, Reason: "bounded syntax analysis; runtime dispatch remains observational"}, {Provider: "python-static", Available: true, Languages: []string{"python"}, Calls: true, Imports: true, DataFlow: true, Confidence: .7, Reason: "bounded syntax analysis; dynamic dispatch remains observational"}, {Provider: "filesystem", Available: true, Confidence: .4}}
}

func Analyze(ctx context.Context, root string, requested []string, changed bool) (Report, error) {
	report := Report{Providers: Discover(root), Files: []File{}, Components: []Component{}, Edges: []Edge{}, Diagnostics: []string{}}
	paths, err := sourcePaths(ctx, root, requested, changed)
	if err != nil {
		return report, err
	}
	for _, path := range paths {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		info, err := os.Stat(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil || info.Size() > MaxFileSize {
			report.Diagnostics = append(report.Diagnostics, "skipped "+path+": unreadable or oversized")
			continue
		}
		layer, reason, confidence := classify(path)
		report.Files = append(report.Files, File{ID: "file:" + path, Path: path, Layer: layer, LayerReason: reason, Confidence: confidence})
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".js" || ext == ".jsx" || ext == ".ts" || ext == ".tsx" || ext == ".py" {
			components, edges, diagnostics := analyzeBoundedLanguage(filepath.Join(root, filepath.FromSlash(path)), path, layer, reason, confidence, ext)
			report.Components = append(report.Components, components...)
			report.Edges = append(report.Edges, edges...)
			report.Diagnostics = append(report.Diagnostics, diagnostics...)
			continue
		}
		if ext != ".go" {
			report.Components = append(report.Components, Component{ID: "file-component:" + path, Name: filepath.Base(path), Kind: "file", Locator: path, File: path, PrimaryLayer: layer, LayerReason: reason, Imports: []string{}, References: []string{}, Calls: []string{}, Inputs: []string{"unknown: language semantic provider unavailable"}, Outputs: []string{"unknown: language semantic provider unavailable"}, Effects: []string{"unknown: runtime observation required"}, Unknown: []string{"declarations", "imports and references", "callers and callees", "input and output shapes", "side effects"}, Provider: "filesystem", Confidence: .4})
			report.Diagnostics = append(report.Diagnostics, "degraded filesystem analysis for "+path+": semantic provider unavailable; Six Questions remain explicit unknowns")
			continue
		}
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

var (
	tsImportPattern   = regexp.MustCompile(`(?:import\s+.*?\s+from\s+|require\s*\()?["']([^"']+)["']`)
	tsFunctionPattern = regexp.MustCompile(`(?:export\s+)?(?:default\s+)?(?:async\s+)?function\s+([A-Za-z_$][\w$]*)\s*\(([^)]*)\)|(?:export\s+)?(?:const|let|var)\s+([A-Za-z_$][\w$]*)\s*=\s*(?:async\s*)?\(([^)]*)\)\s*=>`)
	pyImportPattern   = regexp.MustCompile(`^\s*(?:from\s+([\w.]+)\s+import|import\s+([\w.]+))`)
	pyFunctionPattern = regexp.MustCompile(`^\s*(?:async\s+)?def\s+([A-Za-z_]\w*)\s*\(([^)]*)\)\s*(?:->\s*([^:]+))?:`)
	callPattern       = regexp.MustCompile(`([A-Za-z_$][\w$]*(?:\.[A-Za-z_$][\w$]*)*)\s*\(`)
)

// analyzeBoundedLanguage intentionally extracts only stable, reviewable syntax. It never
// claims dynamic dispatch; runtime QA supplies that missing edge.
func analyzeBoundedLanguage(abs, path, layer, reason string, confidence float64, ext string) ([]Component, []Edge, []string) {
	b, err := os.ReadFile(abs)
	if err != nil {
		return nil, nil, []string{"read " + path + ": " + err.Error()}
	}
	provider := "typescript-static"
	imports := []string{}
	if ext == ".py" {
		provider = "python-static"
	}
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		if ext == ".py" {
			if m := pyImportPattern.FindStringSubmatch(line); len(m) > 0 {
				imports = append(imports, firstNonEmpty(m[1:]...))
			}
		} else if strings.Contains(line, "import") || strings.Contains(line, "require") {
			if m := tsImportPattern.FindStringSubmatch(line); len(m) > 1 {
				imports = append(imports, m[1])
			}
		}
	}
	imports = clean(imports)
	var components []Component
	var edges []Edge
	for lineNo, line := range lines {
		var name, params, output string
		if ext == ".py" {
			m := pyFunctionPattern.FindStringSubmatch(line)
			if len(m) == 0 {
				continue
			}
			name, params, output = m[1], m[2], strings.TrimSpace(m[3])
		} else {
			m := tsFunctionPattern.FindStringSubmatch(line)
			if len(m) == 0 {
				continue
			}
			name, params = firstNonEmpty(m[1], m[3]), firstNonEmpty(m[2], m[4])
			output = tsReturnShape(line)
		}
		inputs := parameterShapes(params)
		if len(inputs) == 0 {
			inputs = []string{"none"}
		}
		if output == "" {
			output = "runtime value or none"
		}
		body := line
		if ext == ".py" {
			indent := len(line) - len(strings.TrimLeft(line, " \t"))
			for i := lineNo + 1; i < len(lines); i++ {
				next := lines[i]
				if strings.TrimSpace(next) == "" {
					continue
				}
				nextIndent := len(next) - len(strings.TrimLeft(next, " \t"))
				if nextIndent <= indent {
					break
				}
				body += "\n" + next
			}
		}
		calls := []string{}
		for _, m := range callPattern.FindAllStringSubmatch(body, -1) {
			if m[1] != name && !isControlCall(m[1]) {
				calls = append(calls, m[1])
			}
		}
		effects := boundedEffects(body)
		component := Component{ID: "symbol:" + path + ":" + name, Name: name, Kind: "function", Locator: fmt.Sprintf("%s:%d", path, lineNo+1), File: path, PrimaryLayer: layer, LayerReason: reason, Imports: imports, References: []string{}, Calls: clean(calls), Inputs: inputs, Outputs: []string{output}, Effects: effects, Unknown: []string{"incoming references and dynamic dispatch require semantic index or runtime observation"}, Provider: provider, Confidence: min(confidence+.1, .75)}
		components = append(components, component)
		edges = append(edges, Edge{From: "file:" + path, To: component.ID, Kind: "declares", Locator: component.Locator, Provider: provider, Confidence: .85})
		for _, call := range component.Calls {
			edges = append(edges, Edge{From: component.ID, To: "call:" + call, Kind: "calls", Locator: component.Locator, Provider: provider, Confidence: .6})
		}
	}
	if len(components) == 0 {
		return nil, nil, []string{"bounded " + provider + " analysis found no function declarations in " + path}
	}
	return components, edges, nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
func parameterShapes(raw string) []string {
	var out []string
	for _, p := range strings.Split(raw, ",") {
		p = strings.TrimSpace(p)
		if p != "" && p != "self" && p != "cls" {
			out = append(out, p)
		}
	}
	return clean(out)
}
func tsReturnShape(line string) string {
	if i := strings.Index(line, "): "); i >= 0 {
		tail := line[i+3:]
		if j := strings.Index(tail, "{"); j >= 0 {
			return strings.TrimSpace(tail[:j])
		}
	}
	if strings.Contains(line, "=>") {
		return "inferred expression or declared return"
	}
	return ""
}
func isControlCall(name string) bool {
	switch name {
	case "if", "for", "while", "switch", "catch", "function", "def":
		return true
	}
	return false
}
func boundedEffects(body string) []string {
	var out []string
	lower := strings.ToLower(body)
	for needle, effect := range map[string]string{"write": "writes persistence or response data", "save": "persists data", "insert": "persists data", "publish": "publishes queued data", "append": "mutates a collection", "fetch": "performs external or HTTP I/O", "request": "performs external or HTTP I/O"} {
		if strings.Contains(lower, needle) {
			out = append(out, effect)
		}
	}
	if len(out) == 0 {
		out = []string{"no bounded static side effect detected"}
	}
	return clean(out)
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
			if n == ".git" || n == ".tene" || n == ".tene-workflow" || n == "vendor" || n == "node_modules" || n == "dist" {
				return filepath.SkipDir
			}
			return nil
		}
		if sourceExtension(filepath.Ext(path)) {
			rel, _ := filepath.Rel(root, path)
			out = append(out, filepath.ToSlash(rel))
		}
		return nil
	})
	return clean(out), err
}

func sourceExtension(ext string) bool {
	switch strings.ToLower(ext) {
	case ".go", ".js", ".jsx", ".ts", ".tsx", ".py", ".java", ".kt", ".rb", ".rs", ".php", ".cs", ".swift":
		return true
	}
	return false
}

func classify(path string) (string, string, float64) {
	p := strings.ToLower(filepath.ToSlash(path))
	rules := []struct {
		keys          []string
		layer, reason string
	}{{[]string{"cmd/", "controller", "route", "handler", "ui/", "frontend/", "api."}, "Interface", "entry-point path/name"}, {[]string{"repository", "storage", "store", "state/", "database", "cache", "queue", "client"}, "Persistence", "data boundary path/name"}, {[]string{".github/", "infra/", "deploy", "config", "server", "auth"}, "Infrastructure", "runtime path/name"}, {[]string{"domain/", "service", "usecase", "workflow/", "worker/", "internal/"}, "Business Logic", "processing-rule path/name"}}
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
