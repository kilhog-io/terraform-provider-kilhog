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

func TestAccNetworkDataSource_byName(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-network-ds")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkDataSourceConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.kilhog_network.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(name),
					),
					statecheck.ExpectKnownValue(
						"data.kilhog_network.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("read by data source"),
					),
				},
			},
		},
	})
}

func testAccNetworkDataSourceConfig(name string) string {
	return fmt.Sprintf(`
%s

resource "kilhog_network" "test" {
  name        = %q
  description = "read by data source"
}

data "kilhog_network" "test" {
  name = kilhog_network.test.name
}
`, testAccProviderConfig(), name)
}
