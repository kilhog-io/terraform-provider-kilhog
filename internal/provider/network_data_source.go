// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	kilhogsdk "github.com/kilhog-io/kilhog/pkg/kilhog"
)

var _ datasource.DataSource = &NetworkDataSource{}

func NewNetworkDataSource() datasource.DataSource {
	return &NetworkDataSource{}
}

type NetworkDataSource struct {
	client *kilhogsdk.Client
}

type NetworkDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Tags        types.List   `tfsdk:"tags"`
}

func (d *NetworkDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network"
}

func (d *NetworkDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads an existing Kilhog network by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Network UUID. Exactly one of `id` or `name` must be set.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Unique network name. Exactly one of `id` or `name` must be set.",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Network description.",
				Computed:            true,
			},
			"tags": schema.ListNestedAttribute{
				MarkdownDescription: "Key-value metadata tags.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							MarkdownDescription: "Tag key.",
							Computed:            true,
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "Tag value.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *NetworkDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*kilhogsdk.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *kilhog.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *NetworkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data NetworkDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	network, diags := d.findNetwork(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceModel := NetworkResourceModel{}
	resp.Diagnostics.Append(flattenNetwork(ctx, network, &resourceModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = resourceModel.ID
	data.Name = resourceModel.Name
	data.Description = resourceModel.Description
	data.Tags = resourceModel.Tags

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *NetworkDataSource) findNetwork(ctx context.Context, data NetworkDataSourceModel) (*kilhogsdk.Network, diag.Diagnostics) {
	var diags diag.Diagnostics

	hasID := !data.ID.IsNull() && !data.ID.IsUnknown()
	hasName := !data.Name.IsNull() && !data.Name.IsUnknown()

	switch {
	case hasID && hasName:
		diags.AddError(
			"Invalid configuration",
			"Exactly one of `id` or `name` must be set.",
		)
		return nil, diags
	case !hasID && !hasName:
		diags.AddError(
			"Invalid configuration",
			"Exactly one of `id` or `name` must be set.",
		)
		return nil, diags
	case hasID:
		networkID, err := uuid.Parse(data.ID.ValueString())
		if err != nil {
			diags.AddError("Invalid network ID", err.Error())
			return nil, diags
		}

		network, err := d.client.GetNetwork(ctx, networkID)
		if err != nil {
			diags.AddError("Error reading network", apiErrorSummary(err, "Read network"))
			return nil, diags
		}

		return network, diags
	default:
		networks, err := d.client.ListNetworks(ctx)
		if err != nil {
			diags.AddError("Error listing networks", apiErrorSummary(err, "List networks"))
			return nil, diags
		}

		name := data.Name.ValueString()
		for i := range networks {
			if networks[i].Name == name {
				return &networks[i], diags
			}
		}

		diags.AddError(
			"Network not found",
			fmt.Sprintf("No network named %q was found.", name),
		)
		return nil, diags
	}
}
