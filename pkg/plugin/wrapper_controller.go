package plugin

// #cgo CFLAGS: -I../../include
// #include "../../include/vst3/vst3_c_api.h"
// #include "../../bridge/editor_view.h"
import "C"

import (
	"encoding/base64"
	"encoding/json"
	"unsafe"

	"github.com/th-release/vst3go/pkg/vst3"
)

// IEditController callbacks
//
//export GoEditControllerSetComponentState
func GoEditControllerSetComponentState(componentPtr unsafe.Pointer, state unsafe.Pointer) C.Steinberg_tresult {
	// Component state received from processor - apply to edit controller
	return GoComponentSetState(componentPtr, state)
}

//export GoEditControllerSetState
func GoEditControllerSetState(componentPtr unsafe.Pointer, state unsafe.Pointer) C.Steinberg_tresult {
	// Edit controller state is typically the same as component state
	return GoComponentSetState(componentPtr, state)
}

//export GoEditControllerGetState
func GoEditControllerGetState(componentPtr unsafe.Pointer, state unsafe.Pointer) C.Steinberg_tresult {
	// Edit controller state is typically the same as component state
	return GoComponentGetState(componentPtr, state)
}

//export GoEditControllerGetParameterCount
func GoEditControllerGetParameterCount(componentPtr unsafe.Pointer) C.int32_t {
	wrapper := getComponent(uintptr(componentPtr))
	if wrapper == nil {
		return 0
	}

	return C.int32_t(wrapper.component.GetParameterCount())
}

//export GoEditControllerGetParameterInfo
func GoEditControllerGetParameterInfo(componentPtr unsafe.Pointer, paramIndex C.int32_t, info *C.struct_Steinberg_Vst_ParameterInfo) C.Steinberg_tresult {
	wrapper := getComponent(uintptr(componentPtr))
	if wrapper == nil {
		return C.Steinberg_tresult(vst3.ResultFalse)
	}

	paramInfo, err := wrapper.component.GetParameterInfo(int32(paramIndex))
	if err != nil || paramInfo == nil {
		return C.Steinberg_tresult(vst3.ResultFalse)
	}

	// Copy to C struct
	cInfo := info
	cInfo.id = C.Steinberg_Vst_ParamID(paramInfo.ID)

	// Copy title
	titleBytes := []byte(paramInfo.Title)
	if len(titleBytes) > 127 {
		titleBytes = titleBytes[:127]
	}
	for i, b := range titleBytes {
		cInfo.title[i] = C.Steinberg_char16(b)
	}
	cInfo.title[len(titleBytes)] = 0

	// Copy short title
	shortTitleBytes := []byte(paramInfo.ShortTitle)
	if len(shortTitleBytes) > 127 {
		shortTitleBytes = shortTitleBytes[:127]
	}
	for i, b := range shortTitleBytes {
		cInfo.shortTitle[i] = C.Steinberg_char16(b)
	}
	cInfo.shortTitle[len(shortTitleBytes)] = 0

	// Copy units
	unitsBytes := []byte(paramInfo.Units)
	if len(unitsBytes) > 127 {
		unitsBytes = unitsBytes[:127]
	}
	for i, b := range unitsBytes {
		cInfo.units[i] = C.Steinberg_char16(b)
	}
	cInfo.units[len(unitsBytes)] = 0

	cInfo.stepCount = C.Steinberg_int32(paramInfo.StepCount)
	cInfo.defaultNormalizedValue = C.Steinberg_Vst_ParamValue(paramInfo.DefaultValue)
	cInfo.unitId = C.Steinberg_Vst_UnitID(paramInfo.UnitID)
	cInfo.flags = C.Steinberg_int32(paramInfo.Flags)

	// Debug: Print parameter info
	// fmt.Printf("Parameter %d: StepCount=%d, Flags=%d, Name=%s\n", paramInfo.ID, paramInfo.StepCount, paramInfo.Flags, paramInfo.Title)

	return C.Steinberg_tresult(vst3.ResultOK)
}

//export GoEditControllerGetParamStringByValue
func GoEditControllerGetParamStringByValue(componentPtr unsafe.Pointer, id C.Steinberg_Vst_ParamID, valueNormalized C.Steinberg_Vst_ParamValue, string *C.Steinberg_Vst_TChar) C.Steinberg_tresult {
	wrapper := getComponent(uintptr(componentPtr))
	if wrapper == nil {
		return C.Steinberg_tresult(vst3.ResultFalse)
	}

	// Get the formatted string
	str, err := wrapper.component.GetParamStringByValue(uint32(id), float64(valueNormalized))
	if err != nil {
		return C.Steinberg_tresult(vst3.ResultFalse)
	}

	// Convert to UTF16 for VST3
	copyStringToTChar(str, string, 128)

	return C.Steinberg_tresult(vst3.ResultOK)
}

//export GoEditControllerGetParamValueByString
func GoEditControllerGetParamValueByString(componentPtr unsafe.Pointer, id C.Steinberg_Vst_ParamID, string *C.Steinberg_Vst_TChar, valueNormalized *C.Steinberg_Vst_ParamValue) C.Steinberg_tresult {
	wrapper := getComponent(uintptr(componentPtr))
	if wrapper == nil {
		return C.Steinberg_tresult(vst3.ResultFalse)
	}

	// Convert from UTF16
	str := stringFromTChar(string)

	// Parse the value
	value, err := wrapper.component.GetParamValueByString(uint32(id), str)
	if err != nil {
		return C.Steinberg_tresult(vst3.ResultFalse)
	}

	*valueNormalized = C.Steinberg_Vst_ParamValue(value)
	return C.Steinberg_tresult(vst3.ResultOK)
}

//export GoEditControllerNormalizedParamToPlain
func GoEditControllerNormalizedParamToPlain(componentPtr unsafe.Pointer, id C.uint32_t, valueNormalized C.double) C.double {
	wrapper := getComponent(uintptr(componentPtr))
	if wrapper == nil {
		return valueNormalized
	}

	return C.double(wrapper.component.NormalizedParamToPlain(uint32(id), float64(valueNormalized)))
}

//export GoEditControllerPlainParamToNormalized
func GoEditControllerPlainParamToNormalized(componentPtr unsafe.Pointer, id C.uint32_t, plainValue C.double) C.double {
	wrapper := getComponent(uintptr(componentPtr))
	if wrapper == nil {
		return plainValue
	}

	return C.double(wrapper.component.PlainParamToNormalized(uint32(id), float64(plainValue)))
}

//export GoEditControllerGetParamNormalized
func GoEditControllerGetParamNormalized(componentPtr unsafe.Pointer, id C.uint32_t) C.double {
	wrapper := getComponent(uintptr(componentPtr))
	if wrapper == nil {
		return 0
	}

	return C.double(wrapper.component.GetParamNormalized(uint32(id)))
}

//export GoEditControllerSetParamNormalized
func GoEditControllerSetParamNormalized(componentPtr unsafe.Pointer, id C.uint32_t, value C.double) C.Steinberg_tresult {
	wrapper := getComponent(uintptr(componentPtr))
	if wrapper == nil {
		return C.Steinberg_tresult(vst3.ResultFalse)
	}

	err := wrapper.component.SetParamNormalized(uint32(id), float64(value))
	if err != nil {
		return C.Steinberg_tresult(vst3.ResultFalse)
	}
	if params := wrapper.component.GetParameters(); params != nil {
		if p, ok := params.GetOK(uint32(id)); ok {
			wrapper.notifyEditorParameterChanged(uint32(id), p.GetNormalized(), p.GetPlain())
		}
	}
	return C.Steinberg_tresult(vst3.ResultOK)
}

//export GoEditControllerSetParamNormalizedWithNotification
func GoEditControllerSetParamNormalizedWithNotification(componentPtr unsafe.Pointer, id C.uint32_t, value C.double) C.Steinberg_tresult {
	wrapper := getComponent(uintptr(componentPtr))
	if wrapper == nil {
		return C.Steinberg_tresult(vst3.ResultFalse)
	}

	err := wrapper.component.SetParamNormalizedWithNotification(uint32(id), float64(value))
	if err != nil {
		return C.Steinberg_tresult(vst3.ResultFalse)
	}
	return C.Steinberg_tresult(vst3.ResultOK)
}

//export GoEditControllerSetComponentHandler
func GoEditControllerSetComponentHandler(componentPtr unsafe.Pointer, handler unsafe.Pointer) C.Steinberg_tresult {
	wrapper := getComponent(uintptr(componentPtr))
	if wrapper == nil {
		return C.Steinberg_tresult(vst3.ResultFalse)
	}

	// Store the component handler for parameter change notifications
	wrapper.handlerMu.Lock()
	wrapper.componentHandler = handler
	wrapper.handlerMu.Unlock()

	return C.Steinberg_tresult(vst3.ResultOK)
}

//export GoEditControllerCreateView
func GoEditControllerCreateView(componentPtr unsafe.Pointer, name *C.char) unsafe.Pointer {
	view := C.VST3GoCreateEditorView(componentPtr)
	if view == nil {
		return nil
	}

	if wrapper := getComponent(uintptr(componentPtr)); wrapper != nil {
		wrapper.setEditorView(unsafe.Pointer(view))
	}

	return unsafe.Pointer(view)
}

//export GoEditControllerClearEditorView
func GoEditControllerClearEditorView(componentPtr unsafe.Pointer, view unsafe.Pointer) {
	wrapper := getComponent(uintptr(componentPtr))
	if wrapper == nil {
		return
	}

	wrapper.clearEditorView(view)
}

//export GoEditControllerGetEditorHTML
func GoEditControllerGetEditorHTML(componentPtr unsafe.Pointer) *C.char {
	wrapper := getComponent(uintptr(componentPtr))
	if wrapper == nil {
		return nil
	}

	snapshot, err := wrapper.component.EditorSnapshot()
	if err != nil {
		return nil
	}

	snapshotJSON, err := json.Marshal(snapshot)
	if err != nil {
		return nil
	}

	encoded := base64.StdEncoding.EncodeToString(snapshotJSON)
	html := buildEditorHTML(encoded)
	return C.CString(html)
}
