//go:build windows && cgo && amd64

package main

import (
	_ "github.com/cwbudde/vst3go/pkg/plugin"
	_ "github.com/cwbudde/vst3go/pkg/plugin/cbridge"
)

func main() {}
