// Copyright IBM Corp. 2021, 2026
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

var _ datasource.DataSource = &SubnetDataSource{}

func NewSubnetDataSource() datasource.DataSource {
	return &SubnetDataSource{}
}

type SubnetDataSource struct {
	client *kilhogsdk.Client
}

type SubnetDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	NetworkID      types.String `tfsdk:"network_id"`
	ParentSubnetID types.String `tfsdk:"parent_subnet_id"`
	ParentKind     types.String `tfsdk:"parent_kind"`
	ParentID       types.String `tfsdk:"parent_id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Prefix         types.Int64  `tfsdk:"prefix"`
	Address        types.String `tfsdk:"address"`
	Type           types.String `tfsdk:"type"`
}

func (d *SubnetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subnet"
}

func (d *SubnetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads an existing Kilhog subnet by network ID and subnet ID or name.",
		Attributes: map[string]schema.Attribute{
			"network_id": schema.StringAttribute{
				MarkdownDescription: "UUID of the network that owns the subnet.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Subnet UUID. Exactly one of `id` or `name` must be set.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Subnet name within the network. Exactly one of `id` or `name` must be set.",
				Optional:            true,
				Computed:            true,
			},
			"parent_subnet_id": schema.StringAttribute{
				MarkdownDescription: "UUID of the parent subnet when the subnet is nested.",
				Computed:            true,
			},
			"parent_kind": schema.StringAttribute{
				MarkdownDescription: "Parent kind reported by the API (`network` or `subnet`).",
				Computed:            true,
			},
			"parent_id": schema.StringAttribute{
				MarkdownDescription: "UUID of the parent resource reported by the API.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Subnet description.",
				Computed:            true,
			},
			"prefix": schema.Int64Attribute{
				MarkdownDescription: "IPv4 prefix length.",
				Computed:            true,
			},
			"address": schema.StringAttribute{
				MarkdownDescription: "Network address for the subnet.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Address family.",
				Computed:            true,
			},
		},
	}
}

func (d *SubnetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SubnetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SubnetDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subnet, diags := d.findSubnet(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceModel := SubnetResourceModel{
		NetworkID: data.NetworkID,
	}
	flattenSubnet(subnet, &resourceModel)

	data.ID = resourceModel.ID
	data.Name = resourceModel.Name
	data.ParentSubnetID = resourceModel.ParentSubnetID
	data.ParentKind = resourceModel.ParentKind
	data.ParentID = resourceModel.ParentID
	data.Description = resourceModel.Description
	data.Prefix = resourceModel.Prefix
	data.Address = resourceModel.Address
	data.Type = resourceModel.Type

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *SubnetDataSource) findSubnet(ctx context.Context, data SubnetDataSourceModel) (*kilhogsdk.Subnet, diag.Diagnostics) {
	var diags diag.Diagnostics

	networkID, err := uuid.Parse(data.NetworkID.ValueString())
	if err != nil {
		diags.AddError("Invalid network ID", err.Error())
		return nil, diags
	}

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
		subnetID, err := uuid.Parse(data.ID.ValueString())
		if err != nil {
			diags.AddError("Invalid subnet ID", err.Error())
			return nil, diags
		}

		subnet, err := d.client.GetSubnet(ctx, networkID, subnetID)
		if err != nil {
			diags.AddError("Error reading subnet", apiErrorSummary(err, "Read subnet"))
			return nil, diags
		}

		return subnet, diags
	default:
		subnets, err := d.client.ListSubnets(ctx, networkID)
		if err != nil {
			diags.AddError("Error listing subnets", apiErrorSummary(err, "List subnets"))
			return nil, diags
		}

		name := data.Name.ValueString()
		for i := range subnets {
			if subnets[i].Name == name {
				return &subnets[i], diags
			}
		}

		diags.AddError(
			"Subnet not found",
			fmt.Sprintf("No subnet named %q was found in network %q.", name, data.NetworkID.ValueString()),
		)
		return nil, diags
	}
}
