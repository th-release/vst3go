package plugin

import (
	"path/filepath"
	"testing"
)

func TestDarwinBundlePaths(t *testing.T) {
	layout := DarwinBundlePaths("dist/macos", "vst3go")

	if got, want := layout.BundleDir, filepath.Join("dist/macos", "vst3go.vst3"); got != want {
		t.Fatalf("BundleDir = %q, want %q", got, want)
	}
	if got, want := layout.ContentsDir, filepath.Join("dist/macos", "vst3go.vst3", "Contents"); got != want {
		t.Fatalf("ContentsDir = %q, want %q", got, want)
	}
	if got, want := layout.MacOSDir, filepath.Join("dist/macos", "vst3go.vst3", "Contents", "MacOS"); got != want {
		t.Fatalf("MacOSDir = %q, want %q", got, want)
	}
	if got, want := layout.BinaryPath, filepath.Join("dist/macos", "vst3go.vst3", "Contents", "MacOS", "vst3go"); got != want {
		t.Fatalf("BinaryPath = %q, want %q", got, want)
	}
	if got, want := layout.HeaderPath, filepath.Join("dist/macos", "vst3go.vst3", "Contents", "MacOS", "vst3go.h"); got != want {
		t.Fatalf("HeaderPath = %q, want %q", got, want)
	}
	if got, want := layout.PlistPath, filepath.Join("dist/macos", "vst3go.vst3", "Contents", "Info.plist"); got != want {
		t.Fatalf("PlistPath = %q, want %q", got, want)
	}
}

func TestDarwinBundlePathsDefaults(t *testing.T) {
	layout := DarwinBundlePaths("", "")

	if got, want := layout.BundleDir, filepath.Join("dist/macos", "vst3go.vst3"); got != want {
		t.Fatalf("BundleDir = %q, want %q", got, want)
	}
	if got, want := layout.BinaryPath, filepath.Join("dist/macos", "vst3go.vst3", "Contents", "MacOS", "vst3go"); got != want {
		t.Fatalf("BinaryPath = %q, want %q", got, want)
	}
}
