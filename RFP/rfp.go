package rfp

// Recipe stores essential information needed for rendering
type Recipe struct {
	Name           string
	ImagePath      string   // Relative or absolute image path
	PrepTime       string   // Preparation time in minutes
	CookTime       string   // Cooking time in minutes
	AdditionalTime string   // Additional time in minutes (e.g., resting, marinating)
	TotalTime      string   // Total time in minutes
	Servings       string   // Number of servings
	Ingredients    []string // List of ingredients (text only for rendering)
	Steps          []string // Ordered list of preparation steps
}
