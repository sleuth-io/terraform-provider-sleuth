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
	IssueTrackerProvider string	`json:"issueTrackerProvider"`
	BuildProvider string `json:"buildProvider"`
	ChangeFailureRateBoundary string `json:"changeFailureRateBoundary"`
	ImpactSensitivity	string	`json:"impactSensitivity"`
	FailureSensitivity int `json:"failureSensitivity"`
	//ChangeSource  []ChangeSource `json:"changeSources"`
}

// Ingredient -
type ChangeSource struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
}

type ProjectOptionalFields struct {
	Description               string `json:"description,omitempty"`
	IssueTrackerProvider      string `json:"issueTrackerProvider,omitempty"`
	BuildProvider             string `json:"buildProvider,omitempty"`
	ChangeFailureRateBoundary string `json:"changeFailureRateBoundary,omitempty"`
	ImpactSensitivity         string `json:"impactSensitivity,omitempty"`
	FailureSensitivity        int    `json:"failureSensitivity,omitempty"`
}

type CreateProjectMutationInput struct {
	Name string `json:"name"`
	*ProjectOptionalFields
}


type UpdateProjectMutationInput struct {
	Name string `json:"name,omitempty"`
	*ProjectOptionalFields
}


type ErrorsType []struct {
	Field string `json:"field"`
	Messages []string `json:"messages"`
}

