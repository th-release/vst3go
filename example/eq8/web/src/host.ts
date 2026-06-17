export interface ParamChangeMessage {
  type: "param-change";
  id: number;
  normalized: number;
  plain: number;
}

export function postParamChange(message: ParamChangeMessage): void {
  if (window.webkit?.messageHandlers?.vst3go) {
    window.webkit.messageHandlers.vst3go.postMessage(message);
    return;
  }

  if (window.chrome?.webview) {
    window.chrome.webview.postMessage(message);
  }
}
