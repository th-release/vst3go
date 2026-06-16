package plugin

// #cgo CFLAGS: -I../../include
// #include "../../include/vst3/vst3_c_api.h"
// #include "../../bridge/editor_view.h"
import "C"

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"unsafe"

	"github.com/cwbudde/vst3go/pkg/vst3"
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

func buildEditorHTML(encodedSnapshot string) string {
	return fmt.Sprintf(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <style>
    :root { color-scheme: dark; --bg: #0b1020; --panel: #11192f; --border: rgba(148,163,184,.18); --text: #e2e8f0; --muted: #94a3b8; --accent: #7c3aed; }
    * { box-sizing: border-box; }
    body { margin: 0; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; background: var(--bg); color: var(--text); }
    .shell { padding: 18px; display: grid; gap: 16px; }
    .topbar, .section, .card { border: 1px solid var(--border); border-radius: 16px; background: rgba(17,25,47,.92); }
    .topbar { padding: 18px; display: flex; justify-content: space-between; gap: 16px; align-items: center; }
    .eyebrow { margin: 0 0 6px; text-transform: uppercase; letter-spacing: .12em; color: var(--muted); font-size: .72rem; }
    h1,h2,h3,p { margin-top: 0; }
    h1 { margin-bottom: 6px; font-size: 1.7rem; }
    .lede, .card p { color: var(--muted); line-height: 1.5; }
    .section { padding: 16px; }
    .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 12px; }
    .card { padding: 14px; }
    .card__label { margin-bottom: 6px; color: var(--muted); text-transform: uppercase; letter-spacing: .12em; font-size: .68rem; }
    .control { display: grid; gap: 8px; }
    .control__meta { display: flex; justify-content: space-between; gap: 10px; font-size: .88rem; color: var(--muted); }
    input[type=range] { width: 100%%; }
    input[type=number], select { width: 100%%; padding: 10px 12px; border-radius: 10px; border: 1px solid var(--border); background: #0f172a; color: var(--text); }
    .pill { display: inline-flex; padding: 8px 10px; border-radius: 999px; border: 1px solid var(--border); color: var(--muted); }
    .editor-actions { display: flex; flex-wrap: wrap; gap: 8px; justify-content: flex-end; }
    .editor-actions button { border: 1px solid var(--border); background: #0f172a; color: var(--text); border-radius: 10px; padding: 9px 12px; cursor: pointer; }
    .editor-actions button:disabled { cursor: not-allowed; opacity: 0.55; }
  </style>
</head>
<body>
  <div class="shell">
    <header class="topbar">
      <div>
        <p class="eyebrow" id="plugin-vendor">VST3Go</p>
        <h1 id="plugin-name">Plugin Editor</h1>
        <p class="lede" id="plugin-meta">Loading editor model…</p>
      </div>
      <div class="editor-actions">
        <div class="pill" id="snapshot-status">Live snapshot</div>
        <button id="save-snapshot" type="button">Save snapshot</button>
        <button id="restore-snapshot" type="button">Restore snapshot</button>
      </div>
    </header>
    <section class="section">
      <div class="grid" id="sections"></div>
    </section>
  </div>
  <script>
    const snapshot = JSON.parse(atob("%s"));
    const model = snapshot.model;
    const sectionsRoot = document.getElementById("sections");
    const statusText = document.getElementById("snapshot-status");
    const saveButton = document.getElementById("save-snapshot");
    const restoreButton = document.getElementById("restore-snapshot");
    const storageKey = "vst3go.editor.snapshot." + (model.plugin.id || "default");
    document.getElementById("plugin-vendor").textContent = model.plugin.vendor || "VST3Go";
    document.getElementById("plugin-name").textContent = model.plugin.name || "Plugin Editor";
    document.getElementById("plugin-meta").textContent = [model.plugin.category, model.plugin.version].filter(Boolean).join(" · ");

    function sendChange(id, value) {
      if (window.webkit && window.webkit.messageHandlers && window.webkit.messageHandlers.vst3go) {
        window.webkit.messageHandlers.vst3go.postMessage({ type: "param-change", id, value });
      }
    }

    const controlBindings = new Map();
    const controlIndex = new Map();

    function setStatus(message) {
      if (statusText) {
        statusText.textContent = message;
      }
    }

    function persistSnapshot() {
      try {
        window.localStorage.setItem(storageKey, JSON.stringify(snapshot));
        if (restoreButton) {
          restoreButton.disabled = false;
        }
        setStatus("Snapshot saved locally");
      } catch (error) {
        setStatus("Snapshot save failed");
      }
    }

    function loadSavedSnapshot() {
      try {
        const saved = window.localStorage.getItem(storageKey);
        if (!saved) {
          setStatus("No saved snapshot yet");
          return null;
        }

        const parsed = JSON.parse(saved);
        if (!parsed || !parsed.model || !Array.isArray(parsed.model.sections)) {
          setStatus("Saved snapshot is invalid");
          return null;
        }

        return parsed;
      } catch (error) {
        setStatus("Snapshot restore failed");
        return null;
      }
    }

    function updateControl(control, normalized, plain, notifyGo) {
      const binding = controlBindings.get(control.id);
      if (!binding) {
        return;
      }

      control.normalized = normalized;
      control.plain = plain;

      if (binding.select) {
        const selected = binding.steps > 0 ? Math.round(normalized * binding.steps) / binding.steps : normalized;
        binding.select.value = String(selected);
      }
      if (binding.range) {
        binding.range.value = String(normalized);
      }
      if (binding.value) {
        binding.value.value = String(plain);
      }
      if (binding.readout) {
        binding.readout.textContent = plain.toFixed(3);
      }

      if (notifyGo) {
        sendChange(control.id, normalized);
      }
    }

    function findControl(controlId) {
      return controlIndex.get(controlId) || null;
    }

    function applySnapshot(nextSnapshot, notifyGo) {
      if (!nextSnapshot || !nextSnapshot.model || !Array.isArray(nextSnapshot.model.sections)) {
        return false;
      }

      nextSnapshot.model.sections.forEach((section) => {
        section.controls.forEach((control) => {
          const existing = findControl(control.id);
          if (!existing) {
            return;
          }

          updateControl(existing, control.normalized, control.plain, notifyGo);
        });
      });

      setStatus(notifyGo ? "Snapshot restored" : "Snapshot loaded");
      return true;
    }

    window.__vst3goUpdateParameter = function(id, normalized, plain) {
      const control = findControl(id);
      if (!control) {
        return;
      }

      updateControl(control, normalized, plain, false);
      persistSnapshot();
    };

    function restoreFromLocalStorage() {
      const saved = loadSavedSnapshot();
      if (!saved) {
        return;
      }

      applySnapshot(saved, true);
    }

    function renderControl(control) {
      const card = document.createElement("article");
      card.className = "card";
      card.innerHTML = "<p class=\"card__label\">" + (control.shortName || control.name) + "</p>" +
        "<h3>" + control.name + "</h3>" +
        "<p>" + (control.unit || control.kind) + "</p>";

      const field = document.createElement("div");
      field.className = "control";

      if (control.kind === "choice") {
        const select = document.createElement("select");
        const steps = Math.max(control.stepCount || 1, 1);
        for (let index = 0; index <= steps; index += 1) {
          const option = document.createElement("option");
          option.value = String(index / steps);
          option.textContent = control.stepCount > 0 ? String(index) : (index / steps).toFixed(2);
          select.appendChild(option);
        }
        const selected = steps > 0 ? Math.round(control.normalized * steps) / steps : control.normalized;
        select.value = String(selected);
        select.addEventListener("change", () => {
          const normalized = Number(select.value);
          const plain = control.min + normalized * (control.max - control.min);
          updateControl(control, normalized, plain, true);
          persistSnapshot();
        });
        field.appendChild(select);
        controlBindings.set(control.id, { select, readout: null, steps });
      } else {
        const range = document.createElement("input");
        range.type = "range";
        range.min = "0";
        range.max = "1";
        range.step = control.kind === "toggle" ? "1" : "0.001";
        range.value = String(control.normalized);
        const value = document.createElement("input");
        value.type = "number";
        value.min = String(control.min);
        value.max = String(control.max);
        value.step = control.kind === "toggle" ? "1" : "0.001";
        value.value = String(control.plain);

        function updateFromNormalized(normalized) {
          const clamped = Math.max(0, Math.min(1, normalized));
          const plain = control.min + clamped * (control.max - control.min);
          updateControl(control, clamped, plain, true);
          persistSnapshot();
        }

        range.addEventListener("input", () => updateFromNormalized(Number(range.value)));
        value.addEventListener("change", () => {
          const plain = Number(value.value);
          const normalized = control.max > control.min ? (plain - control.min) / (control.max - control.min) : 0;
          updateFromNormalized(normalized);
        });

        field.appendChild(range);
        field.appendChild(value);
        controlBindings.set(control.id, { range, value, readout: null });
      }

        const readout = document.createElement("div");
        readout.className = "control__meta";
        readout.innerHTML = "<span>" + control.kind + "</span><span>" + Number(control.plain).toFixed(3) + "</span>";
        card.appendChild(field);
        card.appendChild(readout);

      controlBindings.get(control.id).readout = readout.querySelector('span:last-child');
      controlIndex.set(control.id, control);
      return card;
    }

    model.sections.forEach((section) => {
      const wrapper = document.createElement("section");
      wrapper.className = "card";
      const heading = document.createElement("h2");
      heading.textContent = section.title;
      wrapper.appendChild(heading);
      const grid = document.createElement("div");
      grid.className = "grid";
      section.controls.forEach((control) => grid.appendChild(renderControl(control)));
      wrapper.appendChild(grid);
      sectionsRoot.appendChild(wrapper);
    });

    if (saveButton) {
      saveButton.addEventListener("click", () => persistSnapshot());
    }

    if (restoreButton) {
      restoreButton.addEventListener("click", () => restoreFromLocalStorage());
      let hasSavedSnapshot = false;
      try {
        hasSavedSnapshot = !!window.localStorage.getItem(storageKey);
      } catch (error) {
        hasSavedSnapshot = false;
      }
      restoreButton.disabled = !hasSavedSnapshot;
    }
  </script>
</body>
</html>`, encodedSnapshot)
}
