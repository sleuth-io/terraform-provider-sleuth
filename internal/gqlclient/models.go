package gqlclient

// Project -
type Project struct {
	Slug                      string `json:"slug"`
	Name                      string `json:"name"`
	Description               string `json:"description,omitempty"`
	IssueTrackerProvider      string `json:"issueTrackerProvider,omitempty"`
	BuildProvider             string `json:"buildProvider,omitempty"`
	ChangeFailureRateBoundary string `json:"changeFailureRateBoundary,omitempty"`
	ImpactSensitivity         string `json:"impactSensitivity,omitempty"`
	FailureSensitivity        int    `json:"failureSensitivity,omitempty"`
}

type Environment struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
}

type MutableProject struct {
	Name                      string `json:"name"`
	Description               string `json:"description,omitempty"`
	IssueTrackerProvider      string `json:"issueTrackerProvider,omitempty"`
	BuildProvider             string `json:"buildProvider,omitempty"`
	ChangeFailureRateBoundary string `json:"changeFailureRateBoundary,omitempty"`
	ImpactSensitivity         string `json:"impactSensitivity,omitempty"`
	FailureSensitivity        int    `json:"failureSensitivity,omitempty"`
}

type CreateProjectMutationInput struct {
	*MutableProject
}

type UpdateProjectMutationInput struct {
	Slug string `json:"slug"`
	*MutableProject
}

type DeleteProjectMutationInput struct {
	Slug string `json:"slug"`
}

type MutableEnvironment struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
}

type CreateEnvironmentMutationInput struct {
	ProjectSlug string `json:"projectSlug"`
	*MutableEnvironment
}

type UpdateEnvironmentMutationInput struct {
	ProjectSlug string `json:"projectSlug"`
	Slug        string `json:"slug"`
	*MutableEnvironment
}

type DeleteEnvironmentMutationInput struct {
	ProjectSlug string `json:"projectSlug"`
	Slug        string `json:"slug"`
}

type ErrorsType []struct {
	Field    string   `json:"field"`
	Messages []string `json:"messages"`
}
