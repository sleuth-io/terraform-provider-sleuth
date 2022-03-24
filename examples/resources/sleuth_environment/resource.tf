resource "sleuth_environment" "prod" {
  project_slug = "example_tf_app"
  name = "Production"
  color = "#279694"
}

resource "sleuth_environment" "stage" {
  project_slug = "example_tf_app"
  name = "Staging"
  color = "#58b94b"
}
