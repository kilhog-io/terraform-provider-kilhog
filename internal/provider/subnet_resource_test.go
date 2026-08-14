// Copyright IBM Corp. 2021, 2026
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

func TestAccSubnetResource_basic(t *testing.T) {
	networkName := acctest.RandomWithPrefix("tf-network")
	subnetName := acctest.RandomWithPrefix("tf-subnet")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetConfig(networkName, subnetName, "10.10.0.0", 24, "created by terraform"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"kilhog_subnet.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(subnetName),
					),
					statecheck.ExpectKnownValue(
						"kilhog_subnet.test",
						tfjsonpath.New("address"),
						knownvalue.StringExact("10.10.0.0"),
					),
					statecheck.ExpectKnownValue(
						"kilhog_subnet.test",
						tfjsonpath.New("prefix"),
						knownvalue.Int64Exact(24),
					),
					statecheck.ExpectKnownValue(
						"kilhog_subnet.test",
						tfjsonpath.New("parent_kind"),
						knownvalue.StringExact("network"),
					),
				},
			},
			{
				Config: testAccSubnetConfig(networkName, subnetName, "10.10.0.0", 24, "updated by terraform"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"kilhog_subnet.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("updated by terraform"),
					),
				},
			},
		},
	})
}

func testAccSubnetConfig(networkName, subnetName, address string, prefix int, description string) string {
	return fmt.Sprintf(`
%s

resource "kilhog_network" "test" {
  name = %q
}

resource "kilhog_subnet" "test" {
  network_id  = kilhog_network.test.id
  name        = %q
  description = %q
  address     = %q
  prefix      = %d
}
`, testAccProviderConfig(), networkName, subnetName, description, address, prefix)
}
