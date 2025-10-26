package rfp

// Recipe stores essential information needed for rendering
type Recipe struct {
	ImagePath      string   // Relative or absolute image path
	PrepTimeMin    uint16   // Preparation time in minutes
	CookTimeMin    uint16   // Cooking time in minutes
	AdditionalTime uint16   // Additional time in minutes (e.g., resting, marinating)
	TotalTimeMin   uint16   // Total time in minutes
	Servings       uint16   // Number of servings
	Ingredients    []string // List of ingredients (text only for rendering)
	Steps          []string // Ordered list of preparation steps
}
