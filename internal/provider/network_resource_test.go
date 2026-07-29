// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccNetworkResource_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-network")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConfig(name, "created by terraform"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"kilhog_network.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(name),
					),
					statecheck.ExpectKnownValue(
						"kilhog_network.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("created by terraform"),
					),
				},
			},
			{
				Config: testAccNetworkConfig(name, "updated by terraform"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"kilhog_network.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("updated by terraform"),
					),
				},
			},
		},
	})
}

func testAccNetworkConfig(name, description string) string {
	return fmt.Sprintf(`
%s

resource "kilhog_network" "test" {
  name        = %q
  description = %q
}
`, testAccProviderConfig(), name, description)
}
