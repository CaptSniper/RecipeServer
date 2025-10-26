package ars

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"io"
	"os"
	"path/filepath"

	rfp "github.com/CaptSniper/RecipeServer/RFP"
	"github.com/PuerkitoBio/goquery"
)

// Ingredient represents a single ingredient with quantity, unit, and name
type Ingredient struct {
	Quantity string
	Unit     string
	Name     string
}

// ScrapeAllRecipes fetches a URL from AllRecipes and extracts structured recipe data
func ScrapeAllRecipes(url, imagePath string) (*rfp.Recipe, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	data := &rfp.Recipe{}

	// --- 1. Image ---
	doc.Find("div#article__photo-ribbon_1-0 a").First().Find("img").Each(func(i int, s *goquery.Selection) {
		if src, exists := s.Attr("src"); exists {
			DownloadImage(src, imagePath)
			data.ImagePath = src
		}
	})

	// --- 2. Times & Servings ---
	doc.Find("div#mm-recipes-details_1-0 div").Each(func(i int, s *goquery.Selection) {
		if s.HasClass("mm-recipes-details__label") {
			label := strings.TrimSpace(s.Text())
			valueSel := s.Next() // value div should be next sibling
			if valueSel != nil && valueSel.HasClass("mm-recipes-details__value") {
				value := strings.TrimSpace(valueSel.Text())
				switch strings.ToLower(label) {
				case "prep time":
					v, _ := strconv.Atoi(value)
					data.PrepTime = uint16(v)
				case "cook time":
					v, _ := strconv.Atoi(value)
					data.CookTime = uint16(v)
				case "additional time":
					v, _ := strconv.Atoi(value)
					data.AdditionalTime = uint16(v)
				case "total time":
					v, _ := strconv.Atoi(value)
					data.TotalTime = uint16(v)
				case "servings":
					data.Servings = value
				}
			}
		}
	})

	// --- 3. Ingredients ---
	doc.Find("div.mm-recipes-structured-ingredients__list").First().Find("li").Each(func(i int, s *goquery.Selection) {
		p := s.Find("p").First()
		spans := p.Find("span")
		ing := Ingredient{
			Quantity: strings.TrimSpace(spans.Eq(0).Text()),
			Unit:     strings.TrimSpace(spans.Eq(1).Text()),
			Name:     strings.TrimSpace(spans.Eq(2).Text()),
		}

		data.Ingredients = append(data.Ingredients, ing.Quantity, ing.Unit, ing.Name)
	})

	// --- 4. Steps / Directions ---
	doc.Find("div#mm-recipes-steps__content_1-0 ol li").Each(func(i int, s *goquery.Selection) {
		stepText := strings.TrimSpace(s.Find("p").Text())
		if stepText != "" {
			data.Steps = append(data.Steps, stepText)
		}
	})

	return data, nil
}

// DownloadImage downloads an image from the given URL and saves it to the specified directory.
// It returns the full path to the saved image.
func DownloadImage(url, saveDir string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("empty image URL")
	}

	// Ensure the save directory exists
	if _, err := os.Stat(saveDir); os.IsNotExist(err) {
		if err := os.MkdirAll(saveDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory: %v", err)
		}
	}

	// Extract the filename from the URL
	tokens := strings.Split(url, "/")
	filename := tokens[len(tokens)-1]
	if filename == "" {
		filename = "image.jpg"
	}

	savePath := filepath.Join(saveDir, filename)

	// Download the image
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	// Create the file
	out, err := os.Create(savePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	// Copy the data
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save image: %v", err)
	}

	return savePath, nil
}
