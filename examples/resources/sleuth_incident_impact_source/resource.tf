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

