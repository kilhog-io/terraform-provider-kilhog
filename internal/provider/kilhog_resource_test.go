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
)

func TestAccKilhogResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccKilhogResourceConfig("one"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"kilhog_resource.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact("kilhog-id"),
					),
					statecheck.ExpectKnownValue(
						"kilhog_resource.test",
						tfjsonpath.New("defaulted"),
						knownvalue.StringExact("kilhog default value"),
					),
					statecheck.ExpectKnownValue(
						"kilhog_resource.test",
						tfjsonpath.New("configurable_attribute"),
						knownvalue.StringExact("one"),
					),
				},
			},
			// ImportState testing
			{
				ResourceName:      "kilhog_resource.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"configurable_attribute", "defaulted"},
			},
			// Update and Read testing
			{
				Config: testAccKilhogResourceConfig("two"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"kilhog_resource.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact("kilhog-id"),
					),
					statecheck.ExpectKnownValue(
						"kilhog_resource.test",
						tfjsonpath.New("defaulted"),
						knownvalue.StringExact("kilhog default value"),
					),
					statecheck.ExpectKnownValue(
						"kilhog_resource.test",
						tfjsonpath.New("configurable_attribute"),
						knownvalue.StringExact("two"),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccKilhogResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "kilhog_resource" "test" {
  configurable_attribute = %[1]q
}
`, configurableAttribute)
}
