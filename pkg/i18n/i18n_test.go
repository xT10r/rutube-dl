package i18n

import (
	"os"
	"strings"
	"testing"
)

func TestLanguageDetection(t *testing.T) {
	// Save original environment
	originalLang := os.Getenv("LANG")
	originalLCAll := os.Getenv("LC_ALL")
	originalLCMessages := os.Getenv("LC_MESSAGES")

	defer func() {
		// Restore original environment
		os.Setenv("LANG", originalLang)
		os.Setenv("LC_ALL", originalLCAll)
		os.Setenv("LC_MESSAGES", originalLCMessages)
	}()

	testCases := []struct {
		name     string
		envVars  map[string]string
		expected Language
	}{
		{
			name:     "Russian locale",
			envVars:  map[string]string{"LANG": "ru_RU.UTF-8"},
			expected: Russian,
		},
		{
			name:     "Russian variant",
			envVars:  map[string]string{"LC_ALL": "ru_RU"},
			expected: Russian,
		},
		{
			name:     "English locale",
			envVars:  map[string]string{"LANG": "en_US.UTF-8"},
			expected: English,
		},
		{
			name:     "Unknown locale defaults to English",
			envVars:  map[string]string{"LANG": "fr_FR.UTF-8"},
			expected: English,
		},
		{
			name:     "Empty environment defaults to English",
			envVars:  map[string]string{},
			expected: English,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear environment
			os.Unsetenv("LANG")
			os.Unsetenv("LC_ALL")
			os.Unsetenv("LC_MESSAGES")
			os.Unsetenv("LANGUAGE")

			// Set test environment
			for key, value := range tc.envVars {
				os.Setenv(key, value)
			}

			detected := detectSystemLanguage()
			if detected != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, detected)
			}
		})
	}
}

func TestMessageTranslation(t *testing.T) {
	testCases := []struct {
		key     MessageKey
		lang    Language
		hasArgs bool
		args    []interface{}
	}{
		{MsgDownloadingSegments, English, false, nil},
		{MsgDownloadingSegments, Russian, false, nil},
		{MsgDownloadingList, English, true, []interface{}{"test-id"}},
		{MsgDownloadingList, Russian, true, []interface{}{"test-id"}},
		{MsgUsage, English, false, nil},
		{MsgUsage, Russian, false, nil},
	}

	for _, tc := range testCases {
		t.Run(string(tc.key)+"_"+string(tc.lang), func(t *testing.T) {
			// Set language
			originalLang := GetLanguage()
			SetLanguage(tc.lang)
			defer SetLanguage(originalLang)

			var result string
			if tc.hasArgs {
				result = T(tc.key, tc.args...)
			} else {
				result = T(tc.key)
			}

			if result == "" {
				t.Error("Translation should not be empty")
			}

			if result == string(tc.key) {
				t.Error("Translation should not return the key itself")
			}

			// For Russian, check that result contains Cyrillic characters
			if tc.lang == Russian {
				hasCyrillic := false
				for _, r := range result {
					if r >= 0x0400 && r <= 0x04FF { // Cyrillic block
						hasCyrillic = true
						break
					}
				}
				if !hasCyrillic && tc.key != MsgExample { // Some keys might not have Cyrillic
					t.Log("Warning: Russian translation doesn't contain Cyrillic characters:", result)
				}
			}
		})
	}
}

func TestLanguageSwitching(t *testing.T) {
	originalLang := GetLanguage()
	defer SetLanguage(originalLang)

	// Test switching to Russian
	SetLanguage(Russian)
	if GetLanguage() != Russian {
		t.Error("Language should be Russian after SetLanguage(Russian)")
	}

	russianMsg := T(MsgDownloadingSegments)

	// Test switching to English
	SetLanguage(English)
	if GetLanguage() != English {
		t.Error("Language should be English after SetLanguage(English)")
	}

	englishMsg := T(MsgDownloadingSegments)

	// Messages should be different
	if russianMsg == englishMsg {
		t.Error("Russian and English messages should be different")
	}
}

func TestTfAlias(t *testing.T) {
	// Test that Tf is an alias for T
	originalLang := GetLanguage()
	SetLanguage(English)
	defer SetLanguage(originalLang)

	tResult := T(MsgDownloadingSegments)
	tfResult := Tf(MsgDownloadingSegments)

	if tResult != tfResult {
		t.Error("T() and Tf() should return the same result")
	}
}

func TestMissingTranslation(t *testing.T) {
	// Test behavior with non-existent message key
	fakeKey := MessageKey("non_existent_key")

	result := T(fakeKey)
	if result != string(fakeKey) {
		t.Errorf("Non-existent key should return the key itself, got: %s", result)
	}
}

func TestMessageWithArgs(t *testing.T) {
	originalLang := GetLanguage()
	SetLanguage(English)
	defer SetLanguage(originalLang)

	// Test message with arguments
	listID := "test123"
	result := T(MsgDownloadingList, listID)

	if !strings.Contains(result, listID) {
		t.Errorf("Result should contain the argument '%s', got: %s", listID, result)
	}
}

func BenchmarkTranslation(b *testing.B) {
	SetLanguage(English)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		T(MsgDownloadingSegments)
	}
}

func BenchmarkLanguageDetection(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		detectSystemLanguage()
	}
}
