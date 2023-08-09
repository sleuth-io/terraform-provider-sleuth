package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sleuth-io/terraform-provider-sleuth/internal/gqlclient"
)

type ImpactProvider int

const (
	Unknown ImpactProvider = iota
	PagerDuty
	DataDog
	Jira
	CustomIncident
	Blameless
)

func ImpactProviderFromString(s string) ImpactProvider {
	switch strings.ToUpper(s) {
	case "PAGERDUTY":
		return PagerDuty
	case "DATADOG":
		return DataDog
	case "JIRA":
		return Jira
	case "CUSTOM_INCIDENT":
		return CustomIncident
	case "BLAMELESS":
		return Blameless
	}
	return Unknown
}

func (s ImpactProvider) String() string {
	switch s {
	case PagerDuty:
		return "PAGERDUTY"
	case DataDog:
		return "DATADOG"
	case Jira:
		return "JIRA"
	case CustomIncident:
		return "CUSTOM_INCIDENT"
	case Blameless:
		return "BLAMELESS"
	}
	return "unknown"
}

func resourceIncidentImpactSource() *schema.Resource {
	return &schema.Resource{
		Description: "Sleuth incident impact source",

		CreateContext: resourceIncidentImpactSourceCreate,
		ReadContext:   resourceIncidentImpactSourceRead,
		UpdateContext: resourceIncidentImpactSourceUpdate,
		DeleteContext: resourceIncidentImpactSourceDelete,

		Schema: map[string]*schema.Schema{
			"project_slug": {
				Description: "Project slug",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "Impact source name",
				Type:        schema.TypeString,
				Required:    true,
			},
			// can't use `provider` because terraform tries to import provider
			"provider_name": {
				Description: "Impact source provider (options: PAGERDUTY, CUSTOM_INCIDENT)",
				Type:        schema.TypeString,
				Required:    true,
			},
			"environment_name": {
				Description: "Impact source environment name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"register_impact_link": {
				Description: "Impact source webhook registration link (for CUSTOM_INCIDENT only)",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"pagerduty_input": {
				Description: "PagerDuty input",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"remote_services": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "List of remote services, empty string means all",
							Default:     "",
						},
						"remote_urgency": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "PagerDuty remote urgency, options: HIGH, LOW, ANY",
						},
						"historic_init": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Populate with data from the last 30 days",
						},
						"integration_slug": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IntegrationAuthentication slug used",
						},
					},
				},
			},
			"datadog_input": {
				Description: "DataDog input",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query": {
							Type:     schema.TypeString,
							Optional: true,
							Description: `The query to scope the monitors to track. If you are using a custom facet you would need to add @ to the beginning of the facet name. If empty, all monitors in Datadog will be matched regardless of environment or service.
See [DataDog documentation](https://docs.datadoghq.com/monitors/manage/search/) for more information.`,
							Default: "",
						},
						"remote_priority_threshold": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "ALL",
							Description: `Monitor states with matching or higher priorities will be considered a failure in Sleuth. 
Options: ALL, P1, P2, P3, P4, P5. Defaults to ALL`,
						},
						"integration_slug": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "DataDog IntegrationAuthentication slug from app",
						},
					},
				},
			},
			"jira_input": {
				Description: "JIRA input",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"remote_jql": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "JIRA active incidents issues JQL",
						},
						"integration_slug": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "JIRA IntegrationAuthentication slug from app",
						},
					},
				},
			},
			"blameless_input": {
				Description: "Blameless input",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"remote_types": {
							Optional:    true,
							Type:        schema.TypeSet,
							Description: "The types of incidents to the monitors should track",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"remote_severity_threshold": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Incidents with matching or lower severities will be considered a failure in Sleuth",
						},
						"integration_slug": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Blameless IntegrationAuthentication slug from app",
						},
					},
				},
			},
			"slug": {
				Description: "Impact source slug",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceIncidentImpactSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	projectSlug, slug, err := getSlugsFromID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Project slug, incident slug %s, %s", projectSlug, slug))
	iis, err := c.GetIncidentImpactSource(ctx, projectSlug, slug)
	if err != nil {
		return diag.FromErr(err)
	}

	provider := ImpactProviderFromString(iis.Provider)

	setFields(ctx, d, iis, projectSlug, provider)

	return nil
}

func getProviderData(d *schema.ResourceData, i gqlclient.IncidentImpactSourceInputType, provider ImpactProvider) gqlclient.IncidentImpactSourceInputType {
	switch provider {
	case PagerDuty:
		i.PagerDutyInputType = &gqlclient.PagerDutyInputType{
			RemoteServices: d.Get("pagerduty_input.0.remote_services").(string),
			RemoteUrgency:  d.Get("pagerduty_input.0.remote_urgency").(string),
		}
	case DataDog:
		i.DataDogInputType = &gqlclient.DataDogInputType{
			DataDogProviderData: gqlclient.DataDogProviderData{
				Query:                   d.Get("datadog_input.0.query").(string),
				RemotePriorityThreshold: d.Get("datadog_input.0.remote_priority_threshold").(string),
			},
			IntegrationSlug: d.Get("datadog_input.0.integration_slug").(string),
		}
	case Jira:
		i.JiraInputType = &gqlclient.JiraInputType{
			JiraProviderData: gqlclient.JiraProviderData{
				RemoteJql: d.Get("jira_input.0.remote_jql").(string),
			},
		}
	case Blameless:
		i.BlamelessInputType = &gqlclient.BlamelessInputType{
			BlamelessProviderData: gqlclient.BlamelessProviderData{
				RemoteTypes:             expandStringList(d, "blameless_input.0.remote_types"),
				RemoteSeverityThreshold: d.Get("blameless_input.0.remote_severity_threshold").(string),
			},
		}
	}

	return i
}

func resourceIncidentImpactSourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	projectSlug := d.Get("project_slug").(string)
	providerStr := d.Get("provider_name").(string)

	provider := ImpactProviderFromString(providerStr)
	if provider == Unknown {
		return diag.FromErr(fmt.Errorf("unknown provider %s", providerStr))
	}

	input := gqlclient.IncidentImpactSourceInputType{
		ProjectSlug:     projectSlug,
		Name:            d.Get("name").(string),
		Provider:        provider.String(),
		EnvironmentName: strings.ToLower(d.Get("environment_name").(string)),
	}

	input = getProviderData(d, input, provider)

	tflog.Debug(ctx, fmt.Sprintf("CreateIncidentImpactSourceMutationInput: %v", input))

	incidentImpact, err := c.CreateIncidentImpactSource(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s/%s", projectSlug, incidentImpact.Slug))

	setFields(ctx, d, incidentImpact, projectSlug, provider)

	return nil
}

func resourceIncidentImpactSourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	projectSlug, slug, err := getSlugsFromID(d.Id())

	providerStr := d.Get("provider_name").(string)
	provider := ImpactProviderFromString(providerStr)
	if provider == Unknown {
		return diag.FromErr(fmt.Errorf("unknown provider %s", providerStr))
	}

	incident_input := gqlclient.IncidentImpactSourceInputType{
		ProjectSlug:     projectSlug,
		Name:            d.Get("name").(string),
		Provider:        d.Get("provider_name").(string),
		EnvironmentName: strings.ToLower(d.Get("environment_name").(string)),
	}

	incident_input = getProviderData(d, incident_input, provider)

	input := gqlclient.IncidentImpactSourceInputUpdateType{
		Slug:                          slug,
		IncidentImpactSourceInputType: incident_input,
	}

	proj, err := c.UpdateIncidentImpactSource(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("last_updated", time.Now().Format(time.RFC850))

	setFields(ctx, d, proj, projectSlug, provider)

	return nil
}

func resourceIncidentImpactSourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	projectSlug, slug, err := getSlugsFromID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	input := gqlclient.IncidentImpactSourceDeleteInputType{Slug: slug, ProjectSlug: projectSlug}

	succ, err := c.DeleteIncidentImpactSource(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}
	if !succ {
		return diag.FromErr(fmt.Errorf("unsuccessful deletion of incident impact source"))
	}

	return nil
}

func getSlugsFromID(id string) (string, string, error) {
	splits := strings.Split(id, "/")
	if len(splits) != 2 {
		return "", "", fmt.Errorf("invalid resource ID: %s", id)
	}

	return splits[0], splits[1], nil
}

func setProviderDetailsData(ctx context.Context, d *schema.ResourceData, is *gqlclient.IncidentImpactSource, provider ImpactProvider) {
	switch provider {
	case PagerDuty:
		pagerDutyInput := make(map[string]interface{})
		pagerDutyInput["remote_services"] = is.ProviderData.PagerDutyProviderData.RemoteServices
		pagerDutyInput["remote_urgency"] = is.ProviderData.PagerDutyProviderData.RemoteUrgency
		d.Set("pagerduty_input", []map[string]interface{}{pagerDutyInput})
	case DataDog:
		dataDogInput := make(map[string]interface{})
		dataDogInput["query"] = is.ProviderData.DataDogProviderData.Query
		dataDogInput["remote_priority_threshold"] = is.ProviderData.DataDogProviderData.RemotePriorityThreshold
		dataDogInput["integration_auth"] = is.IntegrationAuthSlug

		d.Set("datadog_input", []map[string]interface{}{dataDogInput})
	case Jira:
		jiraInput := make(map[string]interface{})
		jiraInput["remote_jql"] = is.ProviderData.JiraProviderData.RemoteJql
		jiraInput["integration_auth"] = is.IntegrationAuthSlug
	case Blameless:
		blamelessInput := make(map[string]interface{})
		blamelessInput["remote_types"] = flattenStringSet(is.ProviderData.BlamelessProviderData.RemoteTypes)
		blamelessInput["remote_severity_threshold"] = is.ProviderData.BlamelessProviderData.RemoteSeverityThreshold
		blamelessInput["integration_auth"] = is.IntegrationAuthSlug

		d.Set("blameless_input", []map[string]interface{}{blamelessInput})
	}
}

func setFields(ctx context.Context, d *schema.ResourceData, is *gqlclient.IncidentImpactSource, projectSlug string, provider ImpactProvider) {
	d.Set("name", is.Name)
	d.Set("slug", is.Slug)
	d.Set("provider_name", strings.ToUpper(is.Provider))
	d.Set("environment_name", is.Environment.Name)
	d.Set("project_slug", projectSlug)
	d.Set("register_impact_link", is.RegisterImpactLink)

	setProviderDetailsData(ctx, d, is, provider)
}

func flattenStringList(list []string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, v)
	}
	return vs
}

func flattenStringSet(list []string) *schema.Set {
	return schema.NewSet(schema.HashString, flattenStringList(list))
}

func expandStringList(d *schema.ResourceData, v string) []string {
	set := d.Get(v).(*schema.Set)
	list := set.List()
	stringList := make([]string, len(list))
	for i, v := range list {
		stringList[i] = v.(string)
	}

	return stringList
}
