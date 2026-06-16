//go:build windows && cgo

package plugin

/*
#cgo windows CFLAGS: -I../../include
#cgo windows LDFLAGS: -lWebView2Loader
#include "../../include/vst3/vst3_c_api.h"
#include "../../bridge/editor_view.h"
*/
import "C"
