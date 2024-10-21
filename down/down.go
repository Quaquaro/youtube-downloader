package down

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"youtube-downloader/utils"

	"github.com/kkdai/youtube/v2"
)

// DownloadAndConvert lädt ein Video herunter und konvertiert es in MP3
func DownloadAndConvert(youtubeURL, outputDir string) error {
	if err := utils.CheckFFmpeg(); err != nil {
		return err
	}

	client := youtube.Client{}
	video, err := client.GetVideo(youtubeURL)
	if err != nil {
		return fmt.Errorf("Fehler beim Abrufen des Videos: %v", err)
	}

	formats := video.Formats.WithAudioChannels()
	if len(formats) == 0 {
		return fmt.Errorf("Keine geeigneten Audioformate gefunden")
	}
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		return fmt.Errorf("Fehler beim Abrufen des Streams: %v", err)
	}
	defer stream.Close()

	if err := utils.EnsureDir(outputDir); err != nil {
		return fmt.Errorf("Fehler beim Erstellen des Ausgabeverzeichnisses: %v", err)
	}

	outputFilename := filepath.Join(outputDir, utils.SanitizeFilename(video.Title)+".mp4")
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("Fehler beim Erstellen der Ausgabedatei: %v", err)
	}
	defer outputFile.Close()

	_, err = outputFile.ReadFrom(stream)
	if err != nil {
		return fmt.Errorf("Fehler beim Schreiben der Ausgabedatei: %v", err)
	}

	mp3Filename := filepath.Join(outputDir, video.Title+".mp3")
	cmd := exec.Command("ffmpeg", "-i", outputFilename, "-vn", "-b:a", "192k", "-filter:a", "volume=1.0", mp3Filename)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Fehler beim Konvertieren: %v", err)
	}

	if err := os.Remove(outputFilename); err != nil {
		return fmt.Errorf("Fehler beim Löschen der Ausgabedatei: %v", err)
	}

	return nil
}

// DownloadPlaylist lädt eine YouTube-Playlist herunter und konvertiert die Videos
func DownloadPlaylist(playlistURL, outputDir string) error {
	client := youtube.Client{}
	playlist, err := client.GetPlaylist(playlistURL)
	if err != nil {
		return fmt.Errorf("Fehler beim Abrufen der Playlist: %v", err)
	}

	if err := utils.EnsureDir(outputDir); err != nil {
		return fmt.Errorf("Fehler beim Erstellen des Ausgabeverzeichnisses: %v", err)
	}

	for i, video := range playlist.Videos {
		fmt.Printf("Verarbeite Video %d von %d: %s\n", i+1, len(playlist.Videos), video.Title)
		err := DownloadAndConvert(video.ID, outputDir)
		if err != nil {
			return fmt.Errorf("Fehler beim Herunterladen und Konvertieren des Videos: %s: %v", video.Title, err)
		}
	}

	return nil
}
