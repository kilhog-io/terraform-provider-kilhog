// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	kilhogsdk "github.com/kilhog-io/kilhog/pkg/kilhog"
)

var (
	_ resource.Resource                = &SubnetResource{}
	_ resource.ResourceWithImportState = &SubnetResource{}
)

func NewSubnetResource() resource.Resource {
	return &SubnetResource{}
}

type SubnetResource struct {
	client *kilhogsdk.Client
}

type SubnetResourceModel struct {
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

func (r *SubnetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subnet"
}

func (r *SubnetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Kilhog subnet within a network.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Subnet UUID.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "UUID of the network that owns this subnet.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"parent_subnet_id": schema.StringAttribute{
				MarkdownDescription: "UUID of the parent subnet. When set, the subnet is created as a child of that subnet instead of directly under the network.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"parent_kind": schema.StringAttribute{
				MarkdownDescription: "Parent kind reported by the API (`network` or `subnet`).",
				Computed:            true,
			},
			"parent_id": schema.StringAttribute{
				MarkdownDescription: "UUID of the parent resource reported by the API.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Subnet name, unique within the network.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Subnet description.",
				Optional:            true,
			},
			"prefix": schema.Int64Attribute{
				MarkdownDescription: "IPv4 prefix length (1-32).",
				Required:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"address": schema.StringAttribute{
				MarkdownDescription: "Network address for the subnet. Required when the subnet is created directly under a network. Optional when created under a parent subnet; the API allocates an address when omitted.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Address family. Only `ipv4` is supported by the API today.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(string(kilhogsdk.AddressTypeIPv4)),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringOneOf(string(kilhogsdk.AddressTypeIPv4), string(kilhogsdk.AddressTypeIPv6)),
				},
			},
		},
	}
}

func (r *SubnetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*kilhogsdk.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *kilhog.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *SubnetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SubnetResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	networkID, err := uuid.Parse(data.NetworkID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid network ID", err.Error())
		return
	}

	input := expandSubnetInput(data)

	var subnet *kilhogsdk.Subnet
	if !data.ParentSubnetID.IsNull() && !data.ParentSubnetID.IsUnknown() {
		parentSubnetID, parseErr := uuid.Parse(data.ParentSubnetID.ValueString())
		if parseErr != nil {
			resp.Diagnostics.AddError("Invalid parent subnet ID", parseErr.Error())
			return
		}

		subnet, err = r.client.CreateSubnetUnderParent(ctx, networkID, parentSubnetID, input)
	} else {
		if input.Address == "" {
			resp.Diagnostics.AddError(
				"Missing address",
				"`address` is required when `parent_subnet_id` is not set.",
			)
			return
		}

		subnet, err = r.client.CreateSubnetInNetwork(ctx, networkID, input)
	}

	if err != nil {
		resp.Diagnostics.AddError("Error creating subnet", apiErrorSummary(err, "Create subnet"))
		return
	}

	flattenSubnet(subnet, &data)

	tflog.Trace(ctx, "created a subnet", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SubnetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SubnetResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subnet, diags := r.readSubnet(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	flattenSubnet(subnet, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SubnetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SubnetResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	networkID, err := uuid.Parse(data.NetworkID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid network ID", err.Error())
		return
	}

	subnetID, err := uuid.Parse(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid subnet ID", err.Error())
		return
	}

	input := kilhogsdk.UpdateSubnetInput{}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		input.Description = data.Description.ValueString()
	}

	subnet, err := r.client.UpdateSubnet(ctx, networkID, subnetID, input)
	if err != nil {
		resp.Diagnostics.AddError("Error updating subnet", apiErrorSummary(err, "Update subnet"))
		return
	}

	flattenSubnet(subnet, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SubnetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SubnetResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	networkID, err := uuid.Parse(data.NetworkID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid network ID", err.Error())
		return
	}

	subnetID, err := uuid.Parse(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid subnet ID", err.Error())
		return
	}

	if err := r.client.DeleteSubnet(ctx, networkID, subnetID); err != nil {
		resp.Diagnostics.AddError("Error deleting subnet", apiErrorSummary(err, "Delete subnet"))
		return
	}
}

func (r *SubnetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Expected format: `<network_id>/<subnet_id>`",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

func (r *SubnetResource) readSubnet(ctx context.Context, data SubnetResourceModel) (*kilhogsdk.Subnet, diag.Diagnostics) {
	var diags diag.Diagnostics

	networkID, err := uuid.Parse(data.NetworkID.ValueString())
	if err != nil {
		diags.AddError("Invalid network ID", err.Error())
		return nil, diags
	}

	subnetID, err := uuid.Parse(data.ID.ValueString())
	if err != nil {
		diags.AddError("Invalid subnet ID", err.Error())
		return nil, diags
	}

	subnet, err := r.client.GetSubnet(ctx, networkID, subnetID)
	if err != nil {
		diags.AddError("Error reading subnet", apiErrorSummary(err, "Read subnet"))
		return nil, diags
	}

	return subnet, diags
}

func expandSubnetInput(data SubnetResourceModel) kilhogsdk.CreateSubnetInput {
	input := kilhogsdk.CreateSubnetInput{
		Name:   data.Name.ValueString(),
		Prefix: int(data.Prefix.ValueInt64()),
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		input.Description = data.Description.ValueString()
	}
	if !data.Address.IsNull() && !data.Address.IsUnknown() {
		input.Address = data.Address.ValueString()
	}
	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		input.Type = kilhogsdk.AddressType(data.Type.ValueString())
	}

	return input
}

func flattenSubnet(subnet *kilhogsdk.Subnet, data *SubnetResourceModel) {
	data.ID = types.StringValue(subnet.UUID.String())
	data.ParentKind = types.StringValue(string(subnet.Parent.Kind))
	data.ParentID = types.StringValue(subnet.Parent.UUID.String())

	if subnet.Parent.Kind == kilhogsdk.ParentKindSubnet {
		data.ParentSubnetID = types.StringValue(subnet.Parent.UUID.String())
	} else {
		data.ParentSubnetID = types.StringNull()
	}
	data.Name = types.StringValue(subnet.Name)

	if subnet.Description != "" {
		data.Description = types.StringValue(subnet.Description)
	} else {
		data.Description = types.StringNull()
	}

	data.Prefix = types.Int64Value(int64(subnet.Prefix))
	data.Address = types.StringValue(subnet.Address)
	data.Type = types.StringValue(string(subnet.Type))
}
