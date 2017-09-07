package search

// CategoryType that denotes one category
var CategoryType = "category"

// Categories to which recipes belong
type Categories struct {
	BaseCollection
	Data []Category `json:"data"`
}

// Category taxonomy for a recipe
type Category struct {
	BaseIndexable
	Slug     string `json:"slug"`
	Title    string `json:"title"`
	Children []struct {
		ID    string `json:"id"`
		Slug  string `json:"slug"`
		Title string `json:"title"`
	} `json:"children,omitempty"`
}
