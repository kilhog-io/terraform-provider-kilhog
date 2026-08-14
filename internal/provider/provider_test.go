// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"kilhog": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("TF_ACC"); v == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' is set")
	}
}

func testAccProviderConfig() string {
	baseURL := os.Getenv("KILHOG_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	apiKey := os.Getenv("KILHOG_API_KEY")

	return fmtProviderConfig(baseURL, apiKey)
}

func fmtProviderConfig(baseURL, apiKey string) string {
	if apiKey == "" {
		return fmt.Sprintf(`
provider "kilhog" {
  base_url = %q
}
`, baseURL)
	}

	return fmt.Sprintf(`
provider "kilhog" {
  base_url = %q
  api_key  = %q
}
`, baseURL, apiKey)
}
