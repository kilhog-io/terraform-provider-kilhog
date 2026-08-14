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

func TestAccSubnetDataSource_byName(t *testing.T) {
	networkName := acctest.RandomWithPrefix("tf-network-ds")
	subnetName := acctest.RandomWithPrefix("tf-subnet-ds")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetDataSourceConfig(networkName, subnetName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.kilhog_subnet.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(subnetName),
					),
					statecheck.ExpectKnownValue(
						"data.kilhog_subnet.test",
						tfjsonpath.New("address"),
						knownvalue.StringExact("10.20.0.0"),
					),
					statecheck.ExpectKnownValue(
						"data.kilhog_subnet.test",
						tfjsonpath.New("prefix"),
						knownvalue.Int64Exact(24),
					),
				},
			},
		},
	})
}

func testAccSubnetDataSourceConfig(networkName, subnetName string) string {
	return fmt.Sprintf(`
%s

resource "kilhog_network" "test" {
  name = %q
}

resource "kilhog_subnet" "test" {
  network_id = kilhog_network.test.id
  name       = %q
  address    = "10.20.0.0"
  prefix     = 24
}

data "kilhog_subnet" "test" {
  network_id = kilhog_network.test.id
  name       = kilhog_subnet.test.name
}
`, testAccProviderConfig(), networkName, subnetName)
}
