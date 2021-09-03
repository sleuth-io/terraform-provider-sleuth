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
	FailureSensitivity int `json:"failureSensitivity"`
	//ChangeSource  []ChangeSource `json:"changeSources"`
}

// Ingredient -
type ChangeSource struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
}

type ProjectCreationMutationInput struct {
	Name string `json:"name"`
	FailureSensitivity string `json:"failureSensitivity,omitempty"`
}


type ProjectUpdateMutationInput struct {
	Name string `json:"name,omitempty"`
	FailureSensitivity int `json:"failureSensitivity,omitempty"`
}


type ErrorsType []struct {
	Field string `json:"field"`
	Messages []string `json:"messages"`
}

