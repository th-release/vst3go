import { useEffect, useMemo, useState } from "react";
import type { EditorControl, EditorSnapshot, EqBand } from "./types";
import { decodeSnapshot, snapshotStorageKey } from "./snapshot";
import { postParamChange } from "./host";

const BAND_LAYOUT = [
  { index: 1, enable: 100, frequency: 101, gain: 102, q: 103, type: 104 },
  { index: 2, enable: 110, frequency: 111, gain: 112, q: 113, type: 114 },
  { index: 3, enable: 120, frequency: 121, gain: 122, q: 123, type: 124 },
  { index: 4, enable: 130, frequency: 131, gain: 132, q: 133, type: 134 },
  { index: 5, enable: 140, frequency: 141, gain: 142, q: 143, type: 144 },
  { index: 6, enable: 150, frequency: 151, gain: 152, q: 153, type: 154 },
  { index: 7, enable: 160, frequency: 161, gain: 162, q: 163, type: 164 },
  { index: 8, enable: 170, frequency: 171, gain: 172, q: 173, type: 174 },
];

const TYPE_OPTIONS = ["Bell", "Low Shelf", "High Shelf", "Low Cut", "High Cut", "Notch"];

function getControls(snapshot: EditorSnapshot): EditorControl[] {
  return snapshot.model.sections.flatMap((section) => section.controls ?? []);
}

function buildControlMap(snapshot: EditorSnapshot): Map<number, EditorControl> {
  return new Map(getControls(snapshot).map((control) => [control.id, control]));
}

function clamp(value: number, min: number, max: number): number {
  return Math.min(max, Math.max(min, value));
}

function formatFrequency(value: number): string {
  if (value >= 1000) {
    return `${(value / 1000).toFixed(2)} kHz`;
  }
  return `${value.toFixed(1)} Hz`;
}

function formatControlValue(control: EditorControl): string {
  if (control.kind === "choice") {
    return TYPE_OPTIONS[Math.round(control.plain)] ?? `${control.plain}`;
  }
  if (control.unit === "Hz") {
    return formatFrequency(control.plain);
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
  return control.plain.toFixed(3);
}

function normalizedFromPlain(control: EditorControl, plain: number): number {
  if (control.max <= control.min) {
    return 0;
  }
  return clamp((plain - control.min) / (control.max - control.min), 0, 1);
}

function plainFromNormalized(control: EditorControl, normalized: number): number {
  return control.min + clamp(normalized, 0, 1) * (control.max - control.min);
}

function getInitialSnapshot(): EditorSnapshot {
  return decodeSnapshot(window.__EQ8_SNAPSHOT__);
}

function snapshotKey(snapshot: EditorSnapshot): string {
  return snapshotStorageKey(snapshot);
}

function loadSavedSnapshot(key: string): EditorSnapshot | null {
  try {
    const saved = window.localStorage.getItem(key);
    if (!saved) {
      return null;
    }
    return JSON.parse(saved) as EditorSnapshot;
  } catch {
    return null;
  }
}

function saveSnapshot(key: string, snapshot: EditorSnapshot): void {
  window.localStorage.setItem(key, JSON.stringify(snapshot));
}

function controlById(controlMap: Map<number, EditorControl>, id: number): EditorControl {
  const control = controlMap.get(id);
  if (!control) {
    throw new Error(`Missing control ${id}`);
  }
  return control;
}

function updateControl(
  snapshot: EditorSnapshot,
  controlMap: Map<number, EditorControl>,
  id: number,
  normalized: number,
  plain: number,
): EditorSnapshot {
  const nextControls = getControls(snapshot).map((control) => {
    if (control.id !== id) {
      return control;
    }
    return { ...control, normalized, plain };
  });

  const nextSnapshot: EditorSnapshot = {
    model: {
      ...snapshot.model,
      sections: snapshot.model.sections.map((section) => ({
        ...section,
        controls: section.controls.map((control) => {
          const next = nextControls.find((candidate) => candidate.id === control.id);
          return next ?? control;
        }),
      })),
    },
  };

  controlMap.clear();
  for (const control of getControls(nextSnapshot)) {
    controlMap.set(control.id, control);
  }

  return nextSnapshot;
}

function formatBandLabel(index: number): string {
  return `Band ${index}`;
}

export function App() {
  const [snapshot, setSnapshot] = useState<EditorSnapshot>(() => getInitialSnapshot());
  const [status, setStatus] = useState("Live snapshot");

  const controlMap = useMemo(() => buildControlMap(snapshot), [snapshot]);
  const bands = useMemo<EqBand[]>(
    () =>
      BAND_LAYOUT.map((layout) => ({
        id: layout.index,
        label: formatBandLabel(layout.index),
        type: controlById(controlMap, layout.type).plain,
        gain: controlById(controlMap, layout.gain).plain,
        frequency: controlById(controlMap, layout.frequency).plain,
        q: controlById(controlMap, layout.q).plain,
      })),
    [controlMap],
  );

  const inputGain = controlById(controlMap, 1);
  const outputGain = controlById(controlMap, 2);
  const bypass = controlById(controlMap, 3);
  const analyzer = controlById(controlMap, 4);

  useEffect(() => {
    const key = snapshotKey(snapshot);
    try {
      saveSnapshot(key, snapshot);
    } catch {
      setStatus("Snapshot save failed");
    }

    const listener = (event: MessageEvent) => {
      const payload = typeof event.data === "string" ? JSON.parse(event.data) : event.data;
      if (!payload || payload.type !== "param-change") {
        return;
      }

      setSnapshot((current) => {
        const currentMap = buildControlMap(current);
        const control = currentMap.get(Number(payload.id));
        if (!control) {
          return current;
        }
        return updateControl(current, currentMap, control.id, Number(payload.normalized), Number(payload.plain));
      });
    };

    window.addEventListener("message", listener);
    return () => window.removeEventListener("message", listener);
  }, [snapshot]);

  function mutateControl(id: number, normalized: number, plain: number) {
    setSnapshot((current) => {
      const currentMap = buildControlMap(current);
      const next = updateControl(current, currentMap, id, normalized, plain);
      postParamChange({ type: "param-change", id, normalized, plain });
      return next;
    });
    setStatus("Snapshot saved locally");
  }

  function restoreSavedSnapshot() {
    const key = snapshotKey(snapshot);
    const saved = loadSavedSnapshot(key);
    if (!saved) {
      setStatus("No saved snapshot yet");
      return;
    }
    setSnapshot(saved);
    setStatus("Snapshot restored");
  }

  function renderBandCard(layout: (typeof BAND_LAYOUT)[number]) {
    const enable = controlById(controlMap, layout.enable);
    const type = controlById(controlMap, layout.type);
    const frequency = controlById(controlMap, layout.frequency);
    const gain = controlById(controlMap, layout.gain);
    const q = controlById(controlMap, layout.q);

    return (
      <article key={layout.index} className={`band-card${enable.plain < 0.5 ? " is-off" : ""}`}>
        <div className="band-header">
          <div>
            <p className="eyebrow">{layout.index}</p>
            <h3>{formatBandLabel(layout.index)}</h3>
          </div>
          <span className="band-chip">{enable.plain >= 0.5 ? "Live" : "Off"}</span>
        </div>

        <div className="band-body">
          <label>
            <span>Enable</span>
            <div className="toggle-row">
              <input
                type="checkbox"
                checked={enable.plain >= 0.5}
                onChange={(event) => mutateControl(layout.enable, event.target.checked ? 1 : 0, event.target.checked ? 1 : 0)}
              />
              <span className="metric">{formatControlValue(enable)}</span>
            </div>
          </label>

          <label>
            <span>Type</span>
            <select
              value={Math.round(type.plain)}
              onChange={(event) => {
                const plain = Number(event.target.value);
                mutateControl(layout.type, normalizedFromPlain(type, plain), plain);
              }}
            >
              {TYPE_OPTIONS.map((label, index) => (
                <option key={label} value={index}>
                  {label}
                </option>
              ))}
            </select>
            <span className="metric">{formatControlValue(type)}</span>
          </label>

          {[
            [frequency, layout.frequency, "Frequency"],
            [gain, layout.gain, "Gain"],
            [q, layout.q, "Q"],
          ].map(([control, id, label]) => (
            <label key={label as string}>
              <span>{label}</span>
              <input
                type="range"
                min="0"
                max="1"
                step="0.001"
                value={control.normalized}
                onChange={(event) => {
                  const normalized = Number(event.target.value);
                  const plain = plainFromNormalized(control as EditorControl, normalized);
                  mutateControl(id as number, normalized, plain);
                }}
              />
              <input
                type="number"
                value={control.plain}
                min={control.min}
                max={control.max}
                step={label === "Q" ? 0.01 : 0.1}
                onChange={(event) => {
                  const plain = Number(event.target.value);
                  mutateControl(id as number, normalizedFromPlain(control as EditorControl, plain), plain);
                }}
              />
              <span className="metric">{formatControlValue(control as EditorControl)}</span>
            </label>
          ))}

          <div className="band-metrics">
            <span className="metric">{formatFrequency(frequency.plain)}</span>
            <span className="metric">{gain.plain.toFixed(1)} dB</span>
            <span className="metric">Q {q.plain.toFixed(2)}</span>
          </div>
        </div>
      </article>
    );
  }

  function renderGraph() {
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

    const points = BAND_LAYOUT.map((layout) => {
      const frequency = controlById(controlMap, layout.frequency).plain;
      const gain = controlById(controlMap, layout.gain).plain;
      const enabled = controlById(controlMap, layout.enable).plain >= 0.5;
      return {
        index: layout.index,
        enabled,
        x: freqToX(frequency, minFreq, maxFreq, left, right),
        y: gainToY(gain, minGain, maxGain, top, bottom),
      };
    }).sort((a, b) => a.x - b.x);

    return (
      <svg viewBox={`0 0 ${width} ${height}`} className="graph" aria-label="EQ response graph">
        <rect x="0" y="0" width={width} height={height} rx="16" fill="transparent" />
        {Array.from({ length: 5 }, (_, index) => -24 + index * 12).map((gain) => {
          const y = gainToY(gain, minGain, maxGain, top, bottom);
          return (
            <g key={gain}>
              <line x1={left} y1={y} x2={right} y2={y} stroke="rgba(255,255,255,0.08)" strokeWidth="1" />
              <text x="14" y={y + 4} fill="rgba(226,232,240,0.65)" fontSize="11">
                {gain} dB
              </text>
            </g>
          );
        })}
        {[20, 50, 100, 200, 500, 1000, 2000, 5000, 10000, 20000].map((freq) => {
          const x = freqToX(freq, minFreq, maxFreq, left, right);
          return (
            <g key={freq}>
              <line x1={x} y1={top} x2={x} y2={bottom} stroke="rgba(255,255,255,0.06)" strokeWidth="1" />
              <text x={x} y={height - 8} fill="rgba(226,232,240,0.65)" fontSize="11" textAnchor="middle">
                {formatGraphFrequency(freq)}
              </text>
            </g>
          );
        })}

        <path
          d={buildCurvePath(points)}
          fill="none"
          stroke="rgba(167,139,250,0.95)"
          strokeWidth="4"
          strokeLinecap="round"
          strokeLinejoin="round"
        />
        {points.map((point) => (
          <g key={point.index}>
            <circle
              cx={point.x}
              cy={point.y}
              r={point.enabled ? 9 : 6}
              fill={point.enabled ? "rgba(124,58,237,0.98)" : "rgba(148,163,184,0.8)"}
              stroke="rgba(255,255,255,0.9)"
              strokeWidth="2"
            />
            <text x={point.x} y={point.y - 16} fill="rgba(226,232,240,0.9)" fontSize="12" textAnchor="middle">
              B{point.index}
            </text>
          </g>
        ))}
      </svg>
    );
  }

  return (
    <div className="shell">
      <header className="topbar">
        <div className="titleblock">
          <p className="eyebrow">{snapshot.model.plugin.vendor}</p>
          <h1>{snapshot.model.plugin.name}</h1>
          <p className="lede">{snapshot.model.plugin.category} · {snapshot.model.plugin.version}</p>
        </div>
        <div className="editor-actions">
          <div className="pill">{status}</div>
          <button
            type="button"
            onClick={() => {
              saveSnapshot(snapshotKey(snapshot), snapshot);
              setStatus("Snapshot saved locally");
            }}
          >
            Save snapshot
          </button>
          <button type="button" onClick={restoreSavedSnapshot}>Restore snapshot</button>
        </div>
      </header>

      <main className="surface">
        <section className="graph-panel">
          <div className="panel-heading">
            <div>
              <p className="eyebrow">Response</p>
              <h2>Frequency curve</h2>
            </div>
            <div className="panel-tags">
              <span className="pill">{bypass.plain >= 0.5 ? "Bypass on" : "Bypass off"}</span>
              <span className="pill">{analyzer.plain >= 0.5 ? "Analyzer on" : "Analyzer off"}</span>
            </div>
          </div>
          {renderGraph()}
        </section>

        <section className="section">
          <div className="panel-heading">
            <div>
              <p className="eyebrow">Bands</p>
              <h2>Eight EQ bands</h2>
            </div>
          </div>
          <div className="band-grid">{BAND_LAYOUT.map((layout) => renderBandCard(layout))}</div>
        </section>

        <section className="section">
          <div className="panel-heading">
            <div>
              <p className="eyebrow">Global</p>
              <h2>Input, output, and utilities</h2>
            </div>
          </div>
          <div className="global-grid">
            <ControlCard
              control={inputGain}
              onChange={(plain, normalized) => mutateControl(1, normalized, plain)}
            />
            <ControlCard
              control={outputGain}
              onChange={(plain, normalized) => mutateControl(2, normalized, plain)}
            />
            <ToggleCard
              control={bypass}
              onChange={(plain, normalized) => mutateControl(3, normalized, plain)}
            />
            <ToggleCard
              control={analyzer}
              onChange={(plain, normalized) => mutateControl(4, normalized, plain)}
            />
          </div>
        </section>
      </main>
    </div>
  );
}

function ControlCard({
  control,
  onChange,
}: {
  control: EditorControl;
  onChange: (plain: number, normalized: number) => void;
}) {
  return (
    <article className="control-card">
      <h3>{control.name}</h3>
      <div className="control-body">
        <input
          type="range"
          min="0"
          max="1"
          step="0.001"
          value={control.normalized}
          onChange={(event) => {
            const normalized = Number(event.target.value);
            onChange(plainFromNormalized(control, normalized), normalized);
          }}
        />
        <input
          type="number"
          value={control.plain}
          min={control.min}
          max={control.max}
          step="0.01"
          onChange={(event) => {
            const plain = Number(event.target.value);
            onChange(plain, normalizedFromPlain(control, plain));
          }}
        />
        <span className="metric">{formatControlValue(control)}</span>
      </div>
    </article>
  );
}

function ToggleCard({
  control,
  onChange,
}: {
  control: EditorControl;
  onChange: (plain: number, normalized: number) => void;
}) {
  return (
    <article className="control-card">
      <h3>{control.name}</h3>
      <div className="control-body">
        <label className="toggle-row">
          <input
            type="checkbox"
            checked={control.plain >= 0.5}
            onChange={(event) => {
              const plain = event.target.checked ? 1 : 0;
              onChange(plain, plain);
            }}
          />
          <span>Toggle</span>
        </label>
        <span className="metric">{formatControlValue(control)}</span>
      </div>
    </article>
  );
}

function freqToX(freq: number, minFreq: number, maxFreq: number, left: number, right: number): number {
  const minLog = Math.log10(minFreq);
  const maxLog = Math.log10(maxFreq);
  const clamped = clamp(freq, minFreq, maxFreq);
  const ratio = (Math.log10(clamped) - minLog) / (maxLog - minLog);
  return left + ratio * (right - left);
}

function gainToY(gain: number, minGain: number, maxGain: number, top: number, bottom: number): number {
  const clamped = clamp(gain, minGain, maxGain);
  const ratio = (clamped - minGain) / (maxGain - minGain);
  return bottom - ratio * (bottom - top);
}

function buildCurvePath(points: Array<{ x: number; y: number }>): string {
  if (points.length === 0) {
    return "";
  }

  let path = `M ${points[0].x} ${points[0].y}`;
  for (let index = 1; index < points.length; index += 1) {
    const previous = points[index - 1];
    const current = points[index];
    const midX = (previous.x + current.x) / 2;
    path += ` C ${midX} ${previous.y}, ${midX} ${current.y}, ${current.x} ${current.y}`;
  }
  return path;
}

function formatGraphFrequency(freq: number): string {
  if (freq >= 1000) {
    return `${freq / 1000}k`;
  }
  return String(freq);
}
