package ffmpeg

import (
	"runtime"
	"strings"
	"testing"
)

func TestNewFFmpegManager(t *testing.T) {
	mgr := NewFFmpegManager()

	if mgr == nil {
		t.Fatal("NewFFmpegManager() returned nil")
	}

	if len(mgr.Sources) == 0 {
		t.Fatal("FFmpegManager should have at least one source")
	}

	// Verify BtbN-Builds source exists
	foundBtbN := false
	for _, source := range mgr.Sources {
		if source.Name == "BtbN-Builds" {
			foundBtbN = true

			// Test platform support
			platformKey := runtime.GOOS + "/" + runtime.GOARCH
			if _, supported := source.PlatformMap[platformKey]; !supported {
				t.Errorf("BtbN-Builds should support current platform: %s", platformKey)
			}
			break
		}
	}

	if !foundBtbN {
		t.Error("BtbN-Builds source not found in FFmpegManager")
	}
}

func TestFFmpegPlatformDetection(t *testing.T) {
	mgr := NewFFmpegManager()

	testCases := []struct {
		goos     string
		goarch   string
		expected string
	}{
		{"windows", "amd64", "win64"},
		{"windows", "386", "win32"},
		{"linux", "amd64", "linux64"},
		{"linux", "arm64", "linuxarm64"},
		{"darwin", "amd64", "macos64"},
		{"darwin", "arm64", "macos64"},
	}

	for _, tc := range testCases {
		t.Run(tc.goos+"_"+tc.goarch, func(t *testing.T) {
			// Find BtbN-Builds source
			var btbnSource *FFmpegSource
			for i, source := range mgr.Sources {
				if source.Name == "BtbN-Builds" {
					btbnSource = &mgr.Sources[i]
					break
				}
			}

			if btbnSource == nil {
				t.Fatal("BtbN-Builds source not found")
			}

			platformKey := tc.goos + "/" + tc.goarch
			pattern, exists := btbnSource.PlatformMap[platformKey]

			if !exists {
				t.Errorf("Platform %s not supported", platformKey)
				return
			}

			if pattern != tc.expected {
				t.Errorf("Expected pattern %s for %s, got %s", tc.expected, platformKey, pattern)
			}
		})
	}
}

func TestFFmpegBinaryPath(t *testing.T) {
	mgr := NewFFmpegManager()
	mgr.BinaryPath = mgr.getLocalFFmpegPath()

	path := mgr.GetBinaryPath()
	if path == "" {
		t.Error("GetBinaryPath() should not return empty string")
	}

	// Check that path ends with correct executable name
	expectedSuffix := "ffmpeg"
	if runtime.GOOS == "windows" {
		expectedSuffix = "ffmpeg.exe"
	}

	if !strings.HasSuffix(path, expectedSuffix) {
		t.Errorf("Binary path should end with %s, got: %s", expectedSuffix, path)
	}
}

// Integration test that actually tries to get release information
func TestFFmpegReleaseInfo(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	mgr := NewFFmpegManager()

	release, source, err := mgr.getLatestVersionInfo()
	if err != nil {
		t.Logf("Warning: Could not get release info (network issue?): %v", err)
		return // Don't fail the test for network issues
	}

	if release == nil {
		t.Error("Release should not be nil when no error occurred")
		return
	}

	if source == nil {
		t.Error("Source should not be nil when no error occurred")
		return
	}

	if release.TagName == "" {
		t.Error("Release tag should not be empty")
	}

	if len(release.Assets) == 0 {
		t.Error("Release should have at least one asset")
	}

	t.Logf("Found release: %s from %s with %d assets",
		release.TagName, source.Name, len(release.Assets))
}

// Test download URL generation without actually downloading
func TestFFmpegDownloadURL(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	mgr := NewFFmpegManager()

	release, source, err := mgr.getLatestVersionInfo()
	if err != nil {
		t.Logf("Warning: Could not get release info (network issue?): %v", err)
		return
	}

	downloadURL := mgr.getDownloadURLForPlatform(release, source)
	if downloadURL == "" {
		t.Error("Should find a download URL for current platform")

		// Log available assets for debugging
		t.Logf("Available assets for %s:", source.Name)
		for i, asset := range release.Assets {
			t.Logf("  %d: %s", i+1, asset.Name)
		}
		return
	}

	if !strings.HasPrefix(downloadURL, "http") {
		t.Errorf("Download URL should be a valid HTTP URL, got: %s", downloadURL)
	}

	t.Logf("Found download URL: %s", downloadURL)
}

// Benchmark for performance testing
func BenchmarkNewFFmpegManager(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mgr := NewFFmpegManager()
		_ = mgr
	}
}

func BenchmarkPlatformDetection(b *testing.B) {
	mgr := NewFFmpegManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		platformKey := runtime.GOOS + "/" + runtime.GOARCH
		for _, source := range mgr.Sources {
			_, _ = source.PlatformMap[platformKey]
		}
	}
}
