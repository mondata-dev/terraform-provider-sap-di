package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/mondata-dev/terraform-provider-sap-di/internal/sap_di"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &sapDiProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &sapDiProvider{
			version: version,
		}
	}
}

// sapDiProvider is the provider implementation.
type sapDiProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// sapDiProviderModel maps provider schema data to a Go type.
type sapDiProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// Metadata returns the provider type name.
func (p *sapDiProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sapDi"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *sapDiProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with SAP DI",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional:    true,
				Description: "URI for SAP DI. May also be provided via SAP_DI_HOST environment variable.",
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Description: "Username for SAP DI. May also be provided via SAP_DI_USERNAME environment variable.",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Password for SAP DI. May also be provided via SAP_DI_PASSWORD environment variable.",
			},
		},
	}
}

// Configure prepares a SAP DI API client for data sources and resources.
func (p *sapDiProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring SAP DI client")

	// Retrieve provider data from configuration
	var config sapDiProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown SAP DI API Host",
			"The provider cannot create the SAP DI API client as there is an unknown configuration value for the SAP DI API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SAP_DI_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown SAP DI API Username",
			"The provider cannot create the SAP DI API client as there is an unknown configuration value for the SAP DI API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SAP_DI_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown SAP DI API Password",
			"The provider cannot create the SAP DI API client as there is an unknown configuration value for the SAP DI API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SAP_DI_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("SAP_DI_HOST")
	username := os.Getenv("SAP_DI_USERNAME")
	password := os.Getenv("SAP_DI_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing SAP DI API Host",
			"The provider cannot create the SAP DI API client as there is a missing or empty value for the SAP DI API host. "+
				"Set the host value in the configuration or use the SAP_DI_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing SAP DI API Username",
			"The provider cannot create the SAP DI API client as there is a missing or empty value for the SAP DI API username. "+
				"Set the username value in the configuration or use the SAP_DI_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing SAP DI API Password",
			"The provider cannot create the SAP DI API client as there is a missing or empty value for the SAP DI API password. "+
				"Set the password value in the configuration or use the SAP_DI_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "sapDi_host", host)
	ctx = tflog.SetField(ctx, "sapDi_username", username)
	ctx = tflog.SetField(ctx, "sapDi_password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "sapDi_password")

	tflog.Debug(ctx, "Creating SAP DI client")

	// Create a new SAP DI client using the configuration values
	client, err := sap_di.NewClient(&host, &username, &password)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create SAP DI API Client",
			"An unexpected error occurred when creating the SAP DI API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"SAP DI Client Error: "+err.Error(),
		)
		return
	}

	// Make the SAP DI client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured SAP DI client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *sapDiProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *sapDiProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
