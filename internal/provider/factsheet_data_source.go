package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/mondata-dev/terraform-provider-sap-di/internal/sap_di"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &factsheetDataSource{}
	_ datasource.DataSourceWithConfigure = &factsheetDataSource{}
)

// NewFactsheetDataSource is a helper function to simplify the provider implementation.
func NewFactsheetDataSource() datasource.DataSource {
	return &factsheetDataSource{}
}

// factsheetDataSource is the data source implementation.
type factsheetDataSource struct {
	client *sap_di.Client
}

// Configure adds the provider configured client to the data source.
func (d *factsheetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Info(ctx, "Configuring SAP DI Factsheet data source")

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sap_di.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *sap_di.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client

	tflog.Info(ctx, "Configured SAP DI Factsheet data source", map[string]any{"success": true})
}

// Metadata returns the data source type name.
func (d *factsheetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_factsheet"
}

// Schema defines the schema for the data source.
func (d *factsheetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a factsheet.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier attribute.",
				Computed:    true,
			},

			"metadata": schema.SingleNestedAttribute{
				Description: "Metadata of the factsheet.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "Name of the factsheet.",
						Computed:    true,
					},
					"uri": schema.StringAttribute{
						Description: "URI for the factsheet.",
						Required:    true,
					},
					"connection_id": schema.StringAttribute{
						Description: "Connection ID for the factsheet.",
						Required:    true,
					},
				},
			},
		},
	}
}

// factsheetDataSourceModel maps the data source schema data.
type factsheetDataSourceModel struct {
	ID       types.String           `tfsdk:"id"`
	Metadata factsheetMetadataModel `tfsdk:"metadata"`
}

// factsheetModel maps factsheet schema data.
type factsheetMetadataModel struct {
	Name         types.String `tfsdk:"name"`
	Uri          types.String `tfsdk:"uri"`
	ConnectionId types.String `tfsdk:"connection_id"`
}

// Read refreshes the Terraform state with the latest data.
func (d *factsheetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state factsheetDataSourceModel

	// log req
	req.Config.Get(ctx, &state)

	tflog.Info(ctx, "Reading SAP DI Factsheet data source", map[string]any{
		"input": fmt.Sprintf("%+v", state),
	})

	// diags := req.State.Get(ctx, &state)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	factsheet, err := d.client.GetFactsheet(
		state.Metadata.ConnectionId.ValueString(),
		state.Metadata.Uri.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read SAP DI factsheets",
			err.Error(),
		)
		return
	}

	// Map response body to model
	state.Metadata = factsheetMetadataModel{
		Name:         types.StringValue(factsheet.Metadata.Name),
		Uri:          types.StringValue(factsheet.Metadata.Uri),
		ConnectionId: types.StringValue(factsheet.Metadata.ConnectionId),
	}

	state.ID = types.StringValue("placeholder")

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
