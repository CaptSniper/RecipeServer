package rfp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DefaultRecipePath          string `json:"default_recipe_path"`
	DisplayIngredientsNumbered bool   `json:"display_ingredients_numbered"`
	DisplayStepsNumbered       bool   `json:"display_steps_numbered"`
	MaxIngredients             int    `json:"max_ingredients"`
	MaxSteps                   int    `json:"max_steps"`
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
		DefaultRecipePath:          "recipes/",
		DisplayIngredientsNumbered: true,
		DisplayStepsNumbered:       true,
		MaxIngredients:             50,
		MaxSteps:                   100,
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
