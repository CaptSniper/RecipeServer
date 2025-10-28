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

// Context carries shared info across handlers
type Context struct {
	Config  rfp.Config
	Running bool
	Reader  *bufio.Reader
}

type CommandHandler func(subcommand string, args []string, ctx *Context)
type Command struct {
	Handler        CommandHandler
	Help           string
	SubcommandHelp map[string]string
}

var commands = map[string]Command{}

func main() {
	cfg, _ := rfp.LoadConfig()
	ctx := &Context{Config: *cfg, Running: true, Reader: bufio.NewReader(os.Stdin)}

	// See below to add commands and subcommands
	registerCommands()

	// REPL loop
	for ctx.Running {
		fmt.Print("> ")
		line, _ := ctx.Reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		tokens := strings.Fields(line)
		cmd := strings.ToLower(tokens[0])
		sub := ""
		args := []string{}
		if len(tokens) > 1 {
			sub = tokens[1]
			if len(tokens) > 2 {
				args = tokens[2:]
			}
		}

		if handler, ok := commands[cmd]; ok {
			handler.Handler(sub, args, ctx)
		} else {
			fmt.Println("Unknown command:", cmd)
		}
	}
}

func handleHelp(sub string, args []string, ctx *Context) {
	if sub == "" {
		fmt.Println("Available commands:")
		for k, v := range commands {
			fmt.Printf("  %s: %s\n", k, v.Help)
		}
		fmt.Println("\nUse 'help <command>' for more details on a command.")
		return
	}

	if cmd, ok := commands[sub]; ok {
		fmt.Printf("Help for '%s':\n%s\n", sub, cmd.Help)
		if len(cmd.SubcommandHelp) > 0 {
			fmt.Println("Subcommands:")
			for subcmd, desc := range cmd.SubcommandHelp {
				fmt.Printf("  %s: %s\n", subcmd, desc)
			}
		}
	} else {
		fmt.Println("Unknown command for help:", sub)
	}
}

func handleNew(sub string, args []string, ctx *Context) {
	switch sub {
	case "recipe":
		createRecipe(ctx, args)
	case "config":
		editConfig(ctx, args)
	default:
		fmt.Println("Unknown subcommand for 'new':", sub)
	}
}

func handleRead(sub string, args []string, ctx *Context) {
	switch sub {
	case "recipe":
		readRecipe(ctx, args)
	default:
		fmt.Println("Unknown subcommand for 'read':", sub)
	}
}

func handleScrape(sub string, args []string, ctx *Context) {
	if sub != "recipe" {
		fmt.Println("Unknown subcommand for 'scrape':", sub)
		return
	}
	scrapeRecipe(ctx, args)
}

func handleList(sub string, args []string, ctx *Context) {
	switch sub {
	case "recipes":
		listRecipes(ctx)
	case "servers":
		listServers(ctx)
	default:
		fmt.Println("Unknown subcommand for 'list':", sub)
	}
}

func handleDelete(sub string, args []string, ctx *Context) {
	switch sub {
	case "recipe":
		deleteRecipe(ctx, args)
	case "server":
		fmt.Println("Server deletion not implemented yet")
	default:
		fmt.Println("Unknown subcommand for 'delete':", sub)
	}
}

func createRecipe(ctx *Context, args []string) {
	var r rfp.Recipe
	r.CoreProps = make(map[string]string)

	// Parse flags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-n", "--name":
			if i+1 < len(args) {
				r.Name = args[i+1]
				i++
			}
		case "-i", "--image":
			if i+1 < len(args) {
				r.ImagePath = args[i+1]
				i++
			}
		}
	}

	// Prompt if missing
	if r.Name == "" {
		fmt.Print("Recipe name: ")
		r.Name, _ = ctx.Reader.ReadString('\n')
		r.Name = strings.TrimSpace(r.Name)
	}
	if r.ImagePath == "" {
		fmt.Print("Image path: ")
		r.ImagePath, _ = ctx.Reader.ReadString('\n')
		r.ImagePath = strings.TrimSpace(r.ImagePath)
	}

	// Core properties
	fmt.Println("\nEnter recipe properties (e.g., Prep Time: 15 mins, Servings: 4).")
	fmt.Println("Press Enter on empty line to finish.")
	for {
		fmt.Print("> ")
		line, _ := ctx.Reader.ReadString('\n')
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
		line, _ := ctx.Reader.ReadString('\n')
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
		line, _ := ctx.Reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		r.Steps = append(r.Steps, line)
	}

	// Save
	filename := r.Name
	if err := rfp.WriteRecipe(ctx.Config.DefaultRecipePath, filename, r); err != nil {
		fmt.Println("Error writing recipe:", err)
		return
	}
	fmt.Println("Recipe saved as", filename)
}

func readRecipe(ctx *Context, args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: read recipe <filename> [-c] [-i] [-d]")
		return
	}

	filename := ""
	showCore := false
	showIngredients := false
	showSteps := false

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-c", "--core":
			showCore = true
		case "-i", "--ingredients":
			showIngredients = true
		case "-d", "--directions":
			showSteps = true
		default:
			if filename == "" {
				filename = args[i]
			}
		}
	}

	r, err := rfp.ReadRecipeFile(filename, ctx.Config.DefaultRecipePath)
	if err != nil {
		fmt.Println("Error reading recipe:", err)
		return
	}

	if showCore {
		fmt.Printf("Image Path: %s\n", r.ImagePath)
		fmt.Printf("Core Properties:\n")
		for k, v := range r.CoreProps {
			fmt.Printf("  %s: %s\n", k, v)
		}
	}

	if showIngredients {
		fmt.Println("\nIngredients:")
		for i, ing := range r.Ingredients {
			fmt.Printf("%d) %s\n", i+1, ing)
		}
	}

	if showSteps {
		fmt.Println("\nSteps:")
		for i, step := range r.Steps {
			fmt.Printf("%d) %s\n", i+1, step)
		}
	}
}

func editConfig(ctx *Context, args []string) {
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
	ctx.Config = *cfg
}

func scrapeRecipe(ctx *Context, args []string) {
	var url string
	for i := 0; i < len(args); i++ {
		if args[i] == "-u" || args[i] == "--url" {
			if i+1 < len(args) {
				url = args[i+1]
				i++
			}
		}
	}

	if url == "" {
		fmt.Println("No URL provided. Use -u <url>")
		return
	}

	recipe, err := ars.ScrapeAllRecipes(url, ctx.Config.DefaultImagePath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Scraped Recipe:", recipe.Name)
	if err := rfp.WriteRecipe(ctx.Config.DefaultRecipePath, recipe.Name, *recipe); err != nil {
		fmt.Println("Error writing recipe:", err)
		return
	}
	fmt.Println("Recipe saved as:", recipe.Name)
}

// Lists all recipe files in the configured recipe folder and prints their names
func listRecipes(ctx *Context) {
	cfg, err := rfp.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	files, err := os.ReadDir(cfg.DefaultRecipePath)
	if err != nil {
		fmt.Println("Error reading recipe folder:", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("No recipes found.")
		return
	}

	fmt.Println("Recipes:")
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		ext := filepath.Ext(name)
		if ext == ".rfp" {
			// Read the recipe file to get the actual Recipe.Name
			recipe, err := rfp.ReadRecipeFile(name, cfg.DefaultRecipePath)
			if err != nil {
				// Fallback to filename without extension
				fmt.Println(" -", name[:len(name)-len(ext)])
			} else {
				fmt.Println(" -", recipe.Name)
			}
		}
	}
}

// Deletes a recipe by its name
func deleteRecipe(ctx *Context, names []string) {
	if len(names) == 0 {
		fmt.Println("Please specify at least one recipe name")
		return
	}

	cfg, err := rfp.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	for _, name := range names {
		targetPath := filepath.Join(cfg.DefaultRecipePath, name+".rfp")
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			fmt.Println("Could not find", name)
			continue
		}

		if err := os.Remove(targetPath); err != nil {
			fmt.Println("Error deleting", name, ":", err)
			continue
		}

		fmt.Println("Deleted", name)
	}
}

// --- Placeholders for list/delete ---
func listServers(ctx *Context) { fmt.Println("List servers not implemented yet") }

func registerCommands() {
	// --- Main REPL registration ---
	commands["new"] = Command{
		Handler: handleNew,
		Help:    "Create new resources: recipe, config",
		SubcommandHelp: map[string]string{
			"recipe": "Create a new recipe.\n  Flags: -n, --name <name> | -i, --image <path>",
			"config": "Create or edit the configuration.\n  No flags",
		},
	}
	commands["read"] = Command{
		Handler: handleRead,
		Help:    "Read resources: recipe",
		SubcommandHelp: map[string]string{
			"recipe": "Read a recipe file.\n  Flags: -c, --core | -i, --ingredients | -d, --directions",
		},
	}
	commands["scrape"] = Command{
		Handler: handleScrape,
		Help:    "Scrape recipes from web",
		SubcommandHelp: map[string]string{
			"recipe": "Scrape a recipe from a URL.\n  Flags: -u, --url <url>",
		},
	}
	commands["list"] = Command{
		Handler: handleList,
		Help:    "List resources",
		SubcommandHelp: map[string]string{
			"recipes": "List all recipes",
			"servers": "List all running servers",
		},
	}
	commands["delete"] = Command{
		Handler: handleDelete,
		Help:    "Delete resources",
		SubcommandHelp: map[string]string{
			"recipe": "Delete a recipe by name.\n  Flags: -n, --name <name>",
			"server": "Delete a server (not implemented yet)",
		},
	}
	commands["exit"] = Command{
		Handler: func(_ string, _ []string, c *Context) { c.Running = false },
		Help:    "Exit the console",
	}
	commands["help"] = Command{
		Handler: handleHelp,
		Help:    "Show help for commands",
	}
}
