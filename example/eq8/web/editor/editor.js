const snapshot = JSON.parse(atob("__EQ8_SNAPSHOT__"));
const model = snapshot.model;

const root = {
  pluginVendor: document.getElementById("plugin-vendor"),
  pluginName: document.getElementById("plugin-name"),
  pluginMeta: document.getElementById("plugin-meta"),
  snapshotStatus: document.getElementById("snapshot-status"),
  saveSnapshot: document.getElementById("save-snapshot"),
  restoreSnapshot: document.getElementById("restore-snapshot"),
  bandGrid: document.getElementById("band-grid"),
  globalGrid: document.getElementById("global-grid"),
  graph: document.getElementById("response-graph"),
  bypassPill: document.getElementById("bypass-pill"),
  analyzerPill: document.getElementById("analyzer-pill"),
};

const host = window.webkit?.messageHandlers?.vst3go ?? window.chrome?.webview ?? null;
const storageKey = "vst3go.eq8.snapshot." + (model.plugin.id || "default");

const bandLayout = [
  { index: 1, enable: 100, frequency: 101, gain: 102, q: 103, type: 104 },
  { index: 2, enable: 110, frequency: 111, gain: 112, q: 113, type: 114 },
  { index: 3, enable: 120, frequency: 121, gain: 122, q: 123, type: 124 },
  { index: 4, enable: 130, frequency: 131, gain: 132, q: 133, type: 134 },
  { index: 5, enable: 140, frequency: 141, gain: 142, q: 143, type: 144 },
  { index: 6, enable: 150, frequency: 151, gain: 152, q: 153, type: 154 },
  { index: 7, enable: 160, frequency: 161, gain: 162, q: 163, type: 164 },
  { index: 8, enable: 170, frequency: 171, gain: 172, q: 173, type: 174 },
];

const typeOptions = ["Bell", "Low Shelf", "High Shelf", "Low Cut", "High Cut", "Notch"];
const controlBindings = new Map();
const controlIndex = new Map();

function sections() {
  return Array.isArray(model.sections) ? model.sections : [];
}

function controls() {
  const all = [];
  for (const section of sections()) {
    if (Array.isArray(section.controls)) {
      all.push(...section.controls);
    }
  }
  return all;
}

function findControl(id) {
  return controlIndex.get(id) || null;
}

function getControl(id) {
  const control = findControl(id);
  if (!control) {
    throw new Error(`Missing control ${id}`);
  }
  return control;
}

function clamp(value, min, max) {
  return Math.min(max, Math.max(min, value));
}

function setStatus(message) {
  if (root.snapshotStatus) {
    root.snapshotStatus.textContent = message;
  }
}

function sendChange(id, normalized, plain) {
  if (host && typeof host.postMessage === "function") {
    host.postMessage({ type: "param-change", id, normalized, plain });
  }
}

function loadSavedSnapshot() {
  try {
    const saved = window.localStorage.getItem(storageKey);
    if (!saved) {
      return null;
    }
    return JSON.parse(saved);
  } catch (error) {
    return null;
  }
}

function saveSnapshot() {
  try {
    window.localStorage.setItem(storageKey, JSON.stringify(snapshot));
    setStatus("Snapshot saved locally");
    if (root.restoreSnapshot) {
      root.restoreSnapshot.disabled = false;
    }
  } catch (error) {
    setStatus("Snapshot save failed");
  }
}

function restoreSnapshot(nextSnapshot, notifyHost) {
  if (!nextSnapshot || !nextSnapshot.model || !Array.isArray(nextSnapshot.model.sections)) {
    return false;
  }

  for (const section of nextSnapshot.model.sections) {
    if (!Array.isArray(section.controls)) {
      continue;
    }
    for (const control of section.controls) {
      const existing = findControl(control.id);
      if (!existing) {
        continue;
      }
      existing.normalized = control.normalized;
      existing.plain = control.plain;
      applyControlToBinding(existing);
      if (notifyHost) {
        sendChange(existing.id, existing.normalized, existing.plain);
      }
    }
  }

  renderGraph();
  setStatus(notifyHost ? "Snapshot restored" : "Snapshot loaded");
  return true;
}

function updateControl(id, normalized, plain, notifyHost) {
  const control = getControl(id);
  control.normalized = normalized;
  control.plain = plain;
  applyControlToBinding(control);
  updateStatusPills();
  renderGraph();
  if (notifyHost) {
    sendChange(id, normalized, plain);
  }
}

function plainFromNormalized(control, normalized) {
  return control.min + clamp(normalized, 0, 1) * (control.max - control.min);
}

function normalizedFromPlain(control, plain) {
  if (control.max <= control.min) {
    return 0;
  }
  return clamp((plain - control.min) / (control.max - control.min), 0, 1);
}

function applyControlToBinding(control) {
  const binding = controlBindings.get(control.id);
  if (!binding) {
    return;
  }

  if (binding.toggle) {
    binding.toggle.checked = control.plain >= 0.5;
  }
  if (binding.select) {
    binding.select.value = String(control.plain);
  }
  if (binding.range) {
    binding.range.value = String(control.normalized);
  }
  if (binding.number) {
    binding.number.value = String(control.plain);
  }
  if (binding.readout) {
    binding.readout.textContent = formatControlValue(control);
  }
  if (binding.card && binding.affectsCardState) {
    binding.card.classList.toggle("is-off", control.plain < 0.5);
  }
}

function formatControlValue(control) {
  if (control.kind === "choice") {
    return typeOptions[Math.round(control.plain)] || String(control.plain);
  }
  if (control.unit === "Hz") {
    if (control.plain >= 1000) {
      return `${(control.plain / 1000).toFixed(2)} kHz`;
    }
    return `${control.plain.toFixed(1)} Hz`;
  }
  if (control.unit === "dB") {
    return `${control.plain.toFixed(1)} dB`;
  }
  if (control.unit === "%") {
    return `${control.plain.toFixed(0)}%`;
  }
  if (control.kind === "toggle") {
    return control.plain >= 0.5 ? "On" : "Off";
  }
  return Number(control.plain).toFixed(3);
}

function createElement(tag, className, text) {
  const element = document.createElement(tag);
  if (className) {
    element.className = className;
  }
  if (text !== undefined) {
    element.textContent = text;
  }
  return element;
}

function createBandCard(layout) {
  const card = createElement("article", "band-card");
  const title = createElement("div", "band-header");
  const heading = createElement("div");
  heading.appendChild(createElement("p", "eyebrow", `Band ${layout.index}`));
  heading.appendChild(createElement("h3", "", "Eight-band cell"));
  const chip = createElement("span", "band-chip", "Live");
  title.appendChild(heading);
  title.appendChild(chip);
  card.appendChild(title);

  const body = createElement("div", "band-body");
  const enable = makeToggleField(getControl(layout.enable), "Enable");
  const type = makeSelectField(getControl(layout.type), "Type", typeOptions);
  const frequency = makeSliderField(getControl(layout.frequency), "Frequency");
  const gain = makeSliderField(getControl(layout.gain), "Gain");
  const q = makeSliderField(getControl(layout.q), "Q");

  body.appendChild(enable.wrapper);
  body.appendChild(type.wrapper);
  body.appendChild(frequency.wrapper);
  body.appendChild(gain.wrapper);
  body.appendChild(q.wrapper);

  const metrics = createElement("div", "band-metrics");
  const freqMetric = createElement("span", "metric");
  const gainMetric = createElement("span", "metric");
  const qMetric = createElement("span", "metric");
  metrics.appendChild(freqMetric);
  metrics.appendChild(gainMetric);
  metrics.appendChild(qMetric);
  body.appendChild(metrics);

  card.appendChild(body);

  controlBindings.set(layout.enable, { toggle: enable.input, card, readout: enable.readout, affectsCardState: true });
  controlBindings.set(layout.type, { select: type.input, card, readout: type.readout });
  controlBindings.set(layout.frequency, { range: frequency.range, number: frequency.number, card, readout: frequency.readout });
  controlBindings.set(layout.gain, { range: gain.range, number: gain.number, card, readout: gain.readout });
  controlBindings.set(layout.q, { range: q.range, number: q.number, card, readout: q.readout });

  bindMetricUpdate(layout, {
    card,
    chip,
    freqMetric,
    gainMetric,
    qMetric,
  });

  return card;
}

function makeToggleField(control, label) {
  const wrapper = createElement("label");
  wrapper.appendChild(createElement("span", "", label));
  const row = createElement("div", "toggle-row");
  const input = document.createElement("input");
  input.type = "checkbox";
  input.checked = control.plain >= 0.5;
  const readout = createElement("span", "metric", formatControlValue(control));
  row.appendChild(input);
  row.appendChild(readout);
  wrapper.appendChild(row);

  input.addEventListener("change", () => {
    const plain = input.checked ? 1 : 0;
    updateControl(control.id, plain, plain, true);
    saveSnapshot();
  });

  return { wrapper, input, readout };
}

function makeSelectField(control, label, options) {
  const wrapper = createElement("label");
  wrapper.appendChild(createElement("span", "", label));
  const select = document.createElement("select");
  for (let index = 0; index < options.length; index += 1) {
    const option = document.createElement("option");
    option.value = String(index);
    option.textContent = options[index];
    select.appendChild(option);
  }
  select.value = String(Math.round(control.plain));
  const readout = createElement("span", "metric", formatControlValue(control));
  wrapper.appendChild(select);
  wrapper.appendChild(readout);

  select.addEventListener("change", () => {
    const plain = Number(select.value);
    const normalized = normalizedFromPlain(control, plain);
    updateControl(control.id, normalized, plain, true);
    saveSnapshot();
  });

  return { wrapper, input: select, readout };
}

function makeSliderField(control, label) {
  const wrapper = createElement("label");
  wrapper.appendChild(createElement("span", "", label));
  const range = document.createElement("input");
  range.type = "range";
  range.min = "0";
  range.max = "1";
  range.step = control.kind === "toggle" ? "1" : "0.001";
  range.value = String(control.normalized);
  const number = document.createElement("input");
  number.type = "number";
  number.min = String(control.min);
  number.max = String(control.max);
  number.step = control.kind === "toggle" ? "1" : "0.01";
  number.value = String(control.plain);
  const readout = createElement("span", "metric", formatControlValue(control));
  const rangeRow = createElement("div", "control-body");
  rangeRow.appendChild(range);
  rangeRow.appendChild(number);
  wrapper.appendChild(rangeRow);
  wrapper.appendChild(readout);

  range.addEventListener("input", () => {
    const normalized = Number(range.value);
    const plain = plainFromNormalized(control, normalized);
    updateControl(control.id, normalized, plain, true);
    saveSnapshot();
  });

  number.addEventListener("change", () => {
    const plain = Number(number.value);
    const normalized = normalizedFromPlain(control, plain);
    updateControl(control.id, normalized, plain, true);
    saveSnapshot();
  });

  return { wrapper, range, number, readout };
}

function bindMetricUpdate(layout, refs) {
  refs.card.dataset.band = String(layout.index);
  refs.chip.textContent = "Live";
}

function renderGraph() {
  const svg = root.graph;
  while (svg.firstChild) {
    svg.removeChild(svg.firstChild);
  }

  const width = 1000;
  const height = 360;
  const left = 56;
  const right = width - 28;
  const top = 24;
  const bottom = height - 34;
  const minFreq = 20;
  const maxFreq = 20000;
  const minGain = -24;
  const maxGain = 24;

  const background = document.createElementNS("http://www.w3.org/2000/svg", "rect");
  background.setAttribute("x", "0");
  background.setAttribute("y", "0");
  background.setAttribute("width", String(width));
  background.setAttribute("height", String(height));
  background.setAttribute("rx", "16");
  background.setAttribute("fill", "transparent");
  svg.appendChild(background);

  for (let gain = -24; gain <= 24; gain += 12) {
    const y = gainToY(gain, minGain, maxGain, top, bottom);
    const line = svgLine(left, y, right, y, "rgba(255,255,255,0.08)");
    svg.appendChild(line);
    const label = svgText(14, y + 4, `${gain} dB`, "rgba(226,232,240,0.65)");
    svg.appendChild(label);
  }

  const freqMarks = [20, 50, 100, 200, 500, 1000, 2000, 5000, 10000, 20000];
  for (const freq of freqMarks) {
    const x = freqToX(freq, minFreq, maxFreq, left, right);
    const line = svgLine(x, top, x, bottom, "rgba(255,255,255,0.06)");
    svg.appendChild(line);
    const label = svgText(x, height - 8, formatFrequencyLabel(freq), "rgba(226,232,240,0.65)", "middle");
    svg.appendChild(label);
  }

  const bandPoints = [];
  for (const layout of bandLayout) {
    const frequency = getControl(layout.frequency).plain;
    const gain = getControl(layout.gain).plain;
    const enabled = getControl(layout.enable).plain >= 0.5;
    const x = freqToX(frequency, minFreq, maxFreq, left, right);
    const y = gainToY(gain, minGain, maxGain, top, bottom);
    bandPoints.push({ x, y, enabled, index: layout.index, frequency, gain });
  }

  const path = document.createElementNS("http://www.w3.org/2000/svg", "path");
  path.setAttribute("d", buildCurvePath(bandPoints));
  path.setAttribute("fill", "none");
  path.setAttribute("stroke", "rgba(167,139,250,0.95)");
  path.setAttribute("stroke-width", "4");
  path.setAttribute("stroke-linecap", "round");
  path.setAttribute("stroke-linejoin", "round");
  svg.appendChild(path);

  for (const point of bandPoints) {
    const node = document.createElementNS("http://www.w3.org/2000/svg", "circle");
    node.setAttribute("cx", String(point.x));
    node.setAttribute("cy", String(point.y));
    node.setAttribute("r", point.enabled ? "9" : "6");
    node.setAttribute("fill", point.enabled ? "rgba(124,58,237,0.98)" : "rgba(148,163,184,0.8)");
    node.setAttribute("stroke", "rgba(255,255,255,0.9)");
    node.setAttribute("stroke-width", "2");
    svg.appendChild(node);

    const tag = svgText(point.x, point.y - 16, `B${point.index}`, "rgba(226,232,240,0.9)", "middle");
    tag.setAttribute("font-size", "12");
    svg.appendChild(tag);
  }
}

function buildCurvePath(points) {
  if (!points.length) {
    return "";
  }
  const ordered = [...points].sort((a, b) => a.x - b.x);
  let path = `M ${ordered[0].x} ${ordered[0].y}`;
  for (let index = 1; index < ordered.length; index += 1) {
    const prev = ordered[index - 1];
    const current = ordered[index];
    const midX = (prev.x + current.x) / 2;
    path += ` C ${midX} ${prev.y}, ${midX} ${current.y}, ${current.x} ${current.y}`;
  }
  return path;
}

function svgLine(x1, y1, x2, y2, stroke) {
  const line = document.createElementNS("http://www.w3.org/2000/svg", "line");
  line.setAttribute("x1", String(x1));
  line.setAttribute("y1", String(y1));
  line.setAttribute("x2", String(x2));
  line.setAttribute("y2", String(y2));
  line.setAttribute("stroke", stroke);
  line.setAttribute("stroke-width", "1");
  return line;
}

function svgText(x, y, text, fill, anchor = "start") {
  const label = document.createElementNS("http://www.w3.org/2000/svg", "text");
  label.setAttribute("x", String(x));
  label.setAttribute("y", String(y));
  label.setAttribute("fill", fill);
  label.setAttribute("font-size", "11");
  label.setAttribute("text-anchor", anchor);
  label.textContent = text;
  return label;
}

function freqToX(freq, minFreq, maxFreq, left, right) {
  const minLog = Math.log10(minFreq);
  const maxLog = Math.log10(maxFreq);
  const clamped = clamp(freq, minFreq, maxFreq);
  const ratio = (Math.log10(clamped) - minLog) / (maxLog - minLog);
  return left + ratio * (right - left);
}

function gainToY(gain, minGain, maxGain, top, bottom) {
  const clamped = clamp(gain, minGain, maxGain);
  const ratio = (clamped - minGain) / (maxGain - minGain);
  return bottom - ratio * (bottom - top);
}

function formatFrequencyLabel(freq) {
  if (freq >= 1000) {
    return `${freq / 1000}k`;
  }
  return String(freq);
}

for (const control of controls()) {
  controlIndex.set(control.id, control);
}

function renderGlobals() {
  const items = [
    { id: 1, label: "Input Gain" },
    { id: 2, label: "Output Gain" },
    { id: 3, label: "Bypass" },
    { id: 4, label: "Analyzer" },
  ];

  root.globalGrid.innerHTML = "";

  for (const item of items) {
    const control = getControl(item.id);
    const card = createElement("article", "control-card");
    const heading = createElement("h3", "", item.label);
    const body = createElement("div", "control-body");
    const readout = createElement("span", "metric", formatControlValue(control));

    if (control.kind === "toggle") {
      const row = createElement("label", "toggle-row");
      const input = document.createElement("input");
      input.type = "checkbox";
      input.checked = control.plain >= 0.5;
      const caption = createElement("span", "", "Toggle");
      row.appendChild(input);
      row.appendChild(caption);
      body.appendChild(row);
      body.appendChild(readout);

      controlBindings.set(control.id, { toggle: input, card, readout });
      input.addEventListener("change", () => {
        const plain = input.checked ? 1 : 0;
        updateControl(control.id, plain, plain, true);
        saveSnapshot();
      });
    } else {
      const range = document.createElement("input");
      range.type = "range";
      range.min = "0";
      range.max = "1";
      range.step = "0.001";
      range.value = String(control.normalized);
      const number = document.createElement("input");
      number.type = "number";
      number.min = String(control.min);
      number.max = String(control.max);
      number.step = "0.01";
      number.value = String(control.plain);
      body.appendChild(range);
      body.appendChild(number);
      body.appendChild(readout);

      controlBindings.set(control.id, { range, number, card, readout });

      range.addEventListener("input", () => {
        const normalized = Number(range.value);
        const plain = plainFromNormalized(control, normalized);
        updateControl(control.id, normalized, plain, true);
        saveSnapshot();
      });

      number.addEventListener("change", () => {
        const plain = Number(number.value);
        const normalized = normalizedFromPlain(control, plain);
        updateControl(control.id, normalized, plain, true);
        saveSnapshot();
      });
    }

    card.appendChild(heading);
    card.appendChild(body);
    root.globalGrid.appendChild(card);
  }
}

function updateStatusPills() {
  const bypass = getControl(3).plain >= 0.5;
  const analyzer = getControl(4).plain >= 0.5;
  root.bypassPill.textContent = bypass ? "Bypass on" : "Bypass off";
  root.bypassPill.style.color = bypass ? "var(--warn)" : "var(--muted)";
  root.analyzerPill.textContent = analyzer ? "Analyzer on" : "Analyzer off";
  root.analyzerPill.style.color = analyzer ? "var(--good)" : "var(--muted)";
}

function renderBands() {
  root.bandGrid.innerHTML = "";
  for (const layout of bandLayout) {
    root.bandGrid.appendChild(createBandCard(layout));
  }
}

function renderModelMeta() {
  root.pluginVendor.textContent = model.plugin.vendor || "Example Audio";
  root.pluginName.textContent = model.plugin.name || "EQ8 Example";
  root.pluginMeta.textContent = [model.plugin.category, model.plugin.version].filter(Boolean).join(" · ");
}

function attachHostListener() {
  if (window.chrome?.webview) {
    window.chrome.webview.addEventListener("message", (event) => {
      const payload = typeof event.data === "string" ? JSON.parse(event.data) : event.data;
      if (!payload || payload.type !== "param-change") {
        return;
      }
      updateControl(Number(payload.id), Number(payload.normalized), Number(payload.plain), false);
      saveSnapshot();
    });
    return;
  }

  if (window.webkit?.messageHandlers?.vst3go) {
    window.addEventListener("message", (event) => {
      const payload = event.data;
      if (!payload || payload.type !== "param-change") {
        return;
      }
      updateControl(Number(payload.id), Number(payload.normalized), Number(payload.plain), false);
      saveSnapshot();
    });
  }
}

renderModelMeta();
renderBands();
renderGlobals();
updateStatusPills();
renderGraph();
attachHostListener();

if (root.saveSnapshot) {
  root.saveSnapshot.addEventListener("click", saveSnapshot);
}

if (root.restoreSnapshot) {
  root.restoreSnapshot.addEventListener("click", () => {
    const saved = loadSavedSnapshot();
    if (!saved) {
      setStatus("No saved snapshot yet");
      return;
    }
    restoreSnapshot(saved, true);
  });

  root.restoreSnapshot.disabled = !loadSavedSnapshot();
}

setStatus("Editor ready");
