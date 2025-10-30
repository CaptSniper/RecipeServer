package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	rfp "github.com/CaptSniper/RecipeServer/RFP"
	ars "github.com/CaptSniper/RecipeServer/webScraper"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Recipe File Program")
		fmt.Println("-------------------")
		fmt.Println("Choose an option:")
		fmt.Println("1) Create a recipe file")
		fmt.Println("2) Read a recipe file")
		fmt.Println("3) Create/Edit config")
		fmt.Println("4) Scrape AllRecipes")
		fmt.Println("5) Start API Server")
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
		case 5:
			go StartApiServer()
		case 6:
			go StartWebServer()
		default:
			fmt.Println("Unknown option")
		}
	}
}

func createRecipe(reader *bufio.Reader) {
	var r rfp.Recipe
	r.CoreProps = make(map[string]string)

	// Name (required)
	fmt.Print("Recipe name: ")
	r.Name, _ = reader.ReadString('\n')
	r.Name = strings.TrimSpace(r.Name)

	// Image path
	fmt.Print("Image path: ")
	r.ImagePath, _ = reader.ReadString('\n')
	r.ImagePath = strings.TrimSpace(r.ImagePath)

	// Core properties
	fmt.Println("\nEnter recipe properties (e.g. Prep Time: 15 mins, Servings: 4).")
	fmt.Println("Press Enter on an empty line when done.")
	for {
		fmt.Print("> ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			fmt.Println("Invalid format. Use 'Key: Value'")
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		r.CoreProps[key] = value
	}

	// Ingredients
	fmt.Println("\nEnter ingredients (empty line to finish):")
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
	fmt.Println("\nEnter steps (empty line to finish):")
	for {
		fmt.Print("> ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		r.Steps = append(r.Steps, line)
	}

	// Save
	fmt.Print("\nFilename to save (e.g., recipe.rfp): ")
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
		r.CoreProps["prep time"], r.CoreProps["cook time"], r.CoreProps["additional time"], r.CoreProps["total time"])
	fmt.Printf("Servings: %s\n", r.CoreProps["servings"])

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
	fmt.Println("Prep Time:", recipe.CoreProps["prep time"])
	fmt.Println("Cook Time:", recipe.CoreProps["cook time"])
	fmt.Println("Additional Time:", recipe.CoreProps["additional time"])
	fmt.Println("Total Time:", recipe.CoreProps["total time"])
	fmt.Println("Servings:", recipe.CoreProps["servings"])

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
