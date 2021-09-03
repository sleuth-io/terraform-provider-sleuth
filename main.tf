
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
    api_key = "9208af1845ec3517ee80f8113097f57d9da1a6b2"
	org_slug = "myorg"
}

resource "sleuth_project" "myproject" {
	name = "My Project is cool right mate friend"
    failure_sensitivity = 423
}

