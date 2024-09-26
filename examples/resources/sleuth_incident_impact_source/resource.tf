resource "sleuth_incident_impact_source" "pd" {
  project_slug     = "project_slug"
  name             = "PagerDuty TF incident impact"
  environment_name = "environment_name"
  provider_name    = "pagerduty"
  pagerduty_input = {
    remote_services = ""
    remote_urgency  = "ANY"
  }
}

resource "sleuth_incident_impact_source" "dd" {
  project_slug     = "project_slug"
  name             = "DataDog TF incident impact"
  environment_name = "environment_name"
  provider_name    = "datadog"
  datadog_input = {
    query                     = "@query=123" # use @ if you are using facets in DataDog
    remote_priority_threshold = "ALL"        # or P1 to P5
    integration_slug          = "optional_integration_slug"
  }
}

resource "sleuth_incident_impact_source" "jira" {
  project_slug     = "project_slug"
  name             = "JIRA TF incident impact"
  environment_name = "environment_name"
  provider_name    = "jira"
  jira_input = {
    remote_jql       = "status IN (\"Incident\")"
    integration_slug = "optional_integration_slug"
  }
}

resource "sleuth_incident_impact_source" "blameless" {
  project_slug     = "project_slug"
  name             = "Blameless TF incident impact"
  environment_name = "environment_name"
  provider_name    = "blameless"
  blameless_input = {
    remote_types              = ["type1", "type2"]
    remote_severity_threshold = "SEV1"
    integration_slug          = "optional_integration_slug"
  }
}

resource "sleuth_incident_impact_source" "statuspage" {
  project_slug     = "project_slug"
  name             = "Statuspage TF incident impact"
  environment_name = "environment_name"
  provider_name    = "statuspage"
  statuspage_input = {
    remote_page                  = "remote_page"
    remote_component             = "remote_component"
    remote_impact                = "remote_impact"
    ignore_maintenance_incidents = false
  }
}

resource "sleuth_incident_impact_source" "opsgenie" {
  project_slug     = "project_slug"
  name             = "OpsGenie TF incident impact"
  environment_name = "environment_name"
  provider_name    = "opsgeanie"
  opsgenie_input = {
    remote_alert_tags         = "tag1"
    remote_incidents_tags     = "tag1"
    remote_priority_threshold = "P1"
    remote_service            = "test_service"
    remote_use_alerts         = false
  }
}

resource "sleuth_incident_impact_source" "firehydrant" {
  project_slug     = "project_slug"
  name             = "FireHydrant TF incident impact"
  environment_name = "environment_name"
  provider_name    = "firehydrant"
  firehydrant_input = {
    remote_services             = "service_uuid"
    remote_environments         = "environment_uuid"
    remote_mitigated_is_healthy = true
  }
}

resource "sleuth_incident_impact_source" "clubhouse" {
  project_slug     = "project_slug"
  name             = "Clubhouse TF incident impact"
  environment_name = "environment_name"
  provider_name    = "clubhouse"
  clubhouse_input = {
    remote_query     = "id:135"
    integration_slug = "optional_integration_slug"
  }
}

resource "sleuth_incident_impact_source" "rootly" {
  project_slug     = "project_slug"
  name             = "Rootly TF incident impact"
  environment_name = "environment_name"
  provider_name    = "rootly"
  rootly_input = {
    remote_severity      = "ALL" # or "CRITICAL", "HIGH", "MEDIUM", "LOW"
    remote_incident_type = "remote_incident_type_id"
    remote_environment   = "remote_environment_id"
    remote_service       = "remote_service_id"
    remote_team          = "remote_team_id"
    integration_slug     = "optional_integration_slug"
  }
}
