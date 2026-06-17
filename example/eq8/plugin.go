package eq8

import (
	frameworkplugin "github.com/th-release/vst3go/pkg/framework/plugin"
	"github.com/th-release/vst3go/pkg/plugin"
	_ "github.com/th-release/vst3go/pkg/plugin/cbridge"
)

var _ plugin.Plugin = (*Plugin)(nil)

// Plugin is the vst3go entrypoint for the EQ8 example.
type Plugin struct{}

// GetInfo returns the host-visible plugin metadata.
func (p *Plugin) GetInfo() frameworkplugin.Info {
	return frameworkplugin.Info{
		ID:       "com.example.vst3go.eq8",
		Name:     "EQ8 Example",
		Version:  "1.0.0",
		Vendor:   "Example Audio",
		Category: "Fx",
	}
}

// CreateProcessor constructs a new EQ processor instance.
func (p *Plugin) CreateProcessor() plugin.Processor {
	return NewProcessor()
}

// NewPlugin returns a ready-to-use plugin entrypoint.
func NewPlugin() *Plugin {
	return &Plugin{}
}
