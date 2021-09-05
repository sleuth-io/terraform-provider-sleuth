
terraform {	
	required_providers {	
		sleuth = {		
			version = "0.1-dev"
			source	= "sleuth.io/core/sleuth"	
		}	
	}
}
provider "sleuth" {
    baseurl = "http://dev.sleuth.io"
    api_key = "2089d3139678bd250d33d4d8e5ef8c749e3ee588"
	org_slug = "myorg"
}

resource "sleuth_project" "myproject" {
	name = "My Project 2"
    impact_sensitivity = "FINE"
    change_failure_rate_boundary = "HEALTHY"
    description = "blah"
}

