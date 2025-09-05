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
	fromEpisode := flag.Int("from_episode", 1, "download everything starting from this episode number (only with -list_id)")
	transliterate := flag.Bool("transliterate", false, "transliterate video names to Latin characters")

	flag.Parse()

	if *listID != "" {
		log.Printf(i18n.T(i18n.MsgDownloadingList), *listID)
		list, err := rutubedl.GetItemsListFromFeedURI(*listID)
		if err != nil {
			log.Fatal(err)
			return
		}
		for _, file := range list {
			log.Printf(i18n.T(i18n.MsgProcessing), file.Title)
			if file.Episode < *fromEpisode {
				log.Printf(i18n.T(i18n.MsgEpisodeSkipped), file.Title)
				continue
			}
			err := rutubedl.DownloadFile(file.VideoURL, dir, *workers, *withFfmpeg, *transliterate)
			if err != nil {
				log.Println(err)
				continue
			}
		}
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
