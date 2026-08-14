// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	kilhogsdk "github.com/kilhog-io/kilhog/pkg/kilhog"
)

// Ensure KilhogProvider satisfies various provider interfaces.
var _ provider.Provider = &KilhogProvider{}

// KilhogProvider defines the provider implementation.
type KilhogProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// KilhogProviderModel describes the provider data model.
type KilhogProviderModel struct {
	BaseURL types.String `tfsdk:"base_url"`
	APIKey  types.String `tfsdk:"api_key"`
}

func (p *KilhogProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kilhog"
	resp.Version = p.version
}

func (p *KilhogProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Interact with the Kilhog IPAM API.",
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Kilhog API base URL. May also be set via the `KILHOG_BASE_URL` environment variable. Defaults to `http://localhost:8080`.",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "Kilhog API key sent as a Bearer token. May also be set via the `KILHOG_API_KEY` environment variable.",
				Optional:            true,
				Sensitive:           true,
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

	baseURL := data.BaseURL.ValueString()
	if data.BaseURL.IsNull() || data.BaseURL.IsUnknown() {
		if envBaseURL := os.Getenv("KILHOG_BASE_URL"); envBaseURL != "" {
			baseURL = envBaseURL
		}
	}

	apiKey := data.APIKey.ValueString()
	if data.APIKey.IsNull() || data.APIKey.IsUnknown() {
		apiKey = os.Getenv("KILHOG_API_KEY")
	}

	client, err := kilhogsdk.NewClient(kilhogsdk.ClientConfig{
		BaseURL: baseURL,
		APIKey:  apiKey,
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to create Kilhog client", err.Error())
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *KilhogProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewNetworkResource,
		NewSubnetResource,
	}
}

func (p *KilhogProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewNetworkDataSource,
		NewSubnetDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &KilhogProvider{
			version: version,
		}
	}
}
