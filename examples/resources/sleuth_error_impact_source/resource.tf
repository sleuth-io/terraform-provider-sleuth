resource "sleuth_error_impact_source" "sentry_production" {
  project_slug = "example_tf_app"
  environment_slug = "prod"
  name = "Sentry errors"
  provider_type = "sentry"
  error_org_key = "my-sentry-org-key"
  error_project_key = "my-sentry-project-key"
  error_environment = "my-sentry-environment"
}
