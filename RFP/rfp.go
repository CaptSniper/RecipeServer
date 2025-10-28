package rfp

// Recipe stores essential information needed for rendering
type Recipe struct {
	Name        string
	ImagePath   string
	CoreProps   map[string]string // e.g. {"Prep Time": "15 mins", "Servings": "6"}
	Ingredients []string
	Steps       []string
}

func NewRecipe() *Recipe {
	return &Recipe{
		CoreProps:   make(map[string]string),
		Ingredients: []string{},
		Steps:       []string{},
	}
}
