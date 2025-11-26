package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/platylab/terraform-provider-unbound/internal/unboundapiclient"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &local_zoneDataSource{}
	_ datasource.DataSourceWithConfigure = &local_zoneDataSource{}
)

// NewLocalZoneDataSource is a helper function to simplify the provider implementation.
func NewLocalZoneDataSource() datasource.DataSource {
	return &local_zoneDataSource{}
}

// local_zoneDataSource is the data source implementation.
type local_zoneDataSource struct {
	client *unboundapiclient.Client
}

// Metadata returns the data source type name.
func (d *local_zoneDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_zone"
}

// Configure adds the provider configured client to the data source.
func (d *local_zoneDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*unboundapiclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *unboundapiclient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// local_zoneDataSourceModel maps the data source schema data.
type local_zoneDataSourceModel struct {
	Id   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

// Schema defines the schema for the data source.
func (d *local_zoneDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Configured local-zone data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "ID of the zone",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the zone",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of the zone",
				Computed:            true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *local_zoneDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Initialize state
	var state local_zoneDataSourceModel
	diags := resp.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Arguments for the API call
	valueId := int(state.Id.ValueInt64())

	localzone, err := d.client.GetLocalZone(ctx, valueId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Get Local-Zone",
			err.Error(),
		)
		return
	}

	// Set state values based on parsed response
	state.Id = types.Int64Value(int64(localzone.Id))
	state.Name = types.StringValue(localzone.Name)
	state.Type = types.StringValue(localzone.Type)

	// Set the state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
