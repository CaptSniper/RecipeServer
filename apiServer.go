package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	rfp "github.com/CaptSniper/RecipeServer/RFP"
	ars "github.com/CaptSniper/RecipeServer/webScraper"
	"github.com/gorilla/mux"
)

// --- Handlers ---

type RecipeSummary struct {
	ID   string `json:"id"`   // filename without .rfp
	Name string `json:"name"` // recipe.Name from the file
}

type ScrapeRequest struct {
	URL  string `json:"url"`
	Save bool   `json:"save"`
}

func StartApiServer() {
	cfg, err := rfp.LoadConfig()
	if err != nil {
		fmt.Println("Failed to load config. Try running option 3 to create a default config:", err)
		return
	}
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/recipes", listRecipesHandler).Methods("GET")
	r.HandleFunc("/recipes/{id}", getRecipeHandler).Methods("GET")
	r.HandleFunc("/recipes", createRecipeHandler).Methods("POST")
	r.HandleFunc("/recipes/{id}", updateRecipeHandler).Methods("PUT")
	r.HandleFunc("/recipes/{id}", deleteRecipeHandler).Methods("DELETE")
	r.HandleFunc("/scrape", scrapeRecipeHandler).Methods("POST")

	fmt.Println("Server running at http://localhost:" + strconv.Itoa(cfg.DefaultPort))
	http.ListenAndServe(":"+strconv.Itoa(cfg.DefaultPort), r)
}

// listRecipesHandler – lists all recipes by name and ID
func listRecipesHandler(w http.ResponseWriter, r *http.Request) {
	cfg, err := rfp.LoadConfig()
	if err != nil {
		http.Error(w, "Failed to load config", http.StatusInternalServerError)
		return
	}
	files, err := os.ReadDir(cfg.DefaultRecipePath)
	if err != nil {
		http.Error(w, "Failed to read recipe directory", http.StatusInternalServerError)
		return
	}

	var recipes []map[string]string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".rfp" {
			recipe, err := rfp.ReadRecipeFile(cfg.DefaultRecipePath, file.Name())
			if err != nil {
				continue // skip corrupted files
			}
			id := strings.TrimSuffix(file.Name(), ".rfp")
			recipes = append(recipes, map[string]string{
				"id":   id,
				"name": recipe.Name,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

// getRecipeHandler – gets a specific recipe by ID
func getRecipeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Missing recipe ID", http.StatusBadRequest)
		return
	}

	cfg, err := rfp.LoadConfig()
	if err != nil {
		http.Error(w, "Failed to load config", http.StatusInternalServerError)
		return
	}

	recipe, err := rfp.ReadRecipeFile(cfg.DefaultRecipePath, id+".rfp")
	if err != nil {
		http.Error(w, "Failed to read recipe: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

// createRecipeHandler – creates a new recipe
func createRecipeHandler(w http.ResponseWriter, r *http.Request) {
	var recipe rfp.Recipe
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &recipe); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if recipe.Name == "" {
		http.Error(w, "Recipe name is required", http.StatusBadRequest)
		return
	}

	cfg, err := rfp.LoadConfig()
	if err != nil {
		http.Error(w, "Failed to load config", http.StatusInternalServerError)
		return
	}
	id := strings.ReplaceAll(strings.ToLower(recipe.Name), " ", "_")

	if err := rfp.WriteRecipe(cfg.DefaultRecipePath, id+".rfp", recipe); err != nil {
		http.Error(w, "Failed to save recipe: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Recipe created successfully",
		"id":      id,
	})
}

// updateRecipeHandler – updates an existing recipe
func updateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing recipe ID", http.StatusBadRequest)
		return
	}

	cfg, err := rfp.LoadConfig()
	if err != nil {
		http.Error(w, "Failed to load config", http.StatusInternalServerError)
		return
	}
	recipePath := filepath.Join(cfg.DefaultRecipePath, id+".rfp")

	if _, err := os.Stat(recipePath); os.IsNotExist(err) {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}

	var updated rfp.Recipe
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := rfp.WriteRecipe(cfg.DefaultRecipePath, id+".rfp", updated); err != nil {
		http.Error(w, "Failed to update recipe: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Recipe updated successfully",
		"id":      id,
	})
}

// deleteRecipeHandler – deletes a recipe
func deleteRecipeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Missing recipe ID", http.StatusBadRequest)
		return
	}

	cfg, err := rfp.LoadConfig()
	if err != nil {
		http.Error(w, "Failed to load config", http.StatusInternalServerError)
		return
	}
	recipePath := filepath.Join(cfg.DefaultRecipePath, id+".rfp")

	if err := os.Remove(recipePath); err != nil {
		http.Error(w, "Failed to delete recipe: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Recipe deleted successfully",
		"id":      id,
	})
}

func scrapeRecipeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ScrapeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Determine default image path from config
	config, err := rfp.LoadConfig()
	if err != nil {
		http.Error(w, "Failed to load config", http.StatusInternalServerError)
		return
	}
	imagePath := config.DefaultImagePath

	// Scrape the recipe
	recipe, err := ars.ScrapeRecipe(req.URL, imagePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to scrape recipe: %v", err), http.StatusBadRequest)
		return
	}

	// Optionally save the recipe immediately
	if req.Save {
		if err := rfp.WriteRecipe(config.DefaultRecipePath, recipe.Name, *recipe); err != nil {
			http.Error(w, fmt.Sprintf("Failed to save recipe: %v", err), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}
