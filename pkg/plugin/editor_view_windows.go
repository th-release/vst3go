//go:build windows

package plugin

/*
#cgo windows CFLAGS: -I../../include
#include "../../include/vst3/vst3_c_api.h"
#include "../../bridge/editor_view.h"
*/
import "C"
