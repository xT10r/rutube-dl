package utils

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestSanitizeFilename(t *testing.T) {
	testCases := []struct {
		name              string
		input             string
		transliterate     bool
		expectContains    []string
		expectNotContains []string
	}{
		{
			name:           "Russian characters without transliteration",
			input:          "Русское видео",
			transliterate:  false,
			expectContains: []string{"Русское", "видео"},
		},
		{
			name:              "Russian characters with transliteration",
			input:             "Русское видео",
			transliterate:     true,
			expectContains:    []string{"Russkoe", "video"},
			expectNotContains: []string{"Русское"},
		},
		{
			name:              "Forbidden characters removal",
			input:             "Video<>:|?*with/forbidden\\chars",
			transliterate:     false,
			expectContains:    []string{"Video", "with", "forbidden", "chars"},
			expectNotContains: []string{"<", ">", ":", "|", "?", "*", "/", "\\"},
		},
		{
			name:           "Empty string handling",
			input:          "",
			transliterate:  false,
			expectContains: []string{"untitled"},
		},
		{
			name:           "Control characters removal",
			input:          "Video\x00with\x1fcontrol\x7fchars",
			transliterate:  false,
			expectContains: []string{"Video", "with", "control", "chars"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := SanitizeFilename(tc.input, tc.transliterate)

			for _, expected := range tc.expectContains {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected result to contain '%s', got: %s", expected, result)
				}
			}

			for _, notExpected := range tc.expectNotContains {
				if strings.Contains(result, notExpected) {
					t.Errorf("Expected result to NOT contain '%s', got: %s", notExpected, result)
				}
			}

			// Ensure result is not empty
			if result == "" {
				t.Error("Result should not be empty")
			}
		})
	}
}

func TestTransliterateText(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"Привет мир", "Privet mir"},
		{"Русский текст", "Russkiy tekst"},
		{"ЗАГЛАВНЫЕ буквы", "ZAGLAVNYE bukvy"},
		{"Mixed русский and English", "Mixed russkiy and English"},
		{"Цифры 123 и символы !@#", "Tsifry 123 i simvoly !@#"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := TransliterateText(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestGetSafeFilename(t *testing.T) {
	testCases := []struct {
		title         string
		extension     string
		transliterate bool
		expectSuffix  string
	}{
		{"Test Video", "mp4", false, ".mp4"},
		{"Русское видео", "avi", true, ".avi"},
		{"Test", ".mov", false, ".mov"}, // Extension with dot
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			result := GetSafeFilename(tc.title, tc.extension, tc.transliterate)

			if !strings.HasSuffix(result, tc.expectSuffix) {
				t.Errorf("Expected result to end with %s, got: %s", tc.expectSuffix, result)
			}

			// Should not be just the extension
			if result == tc.expectSuffix {
				t.Error("Result should not be just the extension")
			}
		})
	}
}

func TestCreateSafeDirectory(t *testing.T) {
	var basePath string
	if runtime.GOOS == "windows" {
		basePath = "C:\\temp"
	} else {
		basePath = "/tmp"
	}

	dirName := "Тестовая папка"

	result, err := CreateSafeDirectory(basePath, dirName, true)
	if err != nil {
		t.Errorf("CreateSafeDirectory returned error: %v", err)
	}

	// Normalize paths for comparison
	normalizedBasePath := filepath.ToSlash(basePath)
	normalizedResult := filepath.ToSlash(result)

	if !strings.HasPrefix(normalizedResult, normalizedBasePath) {
		t.Errorf("Result should start with base path %s, got: %s", normalizedBasePath, normalizedResult)
	}

	// Should contain transliterated version
	if !strings.Contains(result, "Testovaya") {
		t.Errorf("Result should contain transliterated directory name, got: %s", result)
	}
}

func TestForbiddenCharsForCurrentOS(t *testing.T) {
	// Test that forbidden characters are correctly identified for current OS
	forbidden := getForbiddenCharsForOS()

	switch runtime.GOOS {
	case "windows":
		expectedChars := []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*"}
		for _, char := range expectedChars {
			if !strings.Contains(forbidden, char) {
				t.Errorf("Windows forbidden chars should contain '%s'", char)
			}
		}
	case "darwin":
		if !strings.Contains(forbidden, "/") || !strings.Contains(forbidden, ":") {
			t.Error("macOS forbidden chars should contain '/' and ':'")
		}
	default: // Linux and others
		if !strings.Contains(forbidden, "/") {
			t.Error("Unix forbidden chars should contain '/'")
		}
	}
}

func BenchmarkSanitizeFilename(b *testing.B) {
	input := "Сложное русское название с /запрещёнными\\символами<>:|?*"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SanitizeFilename(input, true)
	}
}

func BenchmarkTransliterateText(b *testing.B) {
	input := "Длинный русский текст для тестирования производительности транслитерации"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TransliterateText(input)
	}
}
