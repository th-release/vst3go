package plugin

import (
	"path/filepath"
	"testing"
)

func TestWindowsBundlePaths(t *testing.T) {
	layout := WindowsBundlePaths("dist/windows", "vst3go")

	if got, want := layout.BundleDir, filepath.Join("dist/windows", "vst3go.vst3"); got != want {
		t.Fatalf("BundleDir = %q, want %q", got, want)
	}
	if got, want := layout.ContentsDir, filepath.Join("dist/windows", "vst3go.vst3", "Contents"); got != want {
		t.Fatalf("ContentsDir = %q, want %q", got, want)
	}
	if got, want := layout.BinaryDir, filepath.Join("dist/windows", "vst3go.vst3", "Contents", "x86_64-win"); got != want {
		t.Fatalf("BinaryDir = %q, want %q", got, want)
	}
	if got, want := layout.BinaryPath, filepath.Join("dist/windows", "vst3go.vst3", "Contents", "x86_64-win", "vst3go.vst3"); got != want {
		t.Fatalf("BinaryPath = %q, want %q", got, want)
	}
	if got, want := layout.LoaderPath, filepath.Join("dist/windows", "vst3go.vst3", "Contents", "x86_64-win", "WebView2Loader.dll"); got != want {
		t.Fatalf("LoaderPath = %q, want %q", got, want)
	}
}
