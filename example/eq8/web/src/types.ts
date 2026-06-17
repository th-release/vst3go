export interface EditorControl {
  id: number;
  name: string;
  shortName: string;
  unit: string;
  kind: string;
  normalized: number;
  plain: number;
  min: number;
  max: number;
  defaultValue: number;
  stepCount: number;
  flags: number;
  readOnly: boolean;
  hidden: boolean;
}

export interface EditorSection {
  title: string;
  controls: EditorControl[];
}

export interface EditorModel {
  plugin: {
    id: string;
    name: string;
    version: string;
    vendor: string;
    category: string;
  };
  sections: EditorSection[];
}

export interface EditorSnapshot {
  model: EditorModel;
}

export interface EqBand {
  id: number;
  label: string;
  type: number;
  gain: number;
  frequency: number;
  q: number;
}

export interface AppState {
  snapshot: EditorSnapshot;
}
