package plugin

import (
	"embed"
	"strings"
)

//go:embed editor_assets/index.html editor_assets/editor.css editor_assets/editor.js
var editorAssets embed.FS

func loadEditorAsset(path string) string {
	data, err := editorAssets.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

func buildEditorHTML(encodedSnapshot string) string {
	html := loadEditorAsset("editor_assets/index.html")
	css := loadEditorAsset("editor_assets/editor.css")
	script := loadEditorAsset("editor_assets/editor.js")

	html = strings.ReplaceAll(html, "__VST3GO_CSS__", css)
	html = strings.ReplaceAll(html, "__VST3GO_JS__", script)
	html = strings.ReplaceAll(html, "__VST3GO_SNAPSHOT__", encodedSnapshot)
	return html
}
