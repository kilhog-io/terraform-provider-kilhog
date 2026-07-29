// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccKilhogEphemeralResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		// Ephemeral resources are only available in 1.10 and later
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
		Steps: []resource.TestStep{
			{
				Config: testAccKilhogEphemeralResourceConfig("kilhog"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.test",
						tfjsonpath.New("data").AtMapKey("value"),
						knownvalue.StringExact("token-123"),
					),
				},
			},
		},
	})
}

func testAccKilhogEphemeralResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
ephemeral "kilhog_ephemeral" "test" {
  configurable_attribute = %[1]q
}

provider "echo" {
  data = ephemeral.kilhog_ephemeral.test
}

resource "echo" "test" {}
`, configurableAttribute)
}
