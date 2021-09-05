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
}

type Environment struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Description string       `json:"description"`
	Color     string `json:"color"`
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

type EnvironmentOptionalFields struct {
	Description               string `json:"description,omitempty"`
	Color				      string `json:"color,omitempty"`
}

type CreateEnvironmentMutationInput struct {
	Name string `json:"name"`
	*EnvironmentOptionalFields
}


type UpdateEnvironmentMutationInput struct {
	Name string `json:"name,omitempty"`
	*EnvironmentOptionalFields
}


type ErrorsType []struct {
	Field string `json:"field"`
	Messages []string `json:"messages"`
}

