package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/StanislavKH/rutube-dl/pkg/i18n"
	"github.com/StanislavKH/rutube-dl/pkg/rutubedl"
)

func main() {
	listID := flag.String("list_id", "", "ID of the list to download")
	fileLink := flag.String("file_link", "", "Direct URL to the file to download")
	dir := flag.String("dir", "", "directory to store files and temporary chunks, default is - downloads")
	workers := flag.Int("workers", 1, "number of workers for chunk download; be careful when using more than 5")
	withFfmpeg := flag.Bool("with_ffmpeg", false, "use external ffmpeg for more reliable chunk concatenation")
	fromEpisode := flag.Int("from_episode", 1, "download everything starting from this position in playlist (1-based, only with -list_id)")
	transliterate := flag.Bool("transliterate", false, "transliterate video names to Latin characters")

	flag.Parse()

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
			log.Printf("No episodes to download (all episodes before episode %d)", *fromEpisode)
			return
		}

		log.Printf("Downloading %d episodes from playlist", len(filteredList))

		// Download with proper numbering
		for i, file := range filteredList {
			playlistIndex := *fromEpisode + i // Start numbering from fromEpisode value
			// Calculate the maximum number for proper formatting
			maxNumber := *fromEpisode + len(filteredList) - 1

			log.Printf(i18n.T(i18n.MsgProcessing), file.Title)

			// Use the new DownloadPlaylistFile function with numbering
			err := rutubedl.DownloadPlaylistFile(file.VideoURL, dir, *workers, *withFfmpeg, *transliterate, playlistIndex, maxNumber, file.Title)
			if err != nil {
				log.Printf("Error downloading episode %d (%s): %v", playlistIndex, file.Title, err)
				continue
			}

			log.Printf("Successfully downloaded episode %d: %s", playlistIndex, file.Title)
		}

		log.Printf("Playlist download completed. Downloaded %d episodes.", len(filteredList))
	} else if *fileLink != "" {
		log.Printf(i18n.T(i18n.MsgDownloadingFile), *fileLink)
		err := rutubedl.DownloadFile(*fileLink, dir, *workers, *withFfmpeg, *transliterate)
		if err != nil {
			log.Fatalf("failed to download file: %v", err)
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
		fmt.Println(i18n.T(i18n.MsgTransliterateOptional))
	}
}
