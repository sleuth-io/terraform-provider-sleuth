package gqlclient

type CLTStartStates struct {
	ID string `json:"id,omitempty"`
}

type Project struct {
	Slug                      string           `json:"slug"`
	Name                      string           `json:"name"`
	Description               string           `json:"description,omitempty"`
	IssueTrackerProvider      string           `json:"issueTrackerProvider,omitempty"`
	BuildProvider             string           `json:"buildProvider,omitempty"`
	ChangeFailureRateBoundary string           `json:"changeFailureRateBoundary,omitempty"`
	ImpactSensitivity         string           `json:"impactSensitivity,omitempty"`
	FailureSensitivity        int              `json:"failureSensitivity,omitempty"`
	CltStartDefinition        string           `json:"cltStartDefinition,omitempty"`
	CltStartStates            []CLTStartStates `json:"cltStartStates,omitempty"`
	StrictIssueMatching       bool             `json:"strictIssueMatching,omitempty"`
	LabelNames                []string         `json:"labelNames"`
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
	ManuallySetHealthThreshold *float64    `json:"manuallySetHealthThreshold,omitempty"`
	IntegrationAuthSlug        string      `json:"integrationAuthSlug,omitempty"`
}

type MetricImpactSource struct {
	Slug                       string      `json:"slug"`
	Environment                Environment `json:"environment"`
	Name                       string      `json:"name"`
	Provider                   string      `json:"provider,omitempty"`
	Query                      string      `json:"query,omitempty"`
	IntegrationAuthSlug        string      `json:"integrationAuthSlug,omitempty"`
	LessIsBetter               bool        `json:"lessIsBetter,omitempty"`
	ManuallySetHealthThreshold *float64    `json:"manuallySetHealthThreshold,omitempty"`
}

type RepositoryBase struct {
	Owner      string `json:"owner"`
	Name       string `json:"name"`
	Provider   string `json:"provider"`
	Url        string `json:"url,omitempty"`
	ProjectUID string `json:"projectUid,omitempty"`
	RepoUID    string `json:"repoUid,omitempty"`
}

type Repository struct {
	RepositoryBase
	IntegrationAuth *IntegrationAuth `json:"integrationAuth,omitempty"`
	Webhook         *Webhook         `json:"webhook,omitempty"`
}

type Webhook struct {
	Url    string `json:"url"`
	Secret string `json:"secret,omitempty"`
}

type IntegrationAuth struct {
	Slug string `json:"slug"`
}

type MutableRepository struct {
	RepositoryBase
	IntegrationSlug string `json:"integrationSlug,omitempty"`
}

type BranchMapping struct {
	EnvironmentSlug string `json:"environmentSlug" tfsdk:"environment_slug"`
	Branch          string `json:"branch" tfsdk:"branch"`
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
	Name                      string   `json:"name"`
	Description               string   `json:"description,omitempty"`
	IssueTrackerProvider      string   `json:"issueTrackerProvider,omitempty"`
	BuildProvider             string   `json:"buildProvider,omitempty"`
	ChangeFailureRateBoundary string   `json:"changeFailureRateBoundary,omitempty"`
	ImpactSensitivity         string   `json:"impactSensitivity,omitempty"`
	FailureSensitivity        int      `json:"failureSensitivity,omitempty"`
	CltStartDefinition        string   `json:"cltStartDefinition,omitempty"`
	CltStartStates            []int    `json:"cltStartStates,omitempty"`
	StrictIssueMatching       bool     `json:"strictIssueMatching,omitempty"`
	Labels                    []string `json:"labels"`
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
	IntegrationSlug            string  `json:"auth,omitempty"`
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

// This represents a build mapping for creation or mutation
type BuildMapping struct {
	EnvironmentSlug          string `json:"environmentSlug"`
	Provider                 string `json:"provider"`
	BuildName                string `json:"buildName"`
	JobName                  string `json:"jobName,omitempty"`
	BuildProjectKey          string `json:"buildProjectKey,omitempty"`
	BuildProjectName         string `json:"buildProjectName,omitempty"`
	IntegrationSlug          string `json:"integrationSlug"`
	BuildBranch              string `json:"buildBranch"`
	MatchBranchToEnvironment bool   `json:"matchBranchToEnvironment,omitempty"`
	IsCustom                 bool   `json:"isCustom,omitempty"`
}

// This represents the build mapping as retrieved from a query
type DeployTrackingBuildMapping struct {
	Environment              Environment `json:"environment"`
	Provider                 string      `json:"provider"`
	IntegrationSlug          string      `json:"integrationSlug"`
	BuildName                string      `json:"buildName"`
	JobName                  string      `json:"jobName,omitempty"`
	BuildProjectKey          string      `json:"buildProjectKey"`
	MatchBranchToEnvironment bool        `json:"matchBranchToEnvironment"`
	IsCustom                 bool        `json:"isCustom"`
}

type MutableCodeChangeSource struct {
	Name                string            `json:"name"`
	Repository          MutableRepository `json:"repository"`
	DeployTrackingType  string            `json:"deployTrackingType"`
	CollectImpact       bool              `json:"collectImpact"`
	PathPrefix          string            `json:"pathPrefix"`
	NotifyInSlack       bool              `json:"notifyInSlack"`
	IncludeInDashboard  bool              `json:"includeInDashboard"`
	AutoTrackingDelay   int               `json:"autoTrackingDelay"`
	EnvironmentMappings []BranchMapping   `json:"environmentMappings"`
	BuildMappings       []BuildMapping    `json:"buildMappings"`
}

type CreateCodeChangeSourceMutationInput struct {
	ProjectSlug       string `json:"projectSlug"`
	InitializeChanges bool   `json:"initializeChanges"`
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

type PagerDutyProviderData struct {
	RemoteServices string `json:"remoteServices"`
	RemoteUrgency  string `json:"remoteUrgency"`
}

type DataDogProviderData struct {
	Query                   string `json:"query"`
	RemotePriorityThreshold string `json:"remotePriorityThreshold"`
}

type JiraProviderData struct {
	RemoteJql string `json:"remoteJql"`
}

type BlamelessProviderData struct {
	RemoteTypes             []string `json:"remoteTypes"`
	RemoteSeverityThreshold string   `json:"remoteSeverityThreshold"`
}

type StatuspageProviderData struct {
	RemotePage                 string `json:"remotePage"`
	RemoteComponent            string `json:"remoteComponent"`
	RemoteImpact               string `json:"remoteImpact"`
	IgnoreMaintenanceIncidents bool   `json:"ignoreMaintenanceIncidents"`
}

type OpsGenieProviderData struct {
	RemoteAlertTags         string `json:"remoteAlertTags,omitempty"`
	RemoteIncidentTags      string `json:"remoteIncidentTags,omitempty"`
	RemotePriorityThreshold string `json:"remotePriorityThreshold"`
	RemoteService           string `json:"remoteService,omitempty"`
	RemoteUseAlerts         bool   `json:"remoteUseAlerts"`
}

type FireHydrantProviderData struct {
	RemoteEnvironments       string `json:"remoteEnvironments"`
	RemoteServices           string `json:"remoteServices"`
	RemoteMitigatedIsHealthy bool   `json:"remoteMitigatedIsHealthy"`
}

type ClubhouseProviderData struct {
	RemoteQuery string `json:"remoteQuery"`
}

type RootlyProviderData struct {
	RemoteSeverity     string `json:"remoteSeverity"`
	RemoteIncidentType string `json:"remoteIncidentType"`
	RemoteEnvironment  string `json:"remoteEnvironment"`
	RemoteService      string `json:"remoteService"`
	RemoteTeam         string `json:"remoteTeam"`
}

type ProviderData struct {
	PagerDutyProviderData   PagerDutyProviderData   `json:"pagerDutyProviderData" graphql:"... on PagerDutyProviderData"`
	DataDogProviderData     DataDogProviderData     `json:"dataDogProviderData" graphql:"... on DataDogProviderData"`
	JiraProviderData        JiraProviderData        `json:"jiraProviderData" graphql:"... on JiraProviderData"`
	BlamelessProviderData   BlamelessProviderData   `json:"blamelessProviderData" graphql:"... on BlamelessProviderData"`
	StatuspageProviderData  StatuspageProviderData  `json:"statuspageProviderData" graphql:"... on StatuspageProviderData"`
	OpsGenieProviderData    OpsGenieProviderData    `json:"opsgenieProviderData" graphql:"... on OpsgenieProviderData"`
	FireHydrantProviderData FireHydrantProviderData `json:"firehydrantProviderData" graphql:"... on FireHydrantProviderData"`
	ClubhouseProviderData   ClubhouseProviderData   `json:"ClubhouseProviderData" graphql:"... on ClubhouseProviderData"`
	RootlyProviderData      RootlyProviderData      `json:"RootlyProviderData" graphql:"... on RootlyProviderData"`
}

type IncidentImpactSource struct {
	Slug                string       `json:"slug"`
	Environment         Environment  `json:"environment"`
	Name                string       `json:"name"`
	Provider            string       `json:"provider"`
	ProviderData        ProviderData `json:"providerData"`
	IntegrationAuthSlug string       `json:"integrationAuthSlug"`
	RegisterImpactLink  string       `json:"registerImpactLink"`
}

type PagerDutyInputType struct {
	RemoteServices string `json:"remoteServices"`
	RemoteUrgency  string `json:"remoteUrgency"`
}

type DataDogInputType struct {
	DataDogProviderData
	IntegrationSlug string `json:"integrationSlug"`
}

type JiraInputType struct {
	JiraProviderData
	IntegrationSlug string `json:"integrationSlug"`
}

type BlamelessInputType struct {
	BlamelessProviderData
	IntegrationSlug string `json:"integrationSlug"`
}

type StatuspageInputType struct {
	StatuspageProviderData
	IntegrationSlug string `json:"integrationSlug"`
}
type FireHydrantInputType struct {
	FireHydrantProviderData
	IntegrationSlug string `json:"integrationSlug"`
}

type ClubhouseInputType struct {
	ClubhouseProviderData
	IntegrationSlug string `json:"integrationSlug"`
}

type RootlyInputType struct {
	RootlyProviderData
	IntegrationSlug string `json:"integrationSlug"`
}

type OpsGenieInputType struct {
	OpsGenieProviderData
	IntegrationSlug string `json:"integrationSlug"`
}

type IncidentImpactSourceInputType struct {
	ProjectSlug          string                `json:"projectSlug"`
	EnvironmentName      string                `json:"environmentName"`
	Name                 string                `json:"name"`
	Provider             string                `json:"provider"`
	PagerDutyInputType   *PagerDutyInputType   `json:"pagerDutyInput"`
	DataDogInputType     *DataDogInputType     `json:"datadogInput"`
	JiraInputType        *JiraInputType        `json:"jiraInput"`
	BlamelessInputType   *BlamelessInputType   `json:"blamelessInput"`
	StatuspageInputType  *StatuspageInputType  `json:"statuspageInput"`
	OpsGenieInputType    *OpsGenieInputType    `json:"opsgenieInput"`
	FireHydrantInputType *FireHydrantInputType `json:"firehydrantInput"`
	ClubhouseInputType   *ClubhouseInputType   `json:"clubhouseInput"`
	RootlyInputType      *RootlyInputType      `json:"rootlyInput"`
}

type IncidentImpactSourceInputUpdateType struct {
	IncidentImpactSourceInputType
	Slug string `json:"slug"`
}

type IncidentImpactSourceDeleteInputType struct {
	Slug        string `json:"slug"`
	ProjectSlug string `json:"projectSlug"`
}

// Team models and mutation inputs/outputs

type Team struct {
	ID     string `json:"id"`
	Slug   string `json:"slug"`
	Name   string `json:"name"`
	Parent *struct {
		Slug string `json:"slug"`
	} `json:"parent"`
}

type MutableTeam struct {
	Name   string  `json:"name"`
	Parent *string `json:"parent,omitempty"`
}

type CreateTeamMutationInput struct {
	Name   string  `json:"name" graphql:"name"`
	Parent *string `json:"parent,omitempty" graphql:"parent"`
}

type UpdateTeamMutationInput struct {
	Slug   string  `json:"slug" graphql:"slug"`
	Name   *string `json:"name,omitempty" graphql:"name"`
	Parent *string `json:"parent,omitempty" graphql:"parent"`
}

type DeleteTeamMutationInput struct {
	Slug string `json:"slug" graphql:"slug"`
}

type AddTeamMembersMutationInput struct {
	Slug    string   `json:"slug"`
	Members []string `json:"members"`
}

type RemoveTeamMembersMutationInput struct {
	Slug    string   `json:"slug"`
	Members []string `json:"members"`
}

// Payloads for mutation responses

type CreateTeamMutationPayload struct {
	Errors ErrorsType `json:"errors"`
	Team   *Team      `json:"team"`
}

type UpdateTeamMutationPayload struct {
	Errors ErrorsType `json:"errors"`
	Team   *Team      `json:"team"`
}

type DeleteTeamMutationPayload struct {
	Errors  ErrorsType `json:"errors"`
	Success bool       `json:"success"`
}

type AddTeamMembersMutationPayload struct {
	Errors  ErrorsType `json:"errors"`
	Success bool       `json:"success"`
}

type RemoveTeamMembersMutationPayload struct {
	Errors  ErrorsType `json:"errors"`
	Success bool       `json:"success"`
}

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	IsActive  bool   `json:"isActive"`
}
