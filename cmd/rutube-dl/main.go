package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/StanislavKH/rutube-dl/pkg/flags"
	"github.com/StanislavKH/rutube-dl/pkg/i18n"
	"github.com/StanislavKH/rutube-dl/pkg/rutubedl"
	"github.com/StanislavKH/rutube-dl/pkg/utils"
)

func main() {
	// Parse command-line flags
	cmdFlags, err := flags.ParseFlags()
	if err != nil {
		log.Fatal(err)
		return
	}

	listID := cmdFlags.ListID
	fileLink := cmdFlags.FileLink
	dir := cmdFlags.Dir
	workers := cmdFlags.Workers
	withFfmpeg := cmdFlags.WithFFmpeg
	fromEpisode := cmdFlags.FromEpisode
	transliterate := cmdFlags.WithTransliteration

	if *listID != "" {
		log.Printf(i18n.T(i18n.MsgDownloadingList), *listID)
		list, err := rutubedl.GetItemsListFromFeedURI(*listID)
		if err != nil {
			log.Fatal(err)
			return
		}

		// Filter episodes based on from_episode parameter
		// Note: from_episode refers to position in playlist (1-based), not episode number
		var filteredList []rutubedl.ProcessedOut
		for i, file := range list {
			playlistPosition := i + 1 // 1-based position in playlist
			if playlistPosition >= *fromEpisode {
				filteredList = append(filteredList, file)
			} else {
				log.Printf(i18n.T(i18n.MsgEpisodeSkipped), file.Title)
			}
		}

		if len(filteredList) == 0 {
			log.Printf(i18n.T(i18n.MsgNoEpisodesToDownload), *fromEpisode)
			return
		}

		// If no custom directory is specified, create a directory based on the playlist name
		outputDir := dir
		if dir == nil || *dir == "" {
			// Get the actual playlist name from metadata
			playlistName := "playlist"
			if meta, err := rutubedl.GetPlaylistMetadata(*listID); err == nil && meta.Name != "" {
				playlistName = meta.Name
			} else if len(list) > 0 && list[0].FeedName != "" {
				// Fallback to channel name if metadata is not available
				playlistName = list[0].FeedName
			}

			// Sanitize the playlist name with the specialized function and apply transliteration if requested
			safePlaylistName := utils.SanitizePlaylistName(playlistName, *transliterate)

			// Create the full path: downloads/[playlist_name]
			outputDir = new(string)
			*outputDir = filepath.Join("downloads", safePlaylistName)
		}

		log.Printf(i18n.T(i18n.MsgDownloadingEpisodes), len(filteredList))

		// Download with proper numbering
		for i, file := range filteredList {
			playlistIndex := *fromEpisode + i // Start numbering from fromEpisode value
			// Calculate the maximum number for proper formatting
			maxNumber := *fromEpisode + len(filteredList) - 1

			log.Printf(i18n.T(i18n.MsgProcessing), file.Title)

			// Use the new DownloadPlaylistFile function with numbering
			err := rutubedl.DownloadPlaylistFile(file.VideoURL, outputDir, *workers, *withFfmpeg, *transliterate, playlistIndex, maxNumber, file.Title)
			if err != nil {
				log.Printf(i18n.T(i18n.MsgErrorDownloading), playlistIndex, file.Title, err)
				continue
			}

			log.Printf(i18n.T(i18n.MsgEpisodeDownloaded), playlistIndex, file.Title)
		}

		log.Printf(i18n.T(i18n.MsgPlaylistCompleted), len(filteredList))
	} else if *fileLink != "" {
		log.Printf(i18n.T(i18n.MsgDownloadingFile), *fileLink)
		err := rutubedl.DownloadFile(*fileLink, dir, *workers, *withFfmpeg, *transliterate)
		if err != nil {
			log.Fatalf(i18n.T(i18n.MsgFailedToDownloadFile), err)
		}
	} else {
		fmt.Println(i18n.T(i18n.MsgUsage))
		fmt.Println(i18n.T(i18n.MsgExample))
		fmt.Println(i18n.T(i18n.MsgDownloadList))
		fmt.Println(i18n.T(i18n.MsgDownloadFile))
		fmt.Println("\n" + i18n.T(i18n.MsgNote))
		fmt.Println(i18n.T(i18n.MsgDirOptional))
		fmt.Println(i18n.T(i18n.MsgWorkersOptional))
		fmt.Println(i18n.T(i18n.MsgWithFfmpegOptional))
		fmt.Println(i18n.T(i18n.MsgFromEpisodeOptional))
		fmt.Println(i18n.T(i18n.MsgWithTransliterationOptional))
	}
}
