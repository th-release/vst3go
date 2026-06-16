package plugin

import (
	"fmt"

	"github.com/cwbudde/vst3go/pkg/framework/param"
	frameworkplugin "github.com/cwbudde/vst3go/pkg/framework/plugin"
)

// EditorModel is the browser-friendly snapshot of the plugin editor.
//
// The web renderer can consume this structure directly to build controls,
// display the current layout, and keep values in sync with the processor.
type EditorModel struct {
	Plugin   EditorPluginInfo `json:"plugin"`
	Sections []EditorSection  `json:"sections"`
}

// EditorPluginInfo mirrors the metadata the editor needs for display.
type EditorPluginInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Version  string `json:"version"`
	Vendor   string `json:"vendor"`
	Category string `json:"category"`
}

// EditorSection groups related controls for the renderer.
type EditorSection struct {
	Title    string          `json:"title"`
	Controls []EditorControl `json:"controls"`
}

// EditorControl describes a single parameter control in the editor surface.
type EditorControl struct {
	ID           uint32  `json:"id"`
	Name         string  `json:"name"`
	ShortName    string  `json:"shortName"`
	Unit         string  `json:"unit"`
	Kind         string  `json:"kind"`
	Normalized   float64 `json:"normalized"`
	Plain        float64 `json:"plain"`
	Min          float64 `json:"min"`
	Max          float64 `json:"max"`
	DefaultValue float64 `json:"defaultValue"`
	StepCount    int32   `json:"stepCount"`
	Flags        uint32  `json:"flags"`
	ReadOnly     bool    `json:"readOnly"`
	Hidden       bool    `json:"hidden"`
}

// BuildEditorModel turns plugin metadata and parameters into a renderable
// browser-friendly snapshot.
func BuildEditorModel(info frameworkplugin.Info, registry *param.Registry) (*EditorModel, error) {
	if registry == nil {
		return nil, fmt.Errorf("build editor model: nil registry")
	}

	model := &EditorModel{
		Plugin: EditorPluginInfo{
			ID:       info.ID,
			Name:     info.Name,
			Version:  info.Version,
			Vendor:   info.Vendor,
			Category: info.Category,
		},
	}

	controls := registry.All()
	section := EditorSection{
		Title:    "Parameters",
		Controls: make([]EditorControl, 0, len(controls)),
	}

	for _, p := range controls {
		if p == nil {
			continue
		}

		normalized := p.GetNormalized()
		control := EditorControl{
			ID:           p.ID,
			Name:         p.Name,
			ShortName:    p.ShortName,
			Unit:         p.Unit,
			Kind:         inferControlKind(p),
			Normalized:   normalized,
			Plain:        p.Denormalize(normalized),
			Min:          p.Min,
			Max:          p.Max,
			DefaultValue: p.DefaultValue,
			StepCount:    p.StepCount,
			Flags:        p.Flags,
			ReadOnly:     p.Flags&param.IsReadOnly != 0,
			Hidden:       p.Flags&param.IsHidden != 0,
		}
		section.Controls = append(section.Controls, control)
	}

	model.Sections = []EditorSection{section}
	return model, nil
}

func inferControlKind(p *param.Parameter) string {
	switch {
	case p == nil:
		return "unknown"
	case p.Flags&param.IsHidden != 0:
		return "hidden"
	case p.Flags&param.IsReadOnly != 0:
		return "meter"
	case p.Flags&param.IsBypass != 0:
		return "toggle"
	case p.Flags&param.IsList != 0 || p.StepCount > 0:
		return "choice"
	default:
		return "slider"
	}
}
