// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	kilhogsdk "github.com/kilhog-io/kilhog/pkg/kilhog"
)

type tagModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func expandTags(ctx context.Context, tags types.List) ([]kilhogsdk.Tag, diag.Diagnostics) {
	if tags.IsNull() || tags.IsUnknown() {
		return nil, nil
	}

	var models []tagModel
	var diags diag.Diagnostics

	diags.Append(tags.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]kilhogsdk.Tag, 0, len(models))
	for _, model := range models {
		result = append(result, kilhogsdk.Tag{
			Key:   model.Key.ValueString(),
			Value: model.Value.ValueString(),
		})
	}

	return result, diags
}

func flattenTags(ctx context.Context, tags []kilhogsdk.Tag) (types.List, diag.Diagnostics) {
	if len(tags) == 0 {
		return types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"key":   types.StringType,
				"value": types.StringType,
			},
		}), nil
	}

	models := make([]tagModel, 0, len(tags))
	for _, tag := range tags {
		models = append(models, tagModel{
			Key:   types.StringValue(tag.Key),
			Value: types.StringValue(tag.Value),
		})
	}

	return types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"key":   types.StringType,
			"value": types.StringType,
		},
	}, models)
}
