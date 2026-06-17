package eq8

import (
	"strings"
	"testing"

	frameworkparam "github.com/th-release/vst3go/pkg/framework/param"
	frameworkplugin "github.com/th-release/vst3go/pkg/framework/plugin"
	vst3plugin "github.com/th-release/vst3go/pkg/plugin"
)

func TestBuildEditorHTMLInjectsAssets(t *testing.T) {
	html := BuildEditorHTML("c2FtcGxl")

	if strings.Contains(html, "__EQ8_CSS__") {
		t.Fatal("CSS placeholder was not replaced")
	}
	if strings.Contains(html, "__EQ8_JS__") {
		t.Fatal("JS placeholder was not replaced")
	}
	if strings.Contains(html, "__EQ8_SNAPSHOT__") {
		t.Fatal("snapshot placeholder was not replaced")
	}
	if !strings.Contains(html, "c2FtcGxl") {
		t.Fatal("encoded snapshot was not embedded")
	}
	if !strings.Contains(html, ".band-grid") {
		t.Fatal("embedded CSS was not included")
	}
	if !strings.Contains(html, "renderGraph()") {
		t.Fatal("embedded JS was not included")
	}
}

func TestRenderEditorHTMLEmbedsSnapshot(t *testing.T) {
	registry := frameworkparam.NewRegistry()
	if err := registry.Add(frameworkparam.GainParameter(1, "Input Gain").Build()); err != nil {
		t.Fatalf("registry.Add() failed: %v", err)
	}

	snapshot, err := vst3plugin.BuildEditorSnapshot(frameworkplugin.Info{
		ID:   "com.example.vst3go.eq8",
		Name: "EQ8 Example",
	}, registry)
	if err != nil {
		t.Fatalf("BuildEditorSnapshot() failed: %v", err)
	}

	html, err := RenderEditorHTML(snapshot)
	if err != nil {
		t.Fatalf("RenderEditorHTML() failed: %v", err)
	}

	if !strings.Contains(html, "EQ8 Example") {
		t.Fatal("plugin metadata was not embedded")
	}
}
