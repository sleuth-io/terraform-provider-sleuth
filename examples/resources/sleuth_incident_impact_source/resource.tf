resource "sleuth_incident_impact_source" "pd" {
    project_slug = "project_slug"
    name = "PagerDuty TF incident impact"
    environment_name = "environment_name"
    provider_name = "PAGERDUTY"
    pagerduty_input {
        remote_services = ""
        remote_urgency = "ANY"
    }
}

resource "sleuth_incident_impact_source" "dd" {
    project_slug = "project_slug"
    name = "DataDog TF incident impact"
    environment_name = "environment_name"
    provider_name = "DATADOG"
    datadog_input {
        query = "@query=123" # use @ if you are using facets in DataDog
        remote_priority_threshold = "ALL" # or P1 to P5
        integration_slug = "optional_integration_slug"
    }
}

resource "sleuth_incident_impact_source" "jira" {
    project_slug = "project_slug"
    name = "JIRA TF incident impact"
    environment_name = "environment_name"
    provider_name = "JIRA"
    jira_input {
        remote_jql = "status IN (\"Incident\")"
        integration_slug = "optional_integration_slug"
    }
}

resource "sleuth_incident_impact_source" "blameless" {
    project_slug = "project_slug"
    name = "Blameless TF incident impact"
    environment_name = "environment_name"
    provider_name = "BLAMELESS"
    blameless_input {
        remote_types = ["type1", "type2"]
	remote_severity_threshold = "SEV1"
        integration_slug = "optional_integration_slug"
    }
}

