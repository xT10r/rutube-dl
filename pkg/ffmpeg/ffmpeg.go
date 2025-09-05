package ffmpeg

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/StanislavKH/rutube-dl/pkg/i18n"
)

const (
	// Primary source - BtbN builds (pre-compiled binaries)
	FFmpegBtbNReleasesAPI = "https://api.github.com/repos/BtbN/FFmpeg-Builds/releases/latest"
	// Alternative source - FFmpeg Static builds
	FFmpegAltReleasesAPI = "https://api.github.com/repos/eugeneware/ffmpeg-static/releases/latest"
	FFmpegTimeout        = 10 * time.Minute
)

type FFmpegManager struct {
	BinaryPath string
	Version    string
	Sources    []FFmpegSource
}

type FFmpegSource struct {
	Name        string
	APIURL      string
	PlatformMap map[string]string
	ExtractPath string
}

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// NewFFmpegManager creates a new FFmpeg manager instance with multiple sources
func NewFFmpegManager() *FFmpegManager {
	return &FFmpegManager{
		Sources: []FFmpegSource{
			{
				Name:   "BtbN-Builds",
				APIURL: FFmpegBtbNReleasesAPI,
				PlatformMap: map[string]string{
					"windows/amd64": "win64",
					"windows/386":   "win32",
					"linux/amd64":   "linux64",
					"linux/arm64":   "linuxarm64",
					"darwin/amd64":  "macos64",
					"darwin/arm64":  "macos64",
				},
				ExtractPath: "/bin/",
			},
			{
				Name:   "FFmpeg-Static",
				APIURL: FFmpegAltReleasesAPI,
				PlatformMap: map[string]string{
					"windows/amd64": "win32-x64",
					"linux/amd64":   "linux-x64",
					"darwin/amd64":  "darwin-x64",
				},
				ExtractPath: "/",
			},
		},
	}
}

// EnsureFFmpeg ensures FFmpeg is available and up to date
func (fm *FFmpegManager) EnsureFFmpeg(forceUpdate bool) error {
	// Get latest version info first to determine proper path
	release, source, err := fm.getLatestVersionInfo()
	if err != nil {
		// If we can't get version info, try to use system FFmpeg as fallback
		fmt.Printf("Failed to get FFmpeg version info: %v\n", err)
		return fm.trySystemFFmpeg()
	}

	// Set binary path based on version
	fm.BinaryPath = fm.getVersionedFFmpegPath(release.TagName)

	if !forceUpdate && fm.isFFmpegAvailable() {
		// Check if it's the current version
		if fm.isCurrentVersion(release.TagName) {
			fmt.Printf("FFmpeg %s is already available at: %s\n", release.TagName, fm.BinaryPath)
			fmt.Println(i18n.T(i18n.MsgFFmpegReady))
			return nil
		}
	}

	// Try to download or update FFmpeg
	if fm.isFFmpegAvailable() {
		fmt.Printf("Updating FFmpeg to version %s\n", release.TagName)
	} else {
		fmt.Printf("Downloading FFmpeg version %s\n", release.TagName)
	}

	err = fm.downloadFFmpegVersion(release, source)
	if err != nil {
		fmt.Printf("Failed to download FFmpeg: %v\n", err)
		// Try to use system FFmpeg as fallback
		return fm.trySystemFFmpeg()
	}

	return nil
}

// trySystemFFmpeg attempts to use system-installed FFmpeg as fallback
func (fm *FFmpegManager) trySystemFFmpeg() error {
	fmt.Println("Trying to use system FFmpeg as fallback...")

	// Try to find FFmpeg in PATH
	systemPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return fmt.Errorf("FFmpeg download failed and no system FFmpeg found: %v", err)
	}

	// Test if system FFmpeg works
	cmd := exec.Command(systemPath, "-version")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("system FFmpeg found but not working: %v", err)
	}

	fm.BinaryPath = systemPath
	fmt.Printf("Using system FFmpeg: %s\n", systemPath)
	fmt.Println(i18n.T(i18n.MsgFFmpegReady))
	return nil
}

// GetBinaryPath returns the path to the FFmpeg binary
func (fm *FFmpegManager) GetBinaryPath() string {
	return fm.BinaryPath
}

// isFFmpegAvailable checks if FFmpeg binary exists and is executable
func (fm *FFmpegManager) isFFmpegAvailable() bool {
	if fm.BinaryPath == "" {
		return false
	}

	if _, err := os.Stat(fm.BinaryPath); os.IsNotExist(err) {
		return false
	}

	// Test if FFmpeg is executable
	cmd := exec.Command(fm.BinaryPath, "-version")
	err := cmd.Run()
	return err == nil
}

// getLocalFFmpegPath returns the expected path for local FFmpeg binary (legacy method)
func (fm *FFmpegManager) getLocalFFmpegPath() string {
	execDir, err := os.Executable()
	if err != nil {
		execDir = "."
	} else {
		execDir = filepath.Dir(execDir)
	}

	ffmpegDir := filepath.Join(execDir, "ffmpeg")

	if runtime.GOOS == "windows" {
		return filepath.Join(ffmpegDir, "ffmpeg.exe")
	}
	return filepath.Join(ffmpegDir, "ffmpeg")
}

// getVersionedFFmpegPath returns the path for version-specific FFmpeg binary
func (fm *FFmpegManager) getVersionedFFmpegPath(version string) string {
	execDir, err := os.Executable()
	if err != nil {
		execDir = "."
	} else {
		execDir = filepath.Dir(execDir)
	}

	// Create version-specific directory: bin/ffmpeg/[version]
	ffmpegDir := filepath.Join(execDir, "bin", "ffmpeg", version)

	if runtime.GOOS == "windows" {
		return filepath.Join(ffmpegDir, "ffmpeg.exe")
	}
	return filepath.Join(ffmpegDir, "ffmpeg")
}

// isCurrentVersion checks if the specified version is already available
func (fm *FFmpegManager) isCurrentVersion(version string) bool {
	if fm.BinaryPath == "" {
		return false
	}

	// Check if the version-specific directory exists and has the binary
	if _, err := os.Stat(fm.BinaryPath); os.IsNotExist(err) {
		return false
	}

	// Test if FFmpeg is executable
	cmd := exec.Command(fm.BinaryPath, "-version")
	err := cmd.Run()
	if err != nil {
		return false
	}

	// Check if version file matches
	versionFile := filepath.Join(filepath.Dir(fm.BinaryPath), "version.txt")
	data, err := os.ReadFile(versionFile)
	if err != nil {
		return false
	}

	savedVersion := strings.TrimSpace(string(data))
	return strings.Contains(savedVersion, version)
}

// getLatestVersionInfo gets the latest FFmpeg version info from GitHub with fallback sources
func (fm *FFmpegManager) getLatestVersionInfo() (*GitHubRelease, *FFmpegSource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for i, source := range fm.Sources {
		fmt.Printf("Trying source %d/%d: %s\n", i+1, len(fm.Sources), source.Name)

		req, err := http.NewRequestWithContext(ctx, "GET", source.APIURL, nil)
		if err != nil {
			fmt.Printf("Error creating request for %s: %v\n", source.Name, err)
			continue
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("Error fetching from %s: %v\n", source.Name, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("%s API returned status: %d\n", source.Name, resp.StatusCode)
			continue
		}

		var release GitHubRelease
		err = json.NewDecoder(resp.Body).Decode(&release)
		if err != nil {
			fmt.Printf("Error decoding response from %s: %v\n", source.Name, err)
			continue
		}

		// Check if this source has builds for our platform
		platformKey := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
		if _, hasKey := source.PlatformMap[platformKey]; hasKey {
			fmt.Printf("Found compatible release from %s: %s\n", source.Name, release.TagName)
			return &release, &source, nil
		}

		fmt.Printf("%s doesn't support platform %s\n", source.Name, platformKey)
	}

	return nil, nil, fmt.Errorf("no compatible FFmpeg source found for platform %s/%s", runtime.GOOS, runtime.GOARCH)
}

// isLatestVersion checks if the current FFmpeg version is the latest
func (fm *FFmpegManager) isLatestVersion() bool {
	if fm.BinaryPath == "" || !fm.isFFmpegAvailable() {
		return false
	}

	// Get current version
	currentVersion := fm.getCurrentVersion()
	if currentVersion == "" {
		return false
	}

	// Get latest version
	release, _, err := fm.getLatestVersionInfo()
	if err != nil {
		// If we can't check latest version, assume current is fine
		return true
	}

	return currentVersion == release.TagName
}

// getCurrentVersion gets the current FFmpeg version
func (fm *FFmpegManager) getCurrentVersion() string {
	if fm.Version != "" {
		return fm.Version
	}

	versionFile := filepath.Join(filepath.Dir(fm.BinaryPath), "version.txt")
	data, err := os.ReadFile(versionFile)
	if err != nil {
		return ""
	}

	fm.Version = strings.TrimSpace(string(data))
	return fm.Version
}

// GetAvailableVersions returns a list of locally available FFmpeg versions
func (fm *FFmpegManager) GetAvailableVersions() []string {
	execDir, err := os.Executable()
	if err != nil {
		execDir = "."
	} else {
		execDir = filepath.Dir(execDir)
	}

	ffmpegBaseDir := filepath.Join(execDir, "bin", "ffmpeg")
	if _, err := os.Stat(ffmpegBaseDir); os.IsNotExist(err) {
		return []string{}
	}

	entries, err := os.ReadDir(ffmpegBaseDir)
	if err != nil {
		return []string{}
	}

	var versions []string
	for _, entry := range entries {
		if entry.IsDir() {
			// Check if this directory contains a valid FFmpeg binary
			version := entry.Name()
			binaryPath := fm.getVersionedFFmpegPath(version)
			if _, err := os.Stat(binaryPath); err == nil {
				versions = append(versions, version)
			}
		}
	}

	return versions
}

// downloadFFmpeg downloads and sets up FFmpeg (legacy method)
func (fm *FFmpegManager) downloadFFmpeg() error {
	release, source, err := fm.getLatestVersionInfo()
	if err != nil {
		return fmt.Errorf("failed to get latest version info: %v", err)
	}

	return fm.downloadFFmpegVersion(release, source)
}

// downloadFFmpegVersion downloads and sets up a specific version of FFmpeg
func (fm *FFmpegManager) downloadFFmpegVersion(release *GitHubRelease, source *FFmpegSource) error {
	// Find the appropriate asset for the current platform
	downloadURL := fm.getDownloadURLForPlatform(release, source)
	if downloadURL == "" {
		return fmt.Errorf("no suitable FFmpeg build found for platform: %s/%s from source: %s", runtime.GOOS, runtime.GOARCH, source.Name)
	}

	fmt.Printf("Downloading FFmpeg %s from %s: %s\n", release.TagName, source.Name, downloadURL)

	// Create version-specific FFmpeg directory
	ffmpegDir := filepath.Dir(fm.BinaryPath)
	err := os.MkdirAll(ffmpegDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create FFmpeg directory: %v", err)
	}

	// Download and extract
	err = fm.downloadAndExtract(downloadURL, ffmpegDir, source)
	if err != nil {
		return fmt.Errorf("failed to download and extract FFmpeg: %v", err)
	}

	// Save version info
	versionFile := filepath.Join(ffmpegDir, "version.txt")
	versionInfo := fmt.Sprintf("%s (from %s)", release.TagName, source.Name)
	err = os.WriteFile(versionFile, []byte(versionInfo), 0644)
	if err != nil {
		return fmt.Errorf("failed to save version info: %v", err)
	}

	fm.Version = release.TagName
	fmt.Printf("FFmpeg %s successfully installed at: %s\n", release.TagName, fm.BinaryPath)
	fmt.Println(i18n.T(i18n.MsgFFmpegReady))
	return nil
}

// getDownloadURLForPlatform finds the appropriate download URL for the current platform
func (fm *FFmpegManager) getDownloadURLForPlatform(release *GitHubRelease, source *FFmpegSource) string {
	platformKey := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	platformPattern, exists := source.PlatformMap[platformKey]
	if !exists {
		return ""
	}

	// First, let's log all available assets for debugging
	fmt.Printf("Available assets for %s:\n", source.Name)
	for i, asset := range release.Assets {
		fmt.Printf("  %d: %s\n", i+1, asset.Name)
	}
	fmt.Println()

	// Try multiple pattern variations for better compatibility
	patterns := []string{
		platformPattern + "-static",                     // win64-gpl-static
		platformPattern,                                 // win64-gpl
		strings.Replace(platformPattern, "-gpl", "", 1), // win64
	}

	// Add more flexible patterns for Windows
	if runtime.GOOS == "windows" {
		patterns = append(patterns,
			"win64",
			"windows-64",
			"windows",
			"win-64",
		)
	}

	// Add more flexible patterns for Linux
	if runtime.GOOS == "linux" {
		patterns = append(patterns,
			"linux64",
			"linux-64",
			"linux",
		)
	}

	// Add more flexible patterns for macOS
	if runtime.GOOS == "darwin" {
		patterns = append(patterns,
			"macos64",
			"macos",
			"osx",
			"darwin",
		)
	}

	// Try each pattern
	for _, pattern := range patterns {
		for _, asset := range release.Assets {
			if strings.Contains(strings.ToLower(asset.Name), strings.ToLower(pattern)) &&
				(strings.HasSuffix(asset.Name, ".zip") ||
					strings.HasSuffix(asset.Name, ".tar.gz") ||
					strings.HasSuffix(asset.Name, ".tar.xz")) {
				fmt.Printf("Found build with pattern '%s': %s\n", pattern, asset.Name)
				return asset.BrowserDownloadURL
			}
		}
	}

	fmt.Printf("No compatible build found for platform %s with any pattern\n", platformKey)
	return ""
}

// downloadAndExtract downloads and extracts FFmpeg
func (fm *FFmpegManager) downloadAndExtract(url, targetDir string, source *FFmpegSource) error {
	ctx, cancel := context.WithTimeout(context.Background(), FFmpegTimeout)
	defer cancel()

	// Download the file
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating download request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error downloading FFmpeg: %v", err)
	}
	defer resp.Body.Close()

	// Determine file extension
	var fileExt string
	if strings.HasSuffix(url, ".zip") {
		fileExt = ".zip"
	} else if strings.HasSuffix(url, ".tar.gz") {
		fileExt = ".tar.gz"
	} else if strings.HasSuffix(url, ".tar.xz") {
		fileExt = ".tar.xz"
	} else {
		fileExt = ".zip" // default
	}

	// Create temporary file
	tempFile, err := os.CreateTemp("", "ffmpeg_*"+fileExt)
	if err != nil {
		return fmt.Errorf("error creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Copy downloaded content to temp file
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return fmt.Errorf("error saving downloaded file: %v", err)
	}

	// Extract the archive
	if fileExt == ".zip" {
		return fm.extractFFmpegFromZip(tempFile.Name(), targetDir, source)
	} else {
		return fmt.Errorf("unsupported archive format: %s", fileExt)
	}
}

// extractFFmpegFromZip extracts the FFmpeg binary and required DLLs from the downloaded zip
func (fm *FFmpegManager) extractFFmpegFromZip(zipPath, targetDir string, source *FFmpegSource) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("error opening zip file: %v", err)
	}
	defer reader.Close()

	// Find the FFmpeg binary in the zip
	var ffmpegFile *zip.File
	var dllFiles []*zip.File // For Windows shared builds
	ffmpegName := "ffmpeg"
	if runtime.GOOS == "windows" {
		ffmpegName = "ffmpeg.exe"
	}

	// Try different patterns based on the source
	patterns := []string{
		source.ExtractPath + ffmpegName,
		"/bin/" + ffmpegName,
		"\\bin\\" + ffmpegName,
		"/" + ffmpegName,
		"\\" + ffmpegName,
		ffmpegName,
	}

	for _, file := range reader.File {
		// Look for FFmpeg binary
		for _, pattern := range patterns {
			if strings.HasSuffix(file.Name, pattern) {
				ffmpegFile = file
				break
			}
		}

		// For Windows, also collect DLL files from bin directory
		if runtime.GOOS == "windows" && strings.Contains(file.Name, "/bin/") && strings.HasSuffix(file.Name, ".dll") {
			dllFiles = append(dllFiles, file)
		}
	}

	if ffmpegFile == nil {
		// List all files for debugging
		fmt.Println("Available files in archive:")
		for _, file := range reader.File {
			fmt.Printf("  %s\n", file.Name)
		}
		return fmt.Errorf("FFmpeg binary not found in the downloaded archive")
	}

	fmt.Printf("Found FFmpeg binary: %s\n", ffmpegFile.Name)
	if len(dllFiles) > 0 {
		fmt.Printf("Found %d DLL files for Windows\n", len(dllFiles))
	}

	// Extract the FFmpeg binary
	err = fm.extractFile(ffmpegFile, fm.BinaryPath)
	if err != nil {
		return fmt.Errorf("error extracting FFmpeg binary: %v", err)
	}

	// Extract DLL files to the same directory as FFmpeg (for Windows shared builds)
	if runtime.GOOS == "windows" && len(dllFiles) > 0 {
		ffmpegDir := filepath.Dir(fm.BinaryPath)
		for _, dllFile := range dllFiles {
			dllName := filepath.Base(dllFile.Name)
			dllPath := filepath.Join(ffmpegDir, dllName)
			err = fm.extractFile(dllFile, dllPath)
			if err != nil {
				fmt.Printf("Warning: Failed to extract DLL %s: %v\n", dllName, err)
				// Continue with other DLLs, don't fail completely
				continue
			}
			fmt.Printf("Extracted DLL: %s\n", dllName)
		}
	}

	// Make the binary executable on Unix-like systems
	if runtime.GOOS != "windows" {
		err = os.Chmod(fm.BinaryPath, 0755)
		if err != nil {
			return fmt.Errorf("error making FFmpeg executable: %v", err)
		}
	}

	return nil
}

// extractFile extracts a single file from zip to the target path
func (fm *FFmpegManager) extractFile(zipFile *zip.File, targetPath string) error {
	src, err := zipFile.Open()
	if err != nil {
		return fmt.Errorf("error opening file in archive: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("error creating target file: %v", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("error copying file: %v", err)
	}

	return nil
}
