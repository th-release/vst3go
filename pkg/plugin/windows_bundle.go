package plugin

import "path/filepath"

// WindowsBundleLayout describes the expected VST3 directory structure for a
// Windows plugin bundle produced from a built DLL.
type WindowsBundleLayout struct {
	BundleDir   string
	ContentsDir string
	BinaryDir   string
	BinaryPath  string
	HeaderPath  string
	LoaderPath  string
}

// WindowsBundlePaths returns the canonical Windows VST3 bundle locations for a
// given bundle root and plugin name.
func WindowsBundlePaths(bundleRoot, pluginName string) WindowsBundleLayout {
	if bundleRoot == "" {
		bundleRoot = "dist/windows"
	}
	if pluginName == "" {
		pluginName = "vst3go"
	}

	bundleDir := filepath.Join(bundleRoot, pluginName+".vst3")
	contentsDir := filepath.Join(bundleDir, "Contents")
	binaryDir := filepath.Join(contentsDir, "x86_64-win")
	return WindowsBundleLayout{
		BundleDir:   bundleDir,
		ContentsDir: contentsDir,
		BinaryDir:   binaryDir,
		BinaryPath:  filepath.Join(binaryDir, pluginName+".vst3"),
		HeaderPath:  filepath.Join(binaryDir, pluginName+".h"),
		LoaderPath:  filepath.Join(binaryDir, "WebView2Loader.dll"),
	}
}
