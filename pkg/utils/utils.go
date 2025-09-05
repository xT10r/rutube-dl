package utils

import (
	"fmt"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"unicode"
)

// ForbiddenChars contains characters that are not allowed in filenames on various OS
var ForbiddenChars = map[string]string{
	"windows": `<>:"/\\|?*`,
	"linux":   `/`,
	"darwin":  `/:`,
}

// TransliterationMap maps Cyrillic characters to Latin equivalents
var TransliterationMap = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo", 'ж': "zh",
	'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m", 'н': "n", 'о': "o",
	'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u", 'ф': "f", 'х': "kh", 'ц': "ts",
	'ч': "ch", 'ш': "sh", 'щ': "shch", 'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
	'А': "A", 'Б': "B", 'В': "V", 'Г': "G", 'Д': "D", 'Е': "E", 'Ё': "Yo", 'Ж': "Zh",
	'З': "Z", 'И': "I", 'Й': "Y", 'К': "K", 'Л': "L", 'М': "M", 'Н': "N", 'О': "O",
	'П': "P", 'Р': "R", 'С': "S", 'Т': "T", 'У': "U", 'Ф': "F", 'Х': "Kh", 'Ц': "Ts",
	'Ч': "Ch", 'Ш': "Sh", 'Щ': "Shch", 'Ъ': "", 'Ы': "Y", 'Ь': "", 'Э': "E", 'Ю': "Yu", 'Я': "Ya",
}

// SanitizeFilename removes or replaces characters that are not allowed in filenames
func SanitizeFilename(filename string, transliterate bool) string {
	// Apply transliteration if requested
	if transliterate {
		filename = TransliterateText(filename)
	}

	// Get forbidden characters for current OS
	forbidden := getForbiddenCharsForOS()

	// Replace forbidden characters with safe alternatives
	for _, char := range forbidden {
		filename = strings.ReplaceAll(filename, string(char), "_")
	}

	// Remove control characters and other problematic Unicode characters
	filename = removeControlCharacters(filename)

	// Trim spaces and dots from the beginning and end
	filename = strings.Trim(filename, " .")

	// Ensure filename is not empty
	if filename == "" {
		filename = "untitled"
	}

	// Limit filename length (most filesystems support 255 characters)
	if len(filename) > 200 {
		filename = filename[:200]
	}

	return filename
}

// SanitizePlaylistName specifically sanitizes playlist names to remove emojis and special characters
// while preserving meaningful text
func SanitizePlaylistName(name string, transliterate bool) string {
	// Apply transliteration if requested
	if transliterate {
		name = TransliterateText(name)
	}

	// Get forbidden characters for current OS
	forbidden := getForbiddenCharsForOS()

	// Replace forbidden characters with safe alternatives
	for _, char := range forbidden {
		name = strings.ReplaceAll(name, string(char), "_")
	}

	// Remove control characters and other problematic Unicode characters including emojis
	name = removeSpecialCharacters(name)

	// Clean up extra spaces and special formatting
	name = cleanWhitespace(name)

	// Trim spaces and dots from the beginning and end
	name = strings.Trim(name, " .")

	// Ensure name is not empty
	if name == "" {
		name = "playlist"
	}

	// Limit name length
	if len(name) > 200 {
		name = name[:200]
	}

	return name
}

// TransliterateText converts Cyrillic text to Latin characters
func TransliterateText(text string) string {
	var result strings.Builder

	for _, char := range text {
		if replacement, exists := TransliterationMap[char]; exists {
			result.WriteString(replacement)
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}

// getForbiddenCharsForOS returns forbidden characters for the current OS
func getForbiddenCharsForOS() string {
	switch runtime.GOOS {
	case "windows":
		return ForbiddenChars["windows"]
	case "darwin":
		return ForbiddenChars["darwin"]
	default: // Linux and other Unix-like systems
		return ForbiddenChars["linux"]
	}
}

// removeControlCharacters removes control characters (0x00-0x1F and 0x7F-0x9F)
func removeControlCharacters(text string) string {
	controlCharsRegex := regexp.MustCompile(`[\x00-\x1F\x7F-\x9F]`)
	return controlCharsRegex.ReplaceAllString(text, "")
}

// removeSpecialCharacters removes emojis and other special Unicode characters
// while preserving letters, numbers, and basic punctuation
func removeSpecialCharacters(text string) string {
	var result strings.Builder

	for _, r := range text {
		// Allow basic Latin and Cyrillic letters, numbers, and common punctuation
		if unicode.IsLetter(r) || unicode.IsNumber(r) ||
			unicode.IsSpace(r) ||
			r == '-' || r == '_' || r == '(' || r == ')' ||
			r == '[' || r == ']' || r == '.' || r == ',' {
			result.WriteRune(r)
		}
		// Optionally allow some other common characters if needed
	}

	return result.String()
}

// cleanWhitespace normalizes whitespace and removes extra spaces
func cleanWhitespace(text string) string {
	// Replace multiple spaces with single space
	spaceRegex := regexp.MustCompile(`\s+`)
	text = spaceRegex.ReplaceAllString(text, " ")
	return text
}

// CreateSafeDirectory creates a directory with a safe name
func CreateSafeDirectory(basePath, dirName string, transliterate bool) (string, error) {
	safeDirName := SanitizeFilename(dirName, transliterate)
	fullPath := filepath.Join(basePath, safeDirName)
	return fullPath, nil
}

// GetSafeFilename generates a safe filename with extension
func GetSafeFilename(title, extension string, transliterate bool) string {
	safeTitle := SanitizeFilename(title, transliterate)
	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}
	return safeTitle + extension
}

// FormatPlaylistIndex formats the playlist index based on total count
// Returns format like "01", "001", etc. depending on total items
func FormatPlaylistIndex(index, total int) string {
	if total < 10 {
		return fmt.Sprintf("%d", index)
	} else if total < 100 {
		return fmt.Sprintf("%02d", index)
	} else if total < 1000 {
		return fmt.Sprintf("%03d", index)
	} else {
		return fmt.Sprintf("%04d", index)
	}
}

// GetPlaylistFilename generates a filename with playlist index
func GetPlaylistFilename(index, total int, title, extension string, transliterate bool) string {
	formattedIndex := FormatPlaylistIndex(index, total)
	safeTitle := SanitizeFilename(title, transliterate)
	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}
	return fmt.Sprintf("%s - %s%s", formattedIndex, safeTitle, extension)
}
