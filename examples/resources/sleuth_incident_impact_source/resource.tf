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
