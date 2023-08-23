resource "sleuth_metric_impact_source" "datadog_app_memory" {
  project_slug = "example_tf_app"
  environment_slug = "prod"
  name = "Application Memory"
  provider_type = "datadog"
  query = "avg:aws.ecs.memory_utilization{servicename:app-web}"
  less_is_better = true
}

resource "sleuth_metric_impact_source" "cloudwatch_rds_cpu" {
  project_slug = "example_tf_app"
  environment_slug = "prod"
  name = "RDS CPU"
  provider_type = "cloudwatch"
  integration_slug = "" /* If left empty or omitted completely, Sleuth will revert to `integration_slug == provider_type` */
  query = jsonencode({
    "metrics": [
      [ "AWS/RDS", "CPUUtilization", "DBInstanceIdentifier", "my-db-identifier", { "id": "m1" } ]
    ],
    "view": "timeSeries",
    "stacked": false,
    "region": "us-east-1",
    "stat": "Average",
    "period": 300
  })
  less_is_better = true
}
