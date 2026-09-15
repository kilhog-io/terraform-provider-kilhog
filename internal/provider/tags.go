// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	kilhogsdk "github.com/kilhog-io/kilhog/pkg/kilhog"
)

// expandTags converts the Terraform map(string) representation of tags into the
// list of key/value pairs the Kilhog SDK expects. Keys are sorted so the request
// payload is deterministic regardless of map iteration order.
func expandTags(ctx context.Context, tags types.Map) ([]kilhogsdk.Tag, diag.Diagnostics) {
	var diags diag.Diagnostics

	if tags.IsNull() || tags.IsUnknown() {
		return nil, diags
	}

	elements := make(map[string]string, len(tags.Elements()))
	diags.Append(tags.ElementsAs(ctx, &elements, false)...)
	if diags.HasError() {
		return nil, diags
	}

	keys := make([]string, 0, len(elements))
	for key := range elements {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	result := make([]kilhogsdk.Tag, 0, len(elements))
	for _, key := range keys {
		result = append(result, kilhogsdk.Tag{
			Key:   key,
			Value: elements[key],
		})
	}

	return result, diags
}

// flattenTags converts the Kilhog SDK list of tags into a Terraform map(string).
// Using a map makes the attribute order-insensitive, so the API returning tags
// in a different order than they were written no longer produces a diff.
func flattenTags(ctx context.Context, tags []kilhogsdk.Tag) (types.Map, diag.Diagnostics) {
	if len(tags) == 0 {
		return types.MapNull(types.StringType), nil
	}

	elements := make(map[string]string, len(tags))
	for _, tag := range tags {
		elements[tag.Key] = tag.Value
	}

	return types.MapValueFrom(ctx, types.StringType, elements)
}
