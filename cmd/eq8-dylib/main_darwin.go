//go:build darwin && cgo

package main

import (
	"github.com/th-release/vst3go/example/eq8"
	_ "github.com/th-release/vst3go/pkg/plugin/cbridge"
	"github.com/th-release/vst3go/pkg/plugin"
)

func init() {
	plugin.Register(eq8.NewPlugin())
}

func main() {}
