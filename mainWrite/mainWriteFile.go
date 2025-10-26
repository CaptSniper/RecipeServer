package mainWrite

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	rfp "github.com/CaptSniper/RecipeServer/RFP"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var recipe rfp.Recipe

	// --- Basic Info ---
	fmt.Print("Enter recipe image path: ")
	recipe.ImagePath = readLine(reader)

	recipe.PrepTimeMin = readUint16(reader, "Enter prep time (minutes): ")
	recipe.CookTimeMin = readUint16(reader, "Enter cook time (minutes): ")
	recipe.AdditionalTime = readUint16(reader, "Enter additional time (minutes): ")
	recipe.TotalTimeMin = readUint16(reader, "Enter total time (minutes): ")
	recipe.Servings = readUint16(reader, "Enter number of servings: ")

	// --- Ingredients ---
	fmt.Println("\nEnter ingredients one by one (empty line to finish):")
	for {
		fmt.Print("Ingredient: ")
		ing := readLine(reader)
		if ing == "" {
			break
		}
		recipe.Ingredients = append(recipe.Ingredients, ing)
	}

	// --- Steps ---
	fmt.Println("\nEnter steps one by one (empty line to finish):")
	for {
		fmt.Print("Step: ")
		step := readLine(reader)
		if step == "" {
			break
		}
		recipe.Steps = append(recipe.Steps, step)
	}

	// --- Output File ---
	fmt.Print("\nEnter output filename (e.g., myrecipe.rfp): ")
	filename := readLine(reader)
	if filename == "" {
		filename = "recipe.rfp"
	}

	// --- Write Recipe ---
	if err := rfp.WriteRecipe(filename, recipe); err != nil {
		log.Fatalf("Error writing recipe: %v", err)
	}

	fmt.Printf("Recipe successfully written to %s\n", filename)
}

// readLine reads a line of input from stdin
func readLine(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

// readUint16 reads a uint16 value from stdin with prompt
func readUint16(reader *bufio.Reader, prompt string) uint16 {
	for {
		fmt.Print(prompt)
		text := readLine(reader)
		val, err := strconv.ParseUint(text, 10, 16)
		if err == nil {
			return uint16(val)
		}
		fmt.Println("Invalid number. Please enter a valid integer.")
	}
}
