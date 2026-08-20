package projectconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProbeCodexDiscoversPortablePluginSurfaces(t *testing.T) {
	r := t.TempDir()
	for _, p := range []string{".codex-plugin/plugin.json", "skills/a/SKILL.md", "hooks/hooks.json", ".codex/agents/evaluator.toml"} {
		f := filepath.Join(r, filepath.FromSlash(p))
		_ = os.MkdirAll(filepath.Dir(f), 0755)
		body := "x"
		if filepath.Ext(f) == ".json" {
			body = "{}"
		}
		_ = os.WriteFile(f, []byte(body), 0644)
	}
	c := ProbeCodex(r)
	if !c.PluginManifest || c.Skills != 1 || !c.Hooks || c.SubagentProfiles != 1 {
		t.Fatalf("%+v", c)
	}
}
