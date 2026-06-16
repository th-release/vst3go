//go:build windows && cgo

package cbridge

/*
#include "../../../bridge/windows_dll.c"
*/
import "C"
