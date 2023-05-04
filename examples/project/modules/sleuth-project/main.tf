resource "sleuth_project" "projects" {
  name                         = var.project.name
  build_provider               = try(var.project.build_provider, null)
  change_failure_rate_boundary = try(var.project.change_failure_rate_boundary, null)
  failure_sensitivity          = try(var.project.failure_sensitivity, null)
  impact_sensitivity           = try(var.project.impact_sensitivity, null)
  issue_tracker_provider_type  = try(var.project.issue_tracker_provider_type, null)
}

resource "sleuth_environment" "envs" {

  for_each     = { for env in var.environments : env.name => env }
  project_slug = sleuth_project.projects.slug
  name         = each.key
  description  = try(each.value.description, null)
  color        = try(each.value.color, null)

  depends_on = [
    sleuth_project.projects
  ]

}

resource "sleuth_code_change_source" "repo_source" {
  for_each = { for repo in var.repositories : repo.name => repo }

  project_slug = sleuth_project.projects.slug
  name         = each.key
  repository {
    name     = each.value.name
    owner    = each.value.owner
    provider = each.value.provider
    url      = each.value.url
  }

  dynamic "environment_mappings" {
    for_each = each.value.environment_mappings
    content {
      environment_slug = sleuth_environment.envs["${environment_mappings.value.environment_name}"].slug
      branch           = environment_mappings.value.branch
    }
  }

  deploy_tracking_type = try(each.value.deploy_tracking_type, "manual")
  collect_impact       = try(each.value.collect_impact, true)

  depends_on = [
    sleuth_project.projects,
    sleuth_environment.envs,
  ]
}
