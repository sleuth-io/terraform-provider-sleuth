package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/sleuth-io/terraform-provider-sleuth/internal/gqlclient"
)

func resourceCodeChangeSource() *schema.Resource {
	return &schema.Resource{
		Description: "Sleuth code change source.",

		CreateContext: resourceCodeChangeSourceCreate,
		ReadContext:   resourceCodeChangeSourceRead,
		UpdateContext: resourceCodeChangeSourceUpdate,
		DeleteContext: resourceCodeChangeSourceDelete,

		Schema: map[string]*schema.Schema{
			"project_slug": {
				Description: "The project for this environment",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "Change source name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"repository": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"owner": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The repository owner, usually the organization or user name",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The repository name",
						},
						"url": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The repository url, used for links",
						},
						"provider": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The repository provider, such as GITHUB",
						},
						"integration_slug": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IntegrationAuthentication slug used",
						},
						"repo_uid": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Repository UID, required only for AZURE provider. You can obtain data from [API](https://learn.microsoft.com/en-us/rest/api/azure/devops/git/repositories/list?view=azure-devops-rest-6.0&tabs=HTTP)",
						},
						"project_uid": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Project UID, required only for AZURE provider. You can obtain data from [API](https://learn.microsoft.com/en-us/rest/api/azure/devops/git/repositories/list?view=azure-devops-rest-6.0&tabs=HTTP)",
						},
					},
				},
			},
			"environment_mappings": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Environment mappings of the environment. They must be ordered by environment_slug ascending to avoid Terraform plan changes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"environment_slug": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The environment slug or id",
						},
						"branch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The repository branch name for the environment",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Computed ID",
						},
					},
				},
			},
			"build_mappings": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"environment_slug": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The environment slug or id",
						},
						"provider": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The repository provider, such as CIRCLECI",
						},
						"integration_slug": {
							Description: "The integration slug",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"build_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The remote build or pipeline name",
						},
						"job_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The job or stage within the build or pipeline, if supported",
						},
						"project_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The build project key",
						},
						"match_branch_to_environment": {
							Type:     schema.TypeBool,
							Optional: true,
							Description: "Whether only builds performed on the branch mapped from the environment are " +
								"tracked or not. Basically if you only want Sleuth to find builds that were triggered" +
								"by a change on the branch that is configured for the environment, set this to false. " +
								"Defaults to true",
							Default: true,
						},
					},
				},
			},
			"deploy_tracking_type": {
				Description: "How to track deploys. Valid choices are build, manual, auto_pr, auto_tag, auto_push",
				Type:        schema.TypeString,
				Required:    true,
			},
			"collect_impact": {
				Description: "Whether to collect impact for its deploys",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"path_prefix": {
				Description: "What code source path to limit this deployment to. Useful for monorepos. Must be used with the [jsonencode()](https://developer.hashicorp.com/terraform/language/functions/jsonencode) function to specify the paths to include and/or exclude respectively. (see example above)",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"notify_in_slack": {
				Description: "Whether to send Slack notifications for deploys or not",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"include_in_dashboard": {
				Description: "Whether to include deploys from this change source in the metrics dashboard",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"auto_tracking_delay": {
				Description: "The delay to add to a deployment event",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func validateCodeChangeInput(ccs gqlclient.MutableCodeChangeSource) error {
	if strings.ToLower(ccs.Repository.Provider) != "azure" {
		return nil
	}

	if ccs.Repository.ProjectUID == "" || ccs.Repository.RepoUID == "" || ccs.Repository.IntegrationSlug == "" {
		return fmt.Errorf("project_uid, repo_uid and integration_slug are required for AZURE provider")
	}
	return nil
}

func resourceCodeChangeSourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	projectSlug := d.Get("project_slug").(string)
	inputFields := gqlclient.MutableCodeChangeSource{}
	input := gqlclient.CreateCodeChangeSourceMutationInput{ProjectSlug: projectSlug, MutableCodeChangeSource: &inputFields}
	input.InitializeChanges = true
	diags = populateCodeChangeSource(d, &inputFields, diags)

	if err := validateCodeChangeInput(inputFields); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error validating repository input",
			Detail:   err.Error(),
		})

		return diags
	}

	src, err := c.CreateCodeChangeSource(input)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s/%s", projectSlug, src.Slug))
	setCodeChangeSourceFields(ctx, d, projectSlug, src)

	return diags
}

func resourceCodeChangeSourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	parsed := strings.Split(d.Id(), "/")
	projectSlug := parsed[0]
	changeSourceSlug := parsed[1]

	inputFields := gqlclient.MutableCodeChangeSource{}
	input := gqlclient.UpdateCodeChangeSourceMutationInput{ProjectSlug: projectSlug, Slug: changeSourceSlug, MutableCodeChangeSource: &inputFields}
	diags = populateCodeChangeSource(d, &inputFields, diags)

	source, err := c.UpdateCodeChangeSource(input)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("last_updated", time.Now().Format(time.RFC850))
	setCodeChangeSourceFields(ctx, d, projectSlug, source)

	return diags
}

func resourceCodeChangeSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	parsed := strings.Split(d.Id(), "/")
	projectSlug := parsed[0]
	changeSlug := parsed[1]

	source, err := c.GetCodeChangeSource(ctx, &projectSlug, &changeSlug)
	if err != nil {
		return diag.FromErr(err)
	} else if source == nil {
		d.SetId("")
	} else {
		setCodeChangeSourceFields(ctx, d, projectSlug, source)
	}
	return diags

}

func setCodeChangeSourceFields(ctx context.Context, d *schema.ResourceData, projectSlug string, source *gqlclient.CodeChangeSource) {
	repository := make(map[string]interface{})
	repository["owner"] = source.Repository.Owner
	repository["name"] = source.Repository.Name
	repository["provider"] = strings.ToUpper(source.Repository.Provider)
	repository["url"] = source.Repository.Url
	repository["repo_uid"] = source.Repository.RepoUID
	repository["project_uid"] = source.Repository.ProjectUID

	if source.Repository.IntegrationAuth != nil {
		repository["integration_slug"] = source.Repository.IntegrationAuth.Slug
	}
	var repositoryList [1]map[string]interface{}
	repositoryList[0] = repository

	environmentMappings := make([]map[string]interface{}, len(source.EnvironmentMappings))
	for idx, v := range source.EnvironmentMappings {
		m := make(map[string]interface{})
		m["branch"] = v.Branch
		m["environment_slug"] = v.EnvironmentSlug
		m["id"] = fmt.Sprintf("%s/%s", projectSlug, v.EnvironmentSlug)
		environmentMappings[idx] = m
	}

	buildMappings := make([]map[string]interface{}, len(source.DeployTrackingBuildMappings))
	for idx, v := range source.DeployTrackingBuildMappings {
		m := make(map[string]interface{})
		m["build_name"] = v.BuildName
		m["job_name"] = v.JobName
		m["provider"] = v.Provider
		m["project_key"] = v.BuildProjectKey
		m["match_branch_to_environment"] = v.MatchBranchToEnvironment
		m["environment_slug"] = v.Environment.Slug

		buildMappings[idx] = m
	}

	d.Set("project_slug", projectSlug)
	d.Set("name", source.Name)
	d.Set("repository", repositoryList)
	d.Set("environment_mappings", environmentMappings)
	d.Set("build_mappings", buildMappings)
	d.Set("auto_tracking_delay", source.AutoTrackingDelay)
	d.Set("include_in_dashboard", source.IncludeInDashboard)
	d.Set("path_prefix", source.PathPrefix)
	d.Set("notify_in_slack", source.NotifyInSlack)
	d.Set("collect_impact", source.CollectImpact)
	d.Set("deploy_tracking_type", source.DeployTrackingType)
}

func populateCodeChangeSource(d *schema.ResourceData, input *gqlclient.MutableCodeChangeSource, diags diag.Diagnostics) diag.Diagnostics {
	repoList := d.Get("repository").([]interface{})
	repoData := repoList[0].(map[string]interface{})
	repo := gqlclient.MutableRepository{
		IntegrationSlug: repoData["integration_slug"].(string),
		Repository: gqlclient.Repository{
			Owner:      repoData["owner"].(string),
			Name:       repoData["name"].(string),
			Provider:   repoData["provider"].(string),
			Url:        repoData["url"].(string),
			RepoUID:    repoData["repo_uid"].(string),
			ProjectUID: repoData["project_uid"].(string),
		},
	}

	environmentMappingsData := d.Get("environment_mappings").([]interface{})
	environmentMappings := make([]gqlclient.BranchMapping, len(environmentMappingsData))
	for idx, v := range environmentMappingsData {
		m := v.(map[string]interface{})
		var envRaw = m["environment_slug"].(string)
		mapping := gqlclient.BranchMapping{EnvironmentSlug: envRaw, Branch: m["branch"].(string)}
		environmentMappings[idx] = mapping
	}

	buildMappingsData := d.Get("build_mappings").([]interface{})
	buildMappings := make([]gqlclient.BuildMapping, len(buildMappingsData))
	for idx, v := range buildMappingsData {
		m := v.(map[string]interface{})
		var envRaw = m["environment_slug"].(string)

		for _, v2 := range environmentMappings {
			var envSlug = v2.EnvironmentSlug
			if strings.Contains(v2.EnvironmentSlug, "/") {
				envSlug = strings.Split(v2.EnvironmentSlug, "/")[1]
			}
			if envRaw == envSlug {
				mapping := gqlclient.BuildMapping{EnvironmentSlug: envSlug,
					BuildName:                m["build_name"].(string),
					JobName:                  m["job_name"].(string),
					Provider:                 m["provider"].(string),
					BuildProjectKey:          m["project_key"].(string),
					MatchBranchToEnvironment: m["match_branch_to_environment"].(bool),
					IntegrationSlug:          m["integration_slug"].(string),
					BuildBranch:              v2.Branch,
				}
				buildMappings[idx] = mapping
				break
			}
		}
	}

	input.Name = d.Get("name").(string)
	input.Repository = repo
	input.EnvironmentMappings = environmentMappings
	input.BuildMappings = buildMappings
	input.AutoTrackingDelay = d.Get("auto_tracking_delay").(int)
	input.IncludeInDashboard = d.Get("include_in_dashboard").(bool)
	input.PathPrefix = d.Get("path_prefix").(string)
	input.NotifyInSlack = d.Get("notify_in_slack").(bool)
	input.CollectImpact = d.Get("collect_impact").(bool)
	input.DeployTrackingType = d.Get("deploy_tracking_type").(string)

	return diags
}

func resourceCodeChangeSourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	parsed := strings.Split(d.Id(), "/")
	projectSlug := parsed[0]
	sourceSlug := parsed[1]

	err := c.DeleteChangeSource(&projectSlug, &sourceSlug)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
