package provider

import (
	"context"

	"github.com/gekoke/terraform-provider-zone/internal/api"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ provider.Provider = &zoneProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &zoneProvider{
			version: version,
		}
	}
}

type zoneProvider struct {
	version string
}

type zoneProviderModel struct {
	Username types.String `tfsdk:"username"`
	APIKey   types.String `tfsdk:"api_key"`
}

func (provider zoneProvider) Metadata(_ context.Context, _ provider.MetadataRequest, response *provider.MetadataResponse) {
	response.TypeName = "zone"
	response.Version = provider.version
}

func (zoneProvider) Schema(_ context.Context, _ provider.SchemaRequest, response *provider.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{Required: true},
			"api_key":  schema.StringAttribute{Required: true, Sensitive: true},
		},
	}

}

func (zoneProvider) Configure(context context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	var config zoneProviderModel

	diagnostics := request.Config.Get(context, &config)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	if config.Username.IsUnknown() {
		response.Diagnostics.AddAttributeError(path.Root("username"), "Unknown username", "")
	}
	if config.APIKey.IsUnknown() {
		response.Diagnostics.AddAttributeError(path.Root("api_key"), "Unknown API key", "")
	}

	if response.Diagnostics.HasError() {
		return
	}

	if config.Username.IsNull() {
		response.Diagnostics.AddAttributeError(path.Root("username"), "Missing Zone username", "")
	}
	if config.APIKey.IsNull() {
		response.Diagnostics.AddAttributeError(path.Root("api_key"), "Missing Zone API key", "")
	}

	if response.Diagnostics.HasError() {
		return
	}

	const apiBaseURL = "https://api.zone.eu/v2/"
	client := api.MakeClient(apiBaseURL, config.Username.ValueString(), config.APIKey.ValueString())

	response.DataSourceData = client
	response.ResourceData = client
}

func (zoneProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (zoneProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRecordAResource,
	}
}
