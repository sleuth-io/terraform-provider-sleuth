package provider

import (
	"errors"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"sleuth": func() (*schema.Provider, error) {
		return New("dev")(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	t.Skip("Skip")
	if err := os.Getenv("SLEUTH_BASEURL"); err == "" {
		t.Fatal("SLEUTH_BASEURL must be set for acceptance tests")
	}
	if err := os.Getenv("SLEUTH_API_KEY"); err == "" {
		t.Fatal("SLEUTH_API_KEY must be set for acceptance tests")
	}
}

func testAccCheckOrganization() error {
	baseUrl := os.Getenv("SLEUTH_BASEURL")
	apiKey := os.Getenv("SLEUTH_API_KEY")

	if baseUrl == "" || apiKey == "" {
		return errors.New("SLEUTH_BASEURL and SLEUTH_API_KEY must be set for acceptance tests")
	}
	return nil
}
