package rfp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	DefaultImagePath           string `json:"default_image_path"`
	DefaultRecipePath          string `json:"default_recipe_path"`
	DisplayIngredientsNumbered bool   `json:"display_ingredients_numbered"`
	DisplayStepsNumbered       bool   `json:"display_steps_numbered"`
	MaxIngredients             int    `json:"max_ingredients"`
	MaxSteps                   int    `json:"max_steps"`
	DefaultPort                int    `json:"default_port"`
}

// LoadConfig reads the JSON config from .config/config.json relative to the repo root
func LoadConfig() (*Config, error) {
	configPath := filepath.Join(".config", "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SaveConfig writes the given Config struct to .config/config.json
func SaveConfig(cfg *Config) error {
	configDir := ".config"
	configPath := filepath.Join(configDir, "config.json")

	// Create .config folder if it doesn't exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.Mkdir(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %v", err)
		}
	}

	// Create recipe folder if doesn't exist
	if _, err := os.Stat(cfg.DefaultRecipePath); os.IsNotExist(err) {
		if err := os.Mkdir(cfg.DefaultRecipePath, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %v", err)
		}
	}

	// Create images folder if doesn't exist
	if _, err := os.Stat(cfg.DefaultImagePath); os.IsNotExist(err) {
		if err := os.Mkdir(cfg.DefaultImagePath, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %v", err)
		}
	}

	// Marshal the config to JSON
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// CreateConfig creates a default config file at .config/config.json
func CreateConfig() (*Config, error) {
	configDir := ".config"
	configPath := filepath.Join(configDir, "config.json")

	// Create .config folder if it doesn't exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.Mkdir(configDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create config directory: %v", err)
		}
	}

	// Default configuration
	cfg := &Config{
		DefaultImagePath:           "images/",
		DefaultRecipePath:          "recipes/",
		DisplayIngredientsNumbered: true,
		DisplayStepsNumbered:       true,
		MaxIngredients:             50,
		MaxSteps:                   100,
		DefaultPort:                26740,
	}

	// Write to config.json
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal default config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write config file: %v", err)
	}

	return cfg, nil
}

func EditConfig(cfg *Config) {
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
	cfg.DefaultImagePath = promptString("Default Image path", cfg.DefaultImagePath)
	cfg.DefaultRecipePath = promptString("Default Recipe path", cfg.DefaultRecipePath)
	cfg.DisplayIngredientsNumbered = promptBool("Display ingredients numbered", cfg.DisplayIngredientsNumbered)
	cfg.DisplayStepsNumbered = promptBool("Display steps numbered", cfg.DisplayStepsNumbered)
	cfg.MaxIngredients = promptInt("Max ingredients", cfg.MaxIngredients)
	cfg.MaxSteps = promptInt("Max steps", cfg.MaxSteps)
	cfg.DefaultPort = promptInt("Default Port", cfg.DefaultPort)

	// Save updates
	if err := SaveConfig(cfg); err != nil {
		fmt.Println("Failed to save config:", err)
		return
	}

	fmt.Println("Configuration updated successfully.")
	fmt.Printf("%+v\n", *cfg)
}
