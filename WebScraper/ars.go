package ars

import (
	"errors"
	"fmt"
	"net/http"
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

func ScrapeRecipe(url, imagePath string) (*rfp.Recipe, error) {
	if strings.Contains(url, "allrecipes.com") {
		return ScrapeAllRecipes(url, imagePath)
	} else {
		err := errors.New("unsupported base website")
		return nil, err
	}
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

	data := rfp.NewRecipe()

	recipeName := strings.TrimSpace(doc.Find("div#article-header--recipe_1-0 h1").First().Text())
	data.Name = recipeName

	// --- 1. Image ---
	doc.Find("div#photo-dialog__item_1-0 img").Each(func(i int, s *goquery.Selection) {
		if src, exists := s.Attr("src"); exists && src != "" {
			DownloadImage(src, imagePath, data.Name)
			data.ImagePath = imagePath + "\\" + data.Name
		}
	})

	// --- 2. Times & Servings ---
	doc.Find("div#mm-recipes-details_1-0 div.mm-recipes-details__item").Each(func(i int, s *goquery.Selection) {
		label := strings.TrimSpace(s.Find("div.mm-recipes-details__label").Text())
		value := strings.TrimSpace(s.Find("div.mm-recipes-details__value").Text())

		switch strings.ToLower(label) {
		case "prep time:":
			data.CoreProps["prep time"] = value
		case "cook time:":
			data.CoreProps["cook time"] = value
		case "total time:":
			data.CoreProps["total time"] = value
		case "additional time:":
			data.CoreProps["additional time"] = value
		case "servings:":
			data.CoreProps["servings"] = value
		}
	})

	// --- 3. Ingredients ---
	doc.Find("div#mm-recipes-structured-ingredients_1-0").Find("ul").First().Find("li").Each(func(i int, s *goquery.Selection) {
		p := s.Find("p").First()
		spans := p.Find("span")
		ing := Ingredient{
			Quantity: strings.TrimSpace(spans.Eq(0).Text()),
			Unit:     strings.TrimSpace(spans.Eq(1).Text()),
			Name:     strings.TrimSpace(spans.Eq(2).Text()),
		}

		data.Ingredients = append(data.Ingredients, ing.Quantity+" "+ing.Unit+" "+ing.Name)
	})

	// --- 4. Steps / Directions ---
	doc.Find("div#mm-recipes-steps__content_1-0 ol li").Each(func(i int, s *goquery.Selection) {
		stepText := strings.TrimSpace(s.Find("p").First().Text())
		if stepText != "" {
			data.Steps = append(data.Steps, stepText)
		}
	})

	return data, nil
}

// DownloadImage downloads an image from the given URL and saves it to the specified directory.
// It returns the full path to the saved image.
func DownloadImage(url, saveDir, saveName string) (string, error) {
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
		filename = saveName + ".jpg"
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
