package i18n

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// Message keys for internationalization
type MessageKey string

const (
	MsgDownloadingList             MessageKey = "downloading_list"
	MsgDownloadingFile             MessageKey = "downloading_file"
	MsgProcessing                  MessageKey = "processing"
	MsgEpisodeSkipped              MessageKey = "episode_skipped"
	MsgUsage                       MessageKey = "usage"
	MsgExample                     MessageKey = "example"
	MsgDownloadList                MessageKey = "download_list"
	MsgDownloadFile                MessageKey = "download_file"
	MsgNote                        MessageKey = "note"
	MsgDirOptional                 MessageKey = "dir_optional"
	MsgWorkersOptional             MessageKey = "workers_optional"
	MsgWithFfmpegOptional          MessageKey = "with_ffmpeg_optional"
	MsgFromEpisodeOptional         MessageKey = "from_episode_optional"
	MsgWithTransliterationOptional MessageKey = "with_transliteration_optional"
	MsgDownloadingSegments         MessageKey = "downloading_segments"
	MsgAllSegmentsMerged           MessageKey = "all_segments_merged"
	MsgVideoDownloaded             MessageKey = "video_downloaded"
	MsgChunksProcessed             MessageKey = "chunks_processed"
	MsgRetryAttempt                MessageKey = "retry_attempt"
	MsgAttemptFailed               MessageKey = "attempt_failed"
	MsgSegmentNotAvailable         MessageKey = "segment_not_available"
	MsgFFmpegDownloading           MessageKey = "ffmpeg_downloading"
	MsgFFmpegUpdating              MessageKey = "ffmpeg_updating"
	MsgFFmpegReady                 MessageKey = "ffmpeg_ready"
	MsgNoEpisodesToDownload        MessageKey = "no_episodes_to_download"
	MsgDownloadingEpisodes         MessageKey = "downloading_episodes"
	MsgErrorDownloading            MessageKey = "error_downloading"
	MsgEpisodeDownloaded           MessageKey = "episode_downloaded"
	MsgPlaylistCompleted           MessageKey = "playlist_completed"
	MsgFailedToDownloadFile        MessageKey = "failed_to_download_file"
	MsgFlagListID                  MessageKey = "flag_list_id"
	MsgFlagFileLink                MessageKey = "flag_file_link"
	MsgFlagDir                     MessageKey = "flag_dir"
	MsgFlagWorkers                 MessageKey = "flag_workers"
	MsgFlagWithFFmpeg              MessageKey = "flag_with_ffmpeg"
	MsgFlagFromEpisode             MessageKey = "flag_from_episode"
	MsgFlagWithTransliteration     MessageKey = "flag_with_transliteration"
)

// Language represents supported languages
type Language string

const (
	English Language = "en"
	Russian Language = "ru"
)

// Messages stores localized messages
var messages = map[Language]map[MessageKey]string{
	English: {
		MsgDownloadingList:             "Downloading list with ID: %s",
		MsgDownloadingFile:             "Downloading file from URL: %s",
		MsgProcessing:                  "processing: %s",
		MsgEpisodeSkipped:              "episode %s skipped",
		MsgUsage:                       "Usage: Provide either a -list_id or a -file_link option",
		MsgExample:                     "Example:",
		MsgDownloadList:                "  - To download a list: rutube-dl -list_id=<list-id> [-dir=<directory>] [-workers=<number>] [-with_ffmpeg] [-from_episode=<episode-number>] [-with_transliteration]",
		MsgDownloadFile:                "  - To download a file: rutube-dl -file_link=<file-url> [-dir=<directory>] [-workers=<number>] [-with_ffmpeg] [-with_transliteration]",
		MsgNote:                        "Note:",
		MsgDirOptional:                 "  - The -dir option is optional. If not specified, it defaults to 'downloads'.",
		MsgWorkersOptional:             "  - The -workers option is optional and defaults to 1. Be cautious when using more than 5.",
		MsgWithFfmpegOptional:          "  - The -with_ffmpeg option is optional. If enabled, external ffmpeg is used for more reliable chunk concatenation.",
		MsgFromEpisodeOptional:         "  - The -from_episode option is optional and can only be used with -list_id. It specifies the position in the playlist (1-based) from which to start downloading.",
		MsgWithTransliterationOptional: "  - The -with_transliteration option is optional. If enabled, video names will be transliterated to Latin characters.",
		MsgDownloadingSegments:         "Downloading segments",
		MsgAllSegmentsMerged:           "All segments merged into: %s",
		MsgVideoDownloaded:             "Video downloaded and saved to: %s",
		MsgChunksProcessed:             "Chunks successfully processed and concatenated.",
		MsgRetryAttempt:                "Attempt %d/%d: Failed to fetch URL, retrying in %v...",
		MsgAttemptFailed:               "Attempt %d: Failed to download segment, retrying...",
		MsgSegmentNotAvailable:         "segments URI not available, will try next one",
		MsgFFmpegDownloading:           "Downloading FFmpeg...",
		MsgFFmpegUpdating:              "Updating FFmpeg to latest version...",
		MsgFFmpegReady:                 "FFmpeg is ready for use.",
		MsgNoEpisodesToDownload:        "No episodes to download (all episodes before episode %d)",
		MsgDownloadingEpisodes:         "Downloading %d episodes from playlist",
		MsgErrorDownloading:            "Error downloading episode %d (%s): %v",
		MsgEpisodeDownloaded:           "Successfully downloaded episode %d: %s",
		MsgPlaylistCompleted:           "Playlist download completed. Downloaded %d episodes.",
		MsgFailedToDownloadFile:        "failed to download file: %v",
		MsgFlagListID:                  "ID of the list to download",
		MsgFlagFileLink:                "Direct URL to the file to download",
		MsgFlagDir:                     "directory to store files and temporary chunks, default is - downloads",
		MsgFlagWorkers:                 "number of workers for chunk download; be careful when using more than 5",
		MsgFlagWithFFmpeg:              "use external ffmpeg for more reliable chunk concatenation",
		MsgFlagFromEpisode:             "download everything starting from this position in playlist (1-based, only with -list_id)",
		MsgFlagWithTransliteration:     "transliterate video names to Latin characters",
	},
	Russian: {
		MsgDownloadingList:             "Загрузка списка с ID: %s",
		MsgDownloadingFile:             "Загрузка файла по URL: %s",
		MsgProcessing:                  "обработка: %s",
		MsgEpisodeSkipped:              "эпизод %s пропущен",
		MsgUsage:                       "Использование: Укажите либо -list_id, либо -file_link",
		MsgExample:                     "Пример:",
		MsgDownloadList:                "  - Для загрузки списка: rutube-dl -list_id=<list-id> [-dir=<каталог>] [-workers=<количество>] [-with_ffmpeg] [-from_episode=<номер-эпизода>] [-with_transliteration]",
		MsgDownloadFile:                "  - Для загрузки файла: rutube-dl -file_link=<url-файла> [-dir=<каталог>] [-workers=<количество>] [-with_ffmpeg] [-with_transliteration]",
		MsgNote:                        "Примечание:",
		MsgDirOptional:                 "  - Опция -dir необязательная. Если не указана, по умолчанию используется 'downloads'.",
		MsgWorkersOptional:             "  - Опция -workers необязательная, по умолчанию 1. Будьте осторожны при использовании более 5.",
		MsgWithFfmpegOptional:          "  - Опция -with_ffmpeg необязательная. Если включена, используется внешний ffmpeg для более надёжного объединения частей.",
		MsgFromEpisodeOptional:         "  - Опция -from_episode необязательная и может использоваться только с -list_id. Указывает позицию в плейлисте (начиная с 1), с которой начать загрузку.",
		MsgWithTransliterationOptional: "  - Опция -with_transliteration необязательная. Если включена, имена видео будут транслитерированы в латинские символы.",
		MsgDownloadingSegments:         "Загрузка сегментов",
		MsgAllSegmentsMerged:           "Все сегменты объединены в: %s",
		MsgVideoDownloaded:             "Видео загружено и сохранено в: %s",
		MsgChunksProcessed:             "Части успешно обработаны и объединены.",
		MsgRetryAttempt:                "Попытка %d/%d: Не удалось получить URL, повтор через %v...",
		MsgAttemptFailed:               "Попытка %d: Не удалось загрузить сегмент, повторяю...",
		MsgSegmentNotAvailable:         "URI сегментов недоступен, попробую следующий",
		MsgFFmpegDownloading:           "Загрузка FFmpeg...",
		MsgFFmpegUpdating:              "Обновление FFmpeg до последней версии...",
		MsgFFmpegReady:                 "FFmpeg готов к использованию.",
		MsgNoEpisodesToDownload:        "Нет эпизодов для загрузки (все эпизоды до эпизода %d)",
		MsgDownloadingEpisodes:         "Загрузка %d эпизодов из плейлиста",
		MsgErrorDownloading:            "Ошибка загрузки эпизода %d (%s): %v",
		MsgEpisodeDownloaded:           "Эпизод %d успешно загружен: %s",
		MsgPlaylistCompleted:           "Загрузка плейлиста завершена. Загружено %d эпизодов.",
		MsgFailedToDownloadFile:        "не удалось загрузить файл: %v",
		MsgFlagListID:                  "ID списка для загрузки",
		MsgFlagFileLink:                "Прямой URL файла для загрузки",
		MsgFlagDir:                     "каталог для хранения файлов и временных частей, по умолчанию - downloads",
		MsgFlagWorkers:                 "количество потоков для загрузки частей; будьте осторожны при использовании более 5",
		MsgFlagWithFFmpeg:              "использовать внешний ffmpeg для более надежного объединения частей",
		MsgFlagFromEpisode:             "загружать всё, начиная с этой позиции в плейлисте (1-based, только с -list_id)",
		MsgFlagWithTransliteration:     "транслитерировать названия видео в латинские символы",
	},
}

var currentLanguage Language

// init automatically detects the OS language
func init() {
	currentLanguage = detectSystemLanguage()
}

// detectSystemLanguage detects the system language based on environment variables
func detectSystemLanguage() Language {
	// For Windows
	if runtime.GOOS == "windows" {
		if lang := os.Getenv("LANG"); lang != "" {
			return parseLanguage(lang)
		}
		// Check Windows-specific environment variables
		if lang := os.Getenv("LC_ALL"); lang != "" {
			return parseLanguage(lang)
		}
		if lang := os.Getenv("LC_MESSAGES"); lang != "" {
			return parseLanguage(lang)
		}
		if lang := os.Getenv("LANGUAGE"); lang != "" {
			return parseLanguage(lang)
		}
	}

	// For Unix-like systems (Linux, macOS, etc.)
	envVars := []string{"LC_ALL", "LC_MESSAGES", "LANG", "LANGUAGE"}
	for _, envVar := range envVars {
		if lang := os.Getenv(envVar); lang != "" {
			return parseLanguage(lang)
		}
	}

	// Default to English if no language detected
	return English
}

// parseLanguage parses language code from environment variable
func parseLanguage(langCode string) Language {
	langCode = strings.ToLower(langCode)

	// Handle common Russian language codes
	if strings.HasPrefix(langCode, "ru") ||
		strings.Contains(langCode, "russian") ||
		strings.Contains(langCode, "rus") {
		return Russian
	}

	// Default to English for all other cases
	return English
}

// SetLanguage allows manual language setting
func SetLanguage(lang Language) {
	currentLanguage = lang
}

// GetLanguage returns the current language
func GetLanguage() Language {
	return currentLanguage
}

// T translates a message key to the current language
func T(key MessageKey, args ...interface{}) string {
	if msgs, exists := messages[currentLanguage]; exists {
		if msg, exists := msgs[key]; exists {
			if len(args) > 0 {
				return fmt.Sprintf(msg, args...)
			}
			return msg
		}
	}

	// Fallback to English if translation not found
	if msgs, exists := messages[English]; exists {
		if msg, exists := msgs[key]; exists {
			if len(args) > 0 {
				return fmt.Sprintf(msg, args...)
			}
			return msg
		}
	}

	// Final fallback
	return string(key)
}

// Tf is an alias for T for shorter usage
func Tf(key MessageKey, args ...interface{}) string {
	return T(key, args...)
}
