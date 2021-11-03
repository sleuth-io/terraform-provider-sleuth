---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sleuth_metric_impact_source Resource - terraform-provider-sleuth"
subcategory: ""
description: |-
  Sleuth error impact source.
---

# sleuth_metric_impact_source (Resource)

Sleuth error impact source.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **environment_slug** (String) The color for the UI
- **name** (String) Environment name
- **project_slug** (String) The project for this environment
- **provider_type** (String) Integration provider type
- **query** (String) The metric query

### Optional

- **id** (String) The ID of this resource.
- **less_is_better** (Boolean) Whether smaller values are better or not
- **manually_set_health_threshold** (Number) The environment of the integration provider

