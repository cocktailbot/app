package search

// CategoryType that denotes one category
var CategoryType = "category"

// Categories to which recipes belong
type Categories struct {
	Data []Category `json:"data"`
}

// Category taxonomy for a recipe
type Category struct {
	ID       string `json:"id"`
	Slug     string `json:"slug"`
	Title    string `json:"title"`
	Children []struct {
		ID    string `json:"id"`
		Slug  string `json:"slug"`
		Title string `json:"title"`
	} `json:"children,omitempty"`
}

// GetData returns collection
func (cs Categories) GetData() []Category {
	return cs.Data
}

// GetID returns unique id
func (c Category) GetID() string {
	return c.ID
}
