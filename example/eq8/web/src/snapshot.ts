import type { EditorSnapshot } from "./types";

declare global {
  interface Window {
    __EQ8_SNAPSHOT__?: string;
    webkit?: {
      messageHandlers?: {
        vst3go?: {
          postMessage(message: unknown): void;
        };
      };
    };
    chrome?: {
      webview?: {
        postMessage(message: unknown): void;
        addEventListener(type: "message", listener: (event: MessageEvent) => void): void;
      };
    };
  }
}

export function decodeSnapshot(encoded: string | undefined): EditorSnapshot {
  if (!encoded) {
    throw new Error("missing editor snapshot");
  }

  const json = atob(encoded);
  return JSON.parse(json) as EditorSnapshot;
}

export function encodeSnapshot(snapshot: EditorSnapshot): string {
  return btoa(JSON.stringify(snapshot));
}

export function snapshotStorageKey(snapshot: EditorSnapshot): string {
  return `vst3go.eq8.snapshot.${snapshot.model.plugin.id || "default"}`;
}
