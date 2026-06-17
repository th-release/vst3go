package eq8

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	vst3plugin "github.com/th-release/vst3go/pkg/plugin"
)

//go:embed web/editor/index.html web/editor/editor.css web/editor/editor.js
var editorAssets embed.FS

func loadEditorAsset(path string) string {
	data, err := editorAssets.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

// BuildEditorHTML injects the encoded snapshot into the example editor shell.
func BuildEditorHTML(encodedSnapshot string) string {
	html := loadEditorAsset("web/editor/index.html")
	css := loadEditorAsset("web/editor/editor.css")
	script := loadEditorAsset("web/editor/editor.js")

	html = strings.ReplaceAll(html, `<link rel="stylesheet" crossorigin href="./editor.css">`, "<style>"+css+"</style>")
	html = strings.ReplaceAll(html, `<script type="module" crossorigin src="./editor.js"></script>`, "<script>"+script+"</script>")
	html = strings.ReplaceAll(html, "__EQ8_SNAPSHOT__", encodedSnapshot)
	return html
}

// RenderEditorHTML marshals a snapshot and returns the full HTML shell.
func RenderEditorHTML(snapshot *vst3plugin.EditorSnapshot) (string, error) {
	if snapshot == nil {
		return "", fmt.Errorf("render editor html: nil snapshot")
	}

	data, err := json.Marshal(snapshot)
	if err != nil {
		return "", fmt.Errorf("render editor html: marshal snapshot: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	return BuildEditorHTML(encoded), nil
}
