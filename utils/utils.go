package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// CheckFFmpeg prüft, ob FFmpeg verfügbar ist
func CheckFFmpeg() error {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		return fmt.Errorf("ffmpeg ist nicht installiert oder nicht im PATH: %v", err)
	}
	return nil
}

// EnsureDir erstellt ein Verzeichnis, wenn es nicht existiert
func EnsureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// SanitizeFilename bereinigt den Dateinamen
func SanitizeFilename(filename string) string {
	return strings.Map(func(r rune) rune {
		if r == '/' || r == '\\' || r == ':' || r == '*' || r == '?' || r == '"' || r == '<' || r == '>' || r == '|' {
			return '-'
		}
		return r
	}, filename)
}
