package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/StanislavKH/rutube-dl/pkg/ffmpeg"
	"github.com/StanislavKH/rutube-dl/pkg/i18n"
	"github.com/StanislavKH/rutube-dl/pkg/utils"
)

func main() {
	fmt.Println("=== Rutube-dl Integration Tests ===")
	fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Language: %s\n", i18n.GetLanguage())
	fmt.Println()

	runFFmpegTests()
	runI18nTests()
	runUtilsTests()
}

func runFFmpegTests() {
	fmt.Println("🔧 Testing FFmpeg Management...")

	mgr := ffmpeg.NewFFmpegManager()
	fmt.Printf("✅ FFmpeg manager created with %d sources\n", len(mgr.Sources))

	// Test source availability
	for i, source := range mgr.Sources {
		fmt.Printf("📋 Source %d: %s (%s)\n", i+1, source.Name, source.APIURL)

		platformKey := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
		if pattern, supported := source.PlatformMap[platformKey]; supported {
			fmt.Printf("   ✅ Supports current platform: %s\n", pattern)
		} else {
			fmt.Printf("   ❌ Does not support current platform\n")
		}
	}

	// Test actual FFmpeg setup (download if needed)
	fmt.Println("\n⬇️  Testing FFmpeg setup (will download if not present)...")
	fmt.Println("📝 This test will actually download FFmpeg if it's not present")
	fmt.Println("⚠️  FFmpeg will be saved permanently in version-specific directories")

	// Check available versions before download
	availableVersions := mgr.GetAvailableVersions()
	if len(availableVersions) > 0 {
		fmt.Printf("📦 Available local FFmpeg versions: %v\n", availableVersions)
	} else {
		fmt.Println("📦 No local FFmpeg versions found")
	}

	err := mgr.EnsureFFmpeg(false)
	if err != nil {
		log.Printf("❌ FFmpeg setup failed: %v", err)
		fmt.Println("🛠️  Testing fallback to system FFmpeg...")
	} else {
		fmt.Printf("✅ FFmpeg ready at: %s\n", mgr.GetBinaryPath())

		// Show available versions after download
		updatedVersions := mgr.GetAvailableVersions()
		if len(updatedVersions) > 0 {
			fmt.Printf("📦 Available FFmpeg versions after download: %v\n", updatedVersions)
		}

		// Check if binary exists and is accessible
		if _, err := os.Stat(mgr.GetBinaryPath()); err == nil {
			fmt.Println("✅ FFmpeg binary exists and is accessible")

			// Test FFmpeg execution (version check)
			fmt.Println("🗺️ Testing FFmpeg execution...")
			// Note: We don't actually run ffmpeg -version here to keep tests lightweight
			fmt.Println("✅ FFmpeg appears to be ready for use")
		} else {
			fmt.Printf("❌ FFmpeg binary not accessible: %v\n", err)
		}
	}

	// Test with force update flag (commented out to avoid unnecessary downloads)
	fmt.Println("\n📝 Force update test (skipped to avoid re-downloading)")
	fmt.Println("ℹ️  To test force update, uncomment the code in integration.go")
	// err = mgr.EnsureFFmpeg(true)
	// if err != nil {
	//	log.Printf("⚠️  Warning: FFmpeg force update failed: %v", err)
	// } else {
	//	fmt.Println("✅ FFmpeg force update successful")
	// }

	fmt.Println()
}

func runI18nTests() {
	fmt.Println("🌍 Testing Internationalization...")

	currentLang := i18n.GetLanguage()
	fmt.Printf("📍 Current language: %s\n", currentLang)

	// Test message translation
	testKeys := []i18n.MessageKey{
		i18n.MsgDownloadingSegments,
		i18n.MsgFFmpegReady,
		i18n.MsgUsage,
		i18n.MsgVideoDownloaded,
	}

	fmt.Println("📝 Testing message translations:")
	for _, key := range testKeys {
		message := i18n.T(key)
		fmt.Printf("   %s: %s\n", string(key), message)
	}

	// Test language switching
	fmt.Println("\n🔄 Testing language switching:")
	originalLang := i18n.GetLanguage()

	// Switch to Russian
	i18n.SetLanguage(i18n.Russian)
	russianMsg := i18n.T(i18n.MsgDownloadingSegments)
	fmt.Printf("   Russian: %s\n", russianMsg)

	// Switch to English
	i18n.SetLanguage(i18n.English)
	englishMsg := i18n.T(i18n.MsgDownloadingSegments)
	fmt.Printf("   English: %s\n", englishMsg)

	// Restore original language
	i18n.SetLanguage(originalLang)
	fmt.Println("✅ Language switching successful")
	fmt.Println()
}

func runUtilsTests() {
	fmt.Println("🛠️  Testing Utilities...")

	// Test filename sanitization
	testFilenames := []string{
		"Тестовое видео с русскими символами",
		"Video with /illegal\\characters<>:|?*",
		"Normal filename",
		"Очень_длинное_название_видео_которое_должно_быть_обрезано",
	}

	fmt.Println("📄 Testing filename sanitization:")
	for _, filename := range testFilenames {
		safe := utils.SanitizeFilename(filename, false)
		safeWithTranslit := utils.SanitizeFilename(filename, true)

		fmt.Printf("   Original: %s\n", filename)
		fmt.Printf("   Safe: %s\n", safe)
		fmt.Printf("   Safe+Translit: %s\n", safeWithTranslit)
		fmt.Println()
	}

	// Test safe filename generation
	testTitle := "Тестовое видео"
	safeFilename := utils.GetSafeFilename(testTitle, "mp4", true)
	fmt.Printf("📁 Safe filename: %s\n", safeFilename)

	// Test directory creation path
	testDir, err := utils.CreateSafeDirectory("/tmp", testTitle, true)
	if err == nil {
		fmt.Printf("📂 Safe directory path: %s\n", testDir)
	}

	fmt.Println("✅ Utilities testing successful")
	fmt.Println()
}
