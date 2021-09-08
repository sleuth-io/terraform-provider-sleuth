
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
}

resource "sleuth_project" "myproject" {
	name = "My project good"
}

resource "sleuth_environment" "myenvironment" {
	project_slug = "${sleuth_project.myproject.id}"
	name = "Prod"
	description = "blah"
	color = "#aa33ff"
}

resource "sleuth_environment" "myenvironmentstg" {
	project_slug = "${sleuth_project.myproject.id}"
	name = "Staging"
	description = "blah 2"
	color = "#3333ff"
}
