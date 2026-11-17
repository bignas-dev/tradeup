// Package images downloads and converts skin and collection images.
package images

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	skinsURL             = "https://raw.githubusercontent.com/ByMykel/CSGO-API/main/public/api/en/skins.json"
	skinImagesFolderPath = "./.tmp/skins"
)

var supportedImageExtensions = []string{"png", "svg"}

type imagesResponse struct {
	Name string `json:"name"`
	URL  string `json:"image"`
}

func SaveImages() error {
	log.Println("images: starting download")

	if err := os.MkdirAll(skinImagesFolderPath, 0o755); err != nil {
		return fmt.Errorf("images: failed to make dir: %w", err)
	}

	if err := getImagesFromURL(skinImagesFolderPath, skinsURL); err != nil {
		return err
	}

	log.Println("images: done")
	return nil
}

func getImagesFromURL(folderPath string, imageURL string) error {
	resp, err := http.Get(imageURL)
	if err != nil {
		return fmt.Errorf("images: failed to make a GET %v: %w", imageURL, err)
	}
	defer resp.Body.Close()

	var images []imagesResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&images); err != nil {
		return fmt.Errorf("images: failed to decode json: %w", err)
	}

	log.Printf("images: %d images in source", len(images))

	if err := extractImages(folderPath, images); err != nil {
		return fmt.Errorf("images: failed to extract images: %w", err)
	}

	return nil
}

func extractImages(folderPath string, imageData []imagesResponse) error {
	for _, img := range imageData {
		rawImgPath := filepath.Join(folderPath, img.Name)
		extension := getExistingImageExtension(rawImgPath)

		if extension == "" {
			log.Printf("images: downloading %s", img.Name)
			if err := downloadRawImage(img.URL, rawImgPath); err != nil {
				return fmt.Errorf("failed to download image: %w", err)
			}
			log.Printf("images: downloaded %s", img.Name)
		}
	}
	return nil
}

func imageExists(imgPath string) bool {
	f, err := os.Open(imgPath)
	if err != nil {
		return false
	}
	defer f.Close()
	return true
}

func getExistingImageExtension(imgPathNoExtension string) string {
	for _, ext := range supportedImageExtensions {
		if imageExists(imgPathNoExtension + "." + ext) {
			return ext
		}
	}
	return ""
}

func downloadRawImage(imgURL, imgPath string) error {
	resp, err := http.Get(imgURL)
	if err != nil {
		return fmt.Errorf("downloading %q: %w", imgURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("URL %v status code: %v", imgURL, resp.StatusCode)
	}

	imgType := strings.ToLower(resp.Header.Get("Content-Type"))
	extension := ""
	for _, ext := range supportedImageExtensions {
		if strings.Contains(imgType, ext) {
			extension = ext
			break
		}
	}

	if extension == "" {
		return fmt.Errorf("invalid image extension: %v", imgType)
	}

	file, err := os.Create(imgPath + "." + extension)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		_ = os.Remove(imgPath)
		return fmt.Errorf("copying response body: %w", err)
	}

	return nil
}
