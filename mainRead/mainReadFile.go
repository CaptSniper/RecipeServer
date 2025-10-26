package mainRead

import (
	"fmt"
	"log"

	"../rfp"
)

func main() {
	// File to read
	filename := "recipe.rfp"

	// Read the recipe using the library
	recipe, err := rfp.ReadRecipeFile(filename)
	if err != nil {
		log.Fatalf("Error reading recipe file: %v", err)
	}

	// Print basic recipe information
	fmt.Println("=== Recipe Information ===")
	fmt.Printf("Image Path      : %s\n", recipe.ImagePath)
	fmt.Printf("Prep Time       : %d min\n", recipe.PrepTimeMin)
	fmt.Printf("Cook Time       : %d min\n", recipe.CookTimeMin)
	fmt.Printf("Additional Time : %d min\n", recipe.AdditionalTime)
	fmt.Printf("Total Time      : %d min\n", recipe.TotalTimeMin)
	fmt.Printf("Servings        : %d\n", recipe.Servings)

	// Print ingredients
	fmt.Println("\nIngredients:")
	for i, ing := range recipe.Ingredients {
		fmt.Printf("%d. %s\n", i+1, ing)
	}

	// Print steps
	fmt.Println("\nSteps:")
	for i, step := range recipe.Steps {
		fmt.Printf("%d. %s\n", i+1, step)
	}
}
