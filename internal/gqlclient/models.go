package gqlclient
//
// // Order -
// type Order struct {
// 	ID    int         `json:"id,omitempty"`
// 	Items []OrderItem `json:"items,omitempty"`
// }
//
// // OrderItem -
// type OrderItem struct {
// 	Coffee   Coffee `json:"coffee"`
// 	Quantity int    `json:"quantity"`
// }

// Project -
type Project struct {
	Name        string       `json:"name"`
	Slug      string       `json:"slug"`
	Description string       `json:"description"`
	ChangeSource  []ChangeSource `json:"changeSources"`
}

// Ingredient -
type ChangeSource struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
}
