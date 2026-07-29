// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ action.Action = &KilhogAction{}
var _ action.ActionWithConfigure = &KilhogAction{}

func NewKilhogAction() action.Action {
	return &KilhogAction{}
}

// KilhogAction defines the action implementation.
type KilhogAction struct {
	client *http.Client
}

// KilhogActionModel describes the action data model.
type KilhogActionModel struct {
	ConfigurableAttribute types.String `tfsdk:"configurable_attribute"`
}

func (e *KilhogAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_action"
}

func (e *KilhogAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Kilhog action",

		Attributes: map[string]schema.Attribute{
			"configurable_attribute": schema.StringAttribute{
				MarkdownDescription: "Configurable attribute",
				Optional:            true,
			},
		},
	}
}

func (e *KilhogAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	e.client = client
}

func (e *KilhogAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	// Send a progress message back to Terraform
	resp.SendProgress(action.InvokeProgressEvent{
		Message: "starting action invocation",
	})

	var data KilhogActionModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := d.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "invoke an action")

	// Send a progress message back to Terraform
	resp.SendProgress(action.InvokeProgressEvent{
		Message: "finished action invocation",
	})
}
