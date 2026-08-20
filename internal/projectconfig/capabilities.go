package projectconfig

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type CodexCapabilities struct {
	CLI              bool     `json:"cli"`
	Version          string   `json:"version,omitempty"`
	PluginManifest   bool     `json:"plugin_manifest"`
	Skills           int      `json:"skills"`
	Hooks            bool     `json:"hooks"`
	SubagentProfiles int      `json:"subagent_profiles"`
	AppServer        bool     `json:"app_server"`
	MCPConfig        bool     `json:"mcp_config"`
	ProbeErrors      []string `json:"probe_errors"`
}

func ProbeCodex(root string) CodexCapabilities {
	c := CodexCapabilities{}
	if p, e := exec.LookPath("codex"); e == nil {
		c.CLI = true
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if b, e := exec.CommandContext(ctx, p, "--version").CombinedOutput(); e == nil {
			c.Version = strings.TrimSpace(string(b))
		} else {
			c.ProbeErrors = append(c.ProbeErrors, "codex version probe failed")
		}
		ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel2()
		if e := exec.CommandContext(ctx2, p, "app-server", "--help").Run(); e == nil {
			c.AppServer = true
		} else {
			c.ProbeErrors = append(c.ProbeErrors, "app-server unavailable")
		}
	}
	manifest := filepath.Join(root, ".codex-plugin", "plugin.json")
	if b, e := os.ReadFile(manifest); e == nil {
		var v map[string]any
		if json.Unmarshal(b, &v) == nil {
			c.PluginManifest = true
		}
	}
	if ds, e := os.ReadDir(filepath.Join(root, "skills")); e == nil {
		for _, d := range ds {
			if d.IsDir() {
				if _, e := os.Stat(filepath.Join(root, "skills", d.Name(), "SKILL.md")); e == nil {
					c.Skills++
				}
			}
		}
	}
	if _, e := os.Stat(filepath.Join(root, "hooks", "hooks.json")); e == nil {
		c.Hooks = true
	}
	if ds, e := os.ReadDir(filepath.Join(root, ".codex", "agents")); e == nil {
		for _, d := range ds {
			if filepath.Ext(d.Name()) == ".toml" {
				c.SubagentProfiles++
			}
		}
	}
	if _, e := os.Stat(filepath.Join(root, ".mcp.json")); e == nil {
		c.MCPConfig = true
	}
	return c
}
