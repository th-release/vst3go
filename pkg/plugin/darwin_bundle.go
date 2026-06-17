package plugin

import "path/filepath"

// DarwinBundleLayout describes the expected VST3 directory structure for a
// macOS plugin bundle produced from a built dylib.
type DarwinBundleLayout struct {
	BundleDir   string
	ContentsDir string
	MacOSDir    string
	BinaryPath  string
	HeaderPath  string
	PlistPath   string
}

// DarwinBundlePaths returns the canonical macOS VST3 bundle locations for a
// given bundle root and plugin name.
func DarwinBundlePaths(bundleRoot, pluginName string) DarwinBundleLayout {
	if bundleRoot == "" {
		bundleRoot = "dist/macos"
	}
	if pluginName == "" {
		pluginName = "vst3go"
	}

	bundleDir := filepath.Join(bundleRoot, pluginName+".vst3")
	contentsDir := filepath.Join(bundleDir, "Contents")
	macOSDir := filepath.Join(contentsDir, "MacOS")
	return DarwinBundleLayout{
		BundleDir:   bundleDir,
		ContentsDir: contentsDir,
		MacOSDir:    macOSDir,
		BinaryPath:  filepath.Join(macOSDir, pluginName),
		HeaderPath:  filepath.Join(macOSDir, pluginName+".h"),
		PlistPath:   filepath.Join(contentsDir, "Info.plist"),
	}
}
