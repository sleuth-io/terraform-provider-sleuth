resource "sleuth_code_change_source" "sleuth-terraform-provider" {
  project_slug = "example_tf_app"
  name         = "terraform-provider-sleuth"
  repository = {
    name     = "terraform-provider-sleuth"
    owner    = "sleuth-io"
    provider = "GITHUB"
    url      = "https://github.com/sleuth-io/terraform-provider-sleuth"
  }
  environment_mappings = [
      {
        environment_slug = "prod"
        branch           = "main"
      },
      {
        environment_slug = "stage"
        branch           = "dev"
      }
  ]
  deploy_tracking_type = "manual"
  collect_impact       = true
  path_prefix = jsonencode({
    excludes = [""]
    includes = [""]
  })
}
