// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure KilhogProvider satisfies various provider interfaces.
var _ provider.Provider = &KilhogProvider{}
var _ provider.ProviderWithFunctions = &KilhogProvider{}
var _ provider.ProviderWithEphemeralResources = &KilhogProvider{}
var _ provider.ProviderWithActions = &KilhogProvider{}

// KilhogProvider defines the provider implementation.
type KilhogProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// KilhogProviderModel describes the provider data model.
type KilhogProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *KilhogProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kilhog"
	resp.Version = p.version
}

func (p *KilhogProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Kilhog API endpoint",
				Optional:            true,
			},
		},
	}
}

func (p *KilhogProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data KilhogProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Client configuration for data sources and resources
	client := http.DefaultClient
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *KilhogProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewKilhogResource,
	}
}

func (p *KilhogProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		NewKilhogEphemeralResource,
	}
}

func (p *KilhogProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewKilhogDataSource,
	}
}

func (p *KilhogProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewKilhogFunction,
	}
}

func (p *KilhogProvider) Actions(ctx context.Context) []func() action.Action {
	return []func() action.Action{
		NewKilhogAction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &KilhogProvider{
			version: version,
		}
	}
}
