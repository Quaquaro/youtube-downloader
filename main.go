package main

import (
	"flag"
	"fmt"
	"os"

	"youtube-downloader/down"
)

func main() {
	youtubeURL := flag.String("u", "", "Die YouTube-URL des Videos")
	playlistURL := flag.String("p", "", "Die YouTube-URL der Playlist")
	outputDir := flag.String("d", ".", "Das Verzeichnis der Ausgabedateien")

	flag.Parse()

	if *youtubeURL == "" && *playlistURL == "" {
		fmt.Println("Bitte gib eine YouTube-URL oder eine Playlist-URL ein")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *youtubeURL != "" {
		fmt.Println("Verarbeite einzelnes Video...")
		err := down.DownloadAndConvert(*youtubeURL, *outputDir)
		if err != nil {
			fmt.Printf("Fehler: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Video erfolgreich heruntergeladen und konvertiert.")
	} else if *playlistURL != "" {
		fmt.Println("Verarbeite Playlist...")
		err := down.DownloadPlaylist(*playlistURL, *outputDir)
		if err != nil {
			fmt.Printf("Fehler: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Playlist erfolgreich verarbeitet.")
	}
}
