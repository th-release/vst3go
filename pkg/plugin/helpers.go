package plugin

// #cgo CFLAGS: -I../../include
// #include "../../include/vst3/vst3_c_api.h"
//
// // Helper to access channelBuffers32 from the union
// static inline float** getChannelBuffers32(struct Steinberg_Vst_AudioBusBuffers* bus) {
//     return bus->Steinberg_Vst_AudioBusBuffers_channelBuffers32;
// }
import "C"
import "unsafe"

// getChannelBuffers32 extracts the 32-bit channel buffers from an audio bus
func getChannelBuffers32(bus *C.struct_Steinberg_Vst_AudioBusBuffers) **C.float {
	return C.getChannelBuffers32(bus)
}

func copyStringToUTF16(src string, dst []uint16, maxLen int) {
	if maxLen <= 0 || len(dst) == 0 {
		return
	}

	runes := []rune(src)
	limit := len(runes)
	if limit > maxLen-1 {
		limit = maxLen - 1
	}
	if limit > len(dst)-1 {
		limit = len(dst) - 1
	}

	for i := 0; i < limit; i++ {
		dst[i] = uint16(runes[i])
	}
	dst[limit] = 0
}

func stringFromUTF16(src []uint16) string {
	length := 0
	for length < len(src) && src[length] != 0 {
		length++
	}

	runes := make([]rune, length)
	for i := 0; i < length; i++ {
		runes[i] = rune(src[i])
	}
	return string(runes)
}

// copyStringToTChar copies a Go string to a VST3 TChar (UTF16) buffer
func copyStringToTChar(src string, dst *C.Steinberg_Vst_TChar, maxLen int) {
	if dst == nil || maxLen <= 0 {
		return
	}

	buffer := unsafe.Slice((*uint16)(unsafe.Pointer(dst)), maxLen)
	copyStringToUTF16(src, buffer, maxLen)
}

// stringFromTChar converts a VST3 TChar (UTF16) buffer to a Go string
func stringFromTChar(src *C.Steinberg_Vst_TChar) string {
	if src == nil {
		return ""
	}

	// Count length
	length := 0
	for {
		ch := *(*C.Steinberg_char16)(unsafe.Pointer(
			uintptr(unsafe.Pointer(src)) + uintptr(length*2)))
		if ch == 0 {
			break
		}
		length++
	}

	buffer := unsafe.Slice((*uint16)(unsafe.Pointer(src)), length)
	return stringFromUTF16(buffer)
}
