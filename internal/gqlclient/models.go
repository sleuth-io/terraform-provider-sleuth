package gqlclient

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

type ErrorImpactSource struct {
	Slug                       string      `json:"slug"`
	Environment                Environment `json:"environment"`
	Name                       string      `json:"name"`
	Provider                   string      `json:"provider"`
	ErrorOrgKey                string      `json:"errorOrgKey"`
	ErrorProjectKey            string      `json:"errorProjectKey"`
	ErrorEnvironment           string      `json:"errorEnvironment"`
	ManuallySetHealthThreshold float64     `json:"manuallySetHealthThreshold,omitempty"`
}

type MetricImpactSource struct {
	Slug                       string      `json:"slug"`
	Environment                Environment `json:"environment"`
	Name                       string      `json:"name"`
	Provider                   string      `json:"provider,omitempty"`
	Query                      string      `json:"query,omitempty"`
	LessIsBetter               bool        `json:"lessIsBetter,omitempty"`
	ManuallySetHealthThreshold float64     `json:"manuallySetHealthThreshold,omitempty"`
}

type Repository struct {
	Owner    string `json:"owner"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
	Url      string `json:"url,omitempty"`
}

type BranchMapping struct {
	EnvironmentSlug string `json:"environmentSlug"`
	Branch          string `json:"branch"`
}

type CodeChangeSource struct {
	Slug                        string                       `json:"slug"`
	Name                        string                       `json:"name"`
	Repository                  Repository                   `json:"repository"`
	DeployTrackingType          string                       `json:"deployTrackingType"`
	CollectImpact               bool                         `json:"collectImpact"`
	PathPrefix                  string                       `json:"pathPrefix"`
	NotifyInSlack               bool                         `json:"notifyInSlack"`
	IncludeInDashboard          bool                         `json:"includeInDashboard"`
	AutoTrackingDelay           int                          `json:"autoTrackingDelay"`
	EnvironmentMappings         []BranchMapping              `json:"environmentMappings"`
	DeployTrackingBuildMappings []DeployTrackingBuildMapping `json:"deployTrackingBuildMappings"`
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

type MutableErrorImpactSource struct {
	EnvironmentSlug            string  `json:"environment"`
	Name                       string  `json:"name"`
	Provider                   string  `json:"provider"`
	ErrorOrgKey                string  `json:"errorOrgKey"`
	ErrorProjectKey            string  `json:"errorProjectKey"`
	ErrorEnvironment           string  `json:"errorEnvironment"`
	ManuallySetHealthThreshold float64 `json:"manuallySetHealthThreshold,omitempty"`
}

type CreateErrorImpactSourceMutationInput struct {
	ProjectSlug string `json:"projectSlug"`
	*MutableErrorImpactSource
}

type UpdateErrorImpactSourceMutationInput struct {
	ProjectSlug string `json:"projectSlug"`
	Slug        string `json:"slug"`
	*MutableErrorImpactSource
}

type MutableMetricImpactSource struct {
	EnvironmentSlug            string  `json:"environment"`
	Name                       string  `json:"name"`
	Provider                   string  `json:"provider"`
	Query                      string  `json:"query,omitempty"`
	IntegrationSlug            string  `json:"auth,omitempty"`
	LessIsBetter               bool    `json:"lessIsBetter,omitempty"`
	ManuallySetHealthThreshold float64 `json:"manuallySetHealthThreshold,omitempty"`
}

type CreateMetricImpactSourceMutationInput struct {
	ProjectSlug string `json:"projectSlug"`
	*MutableMetricImpactSource
}

type UpdateMetricImpactSourceMutationInput struct {
	ProjectSlug string `json:"projectSlug"`
	Slug        string `json:"slug"`
	*MutableMetricImpactSource
}

type DeleteImpactSourceMutationInput struct {
	ProjectSlug string `json:"projectSlug"`
	Slug        string `json:"slug"`
}

type BuildMapping struct {
	EnvironmentSlug string `json:"environmentSlug"`
	Provider        string `json:"provider"`
	BuildName       string `json:"buildName"`
	JobName         string `json:"jobName,omitempty"`
	BuildProjectKey string `json:"buildProjectKey,omitempty"`
	IntegrationSlug string `json:"integrationSlug"`
	BuildBranch     string `json:"buildBranch"`
}

type DeployTrackingBuildMapping struct {
	Environment     Environment `json:"environment"`
	Provider        string      `json:"provider"`
	BuildName       string      `json:"buildName"`
	JobName         string      `json:"jobName,omitempty"`
	BuildProjectKey string      `json:"buildProjectKey"`
}

type MutableCodeChangeSource struct {
	Name                string          `json:"name"`
	Repository          Repository      `json:"repository"`
	DeployTrackingType  string          `json:"deployTrackingType"`
	CollectImpact       bool            `json:"collectImpact"`
	PathPrefix          string          `json:"pathPrefix"`
	NotifyInSlack       bool            `json:"notifyInSlack"`
	IncludeInDashboard  bool            `json:"includeInDashboard"`
	AutoTrackingDelay   int             `json:"autoTrackingDelay"`
	EnvironmentMappings []BranchMapping `json:"environmentMappings"`
	BuildMappings       []BuildMapping  `json:"buildMappings"`
}

type CreateCodeChangeSourceMutationInput struct {
	ProjectSlug       		string `json:"projectSlug"`
	InitializeChanges 		bool   `json:"initializeChanges"`
	InitializeChangesNow 	bool   `json:"initializeChangesNow"`
	*MutableCodeChangeSource
}

type UpdateCodeChangeSourceMutationInput struct {
	ProjectSlug string `json:"projectSlug"`
	Slug        string `json:"slug"`
	*MutableCodeChangeSource
}

type DeleteChangeSourceMutationInput struct {
	ProjectSlug string `json:"projectSlug"`
	Slug        string `json:"slug"`
}

type ErrorsType []struct {
	Field    string   `json:"field"`
	Messages []string `json:"messages"`
}
