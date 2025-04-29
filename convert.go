package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/dhowden/tag"
	"gopkg.in/yaml.v3"
)

type Metadata struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Published   int    `yaml:"published"`
	File        string `yaml:"file"`
	Duration    string `yaml:"duration"`
	Length      int64  `yaml:"length"`
}

func getAudioFiles() ([]string, error) {
	audioDir := "audio"

	// Check if the directory exists
	if _, err := os.Stat(audioDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory does not exist: %s", audioDir)
	}

	files, err := os.ReadDir(audioDir)
	if err != nil {
		return nil, err
	}

	var audioFiles []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".mp3" {
			audioFiles = append(audioFiles, filepath.Join(audioDir, file.Name()))
		}
	}
	return audioFiles, nil
}

func readID3Metadata(filePath string) (*Metadata, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	metadata, err := tag.ReadFrom(file)
	if err != nil {
		return nil, fmt.Errorf("error reading metadata: %v", err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("error getting file info: %v", err)
	}

	duration := "" // Placeholder for actual duration calculation
	// Use an audio processing library to calculate the duration if needed
	return &Metadata{
		Title:       metadata.Title(),
		Description: metadata.Comment(), //"'" + metadata.Comment() + "'", // Using the comment field as a description
		Published:   metadata.Year(),
		File:        filePath,
		Duration:    duration, // Placeholder for actual duration
		Length:      fileInfo.Size(),
	}, nil
}

func writeMetadataToYAML(metadata []*Metadata, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating YAML file: %v", err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	if err := encoder.Encode(metadata); err != nil {
		return fmt.Errorf("error writing YAML: %v", err)
	}

	return nil
}

func main() {
	audioFiles, err := getAudioFiles()
	if err != nil {
		log.Fatalf("Error reading audio files: %v", err)
	}
	fmt.Println("Audio files found:", audioFiles)

	var metadataList []*Metadata
	for _, filePath := range audioFiles {
		fmt.Printf("Reading metadata for: %s\n", filePath)
		metadata, err := readID3Metadata(filePath)
		if err != nil {
			log.Printf("Error reading metadata for %s: %v", filePath, err)
			continue
		}
		metadataList = append(metadataList, metadata)
	}

	outputPath := "episodes_go.yaml"
	if err := writeMetadataToYAML(metadataList, outputPath); err != nil {
		log.Fatalf("Error writing metadata to YAML: %v", err)
	}

	fmt.Printf("Metadata written to %s\n", outputPath)
}
