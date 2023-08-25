package sleuth

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"sleuth": providerserver.NewProtocol6WithError(New()),
	}
	testAccProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
		"sleuth": providerserver.NewProtocol5WithError(New()),
	}
)
