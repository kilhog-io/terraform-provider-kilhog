// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"sort"
	"strings"
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
				// Keys are intentionally not in sorted order: the API returns
				// tags sorted on read, so an order-sensitive attribute would
				// produce a permanent diff and fail the post-apply plan check.
				Config: testAccNetworkConfig(name, "created by terraform", map[string]string{
					"zeta":  "1",
					"alpha": "2",
					"mike":  "3",
				}),
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
					statecheck.ExpectKnownValue(
						"kilhog_network.test",
						tfjsonpath.New("tags"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"zeta":  knownvalue.StringExact("1"),
							"alpha": knownvalue.StringExact("2"),
							"mike":  knownvalue.StringExact("3"),
						}),
					),
				},
			},
			{
				ResourceName:      "kilhog_network.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Change tags (add, remove, and update values) to verify updates.
				Config: testAccNetworkConfig(name, "updated by terraform", map[string]string{
					"alpha": "20",
					"new":   "x",
				}),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"kilhog_network.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("updated by terraform"),
					),
					statecheck.ExpectKnownValue(
						"kilhog_network.test",
						tfjsonpath.New("tags"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"alpha": knownvalue.StringExact("20"),
							"new":   knownvalue.StringExact("x"),
						}),
					),
				},
			},
		},
	})
}

func testAccNetworkConfig(name, description string, tags map[string]string) string {
	return fmt.Sprintf(`
%s

resource "kilhog_network" "test" {
  name        = %q
  description = %q

  tags = %s
}
`, testAccProviderConfig(), name, description, renderTagsHCL(tags))
}

// renderTagsHCL renders a Go map as an HCL map literal with keys in sorted order
// so the generated config is stable across runs.
func renderTagsHCL(tags map[string]string) string {
	if len(tags) == 0 {
		return "{}"
	}

	keys := make([]string, 0, len(tags))
	for key := range tags {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var b strings.Builder
	b.WriteString("{\n")
	for _, key := range keys {
		fmt.Fprintf(&b, "    %q = %q\n", key, tags[key])
	}
	b.WriteString("  }")

	return b.String()
}
