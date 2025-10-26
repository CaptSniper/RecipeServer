package rfp

// Recipe stores essential information needed for rendering
type Recipe struct {
	ImagePath      string   // Relative or absolute image path
	PrepTime       uint16   // Preparation time in minutes
	CookTime       uint16   // Cooking time in minutes
	AdditionalTime uint16   // Additional time in minutes (e.g., resting, marinating)
	TotalTime      uint16   // Total time in minutes
	Servings       string   // Number of servings
	Ingredients    []string // List of ingredients (text only for rendering)
	Steps          []string // Ordered list of preparation steps
}
