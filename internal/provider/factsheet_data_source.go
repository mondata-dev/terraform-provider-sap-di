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
	descriptionsObj := schema.ListNestedAttribute{
		Description: "Descriptions of the factsheet.",
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"origin": schema.StringAttribute{
					Computed: true,
				},
				"type": schema.StringAttribute{
					Computed: true,
				},
				"value": schema.StringAttribute{
					Computed: true,
				},
			},
		},
	}

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
					"descriptions": descriptionsObj,
				},
			},

			"columns": schema.ListNestedAttribute{
				Description: "Columns of the factsheet.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the column.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Type of the column.",
							Computed:    true,
						},
						"descriptions": descriptionsObj,
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
	Columns  []factsheetColumnModel `tfsdk:"columns"`
}

// factsheetModel maps factsheet schema data.
type factsheetMetadataModel struct {
	Name         types.String                `tfsdk:"name"`
	Uri          types.String                `tfsdk:"uri"`
	ConnectionId types.String                `tfsdk:"connection_id"`
	Descriptions []factsheetDescriptionModel `tfsdk:"descriptions"`
}

type factsheetColumnModel struct {
	Name         types.String                `tfsdk:"name"`
	Type         types.String                `tfsdk:"type"`
	Descriptions []factsheetDescriptionModel `tfsdk:"descriptions"`
}

// factsheetDescriptionModel maps factsheet description schema data.
type factsheetDescriptionModel struct {
	Origin types.String `tfsdk:"origin"`
	Type   types.String `tfsdk:"type"`
	Value  types.String `tfsdk:"value"`
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
		Descriptions: []factsheetDescriptionModel{},
	}

	for _, desc := range factsheet.Metadata.Descriptions {
		state.Metadata.Descriptions = append(state.Metadata.Descriptions, factsheetDescriptionModel{
			Origin: types.StringValue(desc.Origin),
			Type:   types.StringValue(desc.Type),
			Value:  types.StringValue(desc.Value),
		})
	}

	for _, column := range factsheet.Columns {
		col := factsheetColumnModel{
			Name:         types.StringValue(column.Name),
			Type:         types.StringValue(column.Type),
			Descriptions: []factsheetDescriptionModel{},
		}

		for _, desc := range column.Descriptions {
			col.Descriptions = append(col.Descriptions, factsheetDescriptionModel{
				Origin: types.StringValue(desc.Origin),
				Type:   types.StringValue(desc.Type),
				Value:  types.StringValue(desc.Value),
			})
		}

		state.Columns = append(state.Columns, col)
	}

	state.ID = types.StringValue("placeholder")

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
