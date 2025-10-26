package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	rfp "github.com/CaptSniper/RecipeServer/RFP"
	ars "github.com/CaptSniper/RecipeServer/WebScraper"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Recipe File Program")
	fmt.Println("-------------------")
	fmt.Println("Choose an option:")
	fmt.Println("1) Create a recipe file")
	fmt.Println("2) Read a recipe file")
	fmt.Println("3) Create/Edit config")
	fmt.Println("4) Scrape AllRecipes")
	fmt.Print("> ")

	var choice int
	fmt.Scan(&choice)

	reader.ReadString('\n')

	switch choice {
	case 1:
		createRecipe(reader)
	case 2:
		readRecipe(reader)
	case 3:
		editConfig()
	case 4:
		ScrapeAS()
	default:
		fmt.Println("Unknown option")
	}
}

func createRecipe(reader *bufio.Reader) {
	var r rfp.Recipe

	fmt.Print("Image path: ")
	r.ImagePath, _ = reader.ReadString('\n')
	r.ImagePath = strings.TrimSpace(r.ImagePath)

	fmt.Print("Prep time (min): ")
	fmt.Scan(&r.PrepTime)
	fmt.Print("Cook time (min): ")
	fmt.Scan(&r.CookTime)
	fmt.Print("Additional time (min): ")
	fmt.Scan(&r.AdditionalTime)
	fmt.Print("Total time (min): ")
	fmt.Scan(&r.TotalTime)
	fmt.Print("Number of servings: ")
	fmt.Scan(&r.Servings)

	reader.ReadString('\n')

	// Ingredients
	fmt.Println("Enter ingredients (empty line to finish):")
	for {
		fmt.Print("> ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		r.Ingredients = append(r.Ingredients, line)
	}

	// Steps
	fmt.Println("Enter steps (empty line to finish):")
	for {
		fmt.Print("> ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		r.Steps = append(r.Steps, line)
	}

	fmt.Print("Filename to save (e.g., recipe.rfp): ")
	var filename string
	fmt.Scan(&filename)
	reader.ReadString('\n')
	path := filepath.Clean(filename)

	cfg, _ := rfp.LoadConfig()
	if err := rfp.WriteRecipe(cfg.DefaultRecipePath, path, r); err != nil {
		fmt.Println("Error writing recipe:", err)
		return
	}
	fmt.Println("Recipe saved to", path)
}

func readRecipe(reader *bufio.Reader) {
	fmt.Print("Recipe to read: ")
	var filename string
	filename, _ = reader.ReadString('\n')
	filename = strings.TrimSpace(filename)
	path := filepath.Clean(filename)
	cfg, _ := rfp.LoadConfig()

	r, err := rfp.ReadRecipeFile(path, cfg.DefaultRecipePath)
	if err != nil {
		fmt.Println("Error reading recipe:", err)
		return
	}

	fmt.Printf("Image Path: %s\n", r.ImagePath)
	fmt.Printf("Prep: %s min, Cook: %s min, Additional: %s min, Total: %s min\n",
		r.PrepTime, r.CookTime, r.AdditionalTime, r.TotalTime)
	fmt.Printf("Servings: %s\n", r.Servings)

	fmt.Println("\nIngredients:")
	for i, ing := range r.Ingredients {
		fmt.Printf("%d) %s\n", i+1, ing)
	}

	fmt.Println("\nSteps:")
	for i, step := range r.Steps {
		fmt.Printf("%d) %s\n", i+1, step)
	}
}

func editConfig() {
	// Try loading existing config
	cfg, err := rfp.LoadConfig()
	if err != nil {
		fmt.Println("Config not found, creating default config...")
		cfg, err = rfp.CreateConfig()
		if err != nil {
			fmt.Println("Failed to create default config:", err)
			return
		}
	}

	rfp.EditConfig(cfg)
}

func ScrapeAS() {
	config, _ := rfp.LoadConfig()

	reader := bufio.NewReader(os.Stdin)
	var url string
	fmt.Scan(&url)
	reader.ReadString('\n')

	recipe, err := ars.ScrapeAllRecipes(url, config.DefaultImagePath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Image Path:", recipe.ImagePath)
	fmt.Println("Prep Time:", recipe.PrepTime)
	fmt.Println("Cook Time:", recipe.CookTime)
	fmt.Println("Additional Time:", recipe.AdditionalTime)
	fmt.Println("Total Time:", recipe.TotalTime)
	fmt.Println("Servings:", recipe.Servings)

	fmt.Println("\nIngredients:")
	for _, ing := range recipe.Ingredients {
		fmt.Printf("%s\n", ing)
	}

	fmt.Println("\nSteps:")
	for i, step := range recipe.Steps {
		fmt.Printf("%d) %s\n", i+1, step)
	}

	rfp.WriteRecipe(config.DefaultRecipePath, recipe.Name, *recipe)
}
