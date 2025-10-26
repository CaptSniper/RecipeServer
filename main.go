package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	rfp "github.com/CaptSniper/RecipeServer/RFP"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Recipe File Program")
	fmt.Println("-------------------")
	fmt.Println("Choose an option:")
	fmt.Println("1) Create a recipe file")
	fmt.Println("2) Read a recipe file")
	fmt.Println("3) Create/Edit config")
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

	if err := rfp.WriteRecipe(path, r); err != nil {
		fmt.Println("Error writing recipe:", err)
		return
	}
	fmt.Println("Recipe saved to", path)
}

func readRecipe(reader *bufio.Reader) {
	fmt.Print("Filename to read (e.g., recipe.rfp): ")
	var filename string
	fmt.Scan(&filename)
	path := filepath.Clean(filename)

	r, err := rfp.ReadRecipeFile(path)
	if err != nil {
		fmt.Println("Error reading recipe:", err)
		return
	}

	fmt.Println("\nRecipe Contents:")
	fmt.Printf("Image Path: %s\n", r.ImagePath)
	fmt.Printf("Prep: %d min, Cook: %d min, Additional: %d min, Total: %d min\n",
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

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Edit configuration (press Enter to keep current value):")

	// Helper functions
	promptString := func(field, current string) string {
		fmt.Printf("%s [%s]: ", field, current)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			return current
		}
		return input
	}

	promptInt := func(field string, current int) int {
		fmt.Printf("%s [%d]: ", field, current)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if val, err := strconv.Atoi(input); err == nil {
			return val
		}
		return current
	}

	promptBool := func(field string, current bool) bool {
		fmt.Printf("%s [%t] (true/false): ", field, current)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "true" {
			return true
		} else if input == "false" {
			return false
		}
		return current
	}

	// Prompt user for each field
	cfg.DefaultRecipePath = promptString("Default Recipe path", cfg.DefaultRecipePath)
	cfg.DisplayIngredientsNumbered = promptBool("Display ingredients numbered", cfg.DisplayIngredientsNumbered)
	cfg.DisplayStepsNumbered = promptBool("Display steps numbered", cfg.DisplayStepsNumbered)
	cfg.MaxIngredients = promptInt("Max ingredients", cfg.MaxIngredients)
	cfg.MaxSteps = promptInt("Max steps", cfg.MaxSteps)

	// Save updates
	if err := rfp.SaveConfig(cfg); err != nil {
		fmt.Println("Failed to save config:", err)
		return
	}

	fmt.Println("Configuration updated successfully.")
	fmt.Printf("%+v\n", *cfg)
}
