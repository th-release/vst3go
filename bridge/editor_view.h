#ifndef VST3GO_EDITOR_VIEW_H
#define VST3GO_EDITOR_VIEW_H

#include "../include/vst3/vst3_c_api.h"

#ifdef __cplusplus
extern "C" {
#endif

typedef struct Steinberg_ViewRect {
    Steinberg_int32 left;
    Steinberg_int32 top;
    Steinberg_int32 right;
    Steinberg_int32 bottom;
} Steinberg_ViewRect;

typedef struct Steinberg_IPlugFrame Steinberg_IPlugFrame;

typedef struct Steinberg_IPlugViewVtbl {
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *queryInterface)(void* thisInterface, const Steinberg_TUID iid, void** obj);
    Steinberg_uint32 (SMTG_STDMETHODCALLTYPE *addRef)(void* thisInterface);
    Steinberg_uint32 (SMTG_STDMETHODCALLTYPE *release)(void* thisInterface);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *isPlatformTypeSupported)(void* thisInterface, Steinberg_FIDString type);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *attached)(void* thisInterface, void* parent, Steinberg_FIDString type);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *removed)(void* thisInterface);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *onWheel)(void* thisInterface, float distance);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *onKeyDown)(void* thisInterface, Steinberg_Vst_TChar key, Steinberg_int32 keyCode, Steinberg_int16 modifiers);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *onKeyUp)(void* thisInterface, Steinberg_Vst_TChar key, Steinberg_int32 keyCode, Steinberg_int16 modifiers);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *getSize)(void* thisInterface, Steinberg_ViewRect* size);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *onSize)(void* thisInterface, Steinberg_ViewRect* size);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *onFocus)(void* thisInterface, Steinberg_TBool state);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *setFrame)(void* thisInterface, Steinberg_IPlugFrame* frame);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *canResize)(void* thisInterface);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *checkSizeConstraint)(void* thisInterface, Steinberg_ViewRect* rect);
    Steinberg_tresult (SMTG_STDMETHODCALLTYPE *onSizeScaleFactorChanged)(void* thisInterface, float factor);
} Steinberg_IPlugViewVtbl;

void* VST3GoCreateEditorView(void* component);
void VST3GoEditorViewUpdateParameter(void* view, Steinberg_Vst_ParamID id, Steinberg_Vst_ParamValue normalized, Steinberg_Vst_ParamValue plain);
extern char* GoEditControllerGetEditorHTML(void* component);
extern void GoEditControllerClearEditorView(void* component, void* view);
extern Steinberg_tresult GoEditControllerSetParamNormalizedWithNotification(void* component, Steinberg_Vst_ParamID id, Steinberg_Vst_ParamValue value);

#ifdef __cplusplus
}
#endif

#endif
