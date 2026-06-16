//go:build windows && cgo && amd64

package cbridge

/*
#include "../../../bridge/windows_dll.c"
*/
import "C"
