package models

// Category taxonomy for a recipe
type Category struct {
	ID       string     `json:"id"`
	Slug     string     `json:"slug"`
	Title    string     `json:"title"`
	Children []Category `json:"children,omitempty"`
}
