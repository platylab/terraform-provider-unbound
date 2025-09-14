
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
  _ datasource.DataSource              = &local_dataDataSource{}
  _ datasource.DataSourceWithConfigure = &local_dataDataSource{}
)

// NewLocalDataDataSource is a helper function to simplify the provider implementation.
func NewLocalDataDataSource() datasource.DataSource {
  return &local_dataDataSource{}
}

// local_dataDataSource is the data source implementation.
type local_dataDataSource struct {
  client *unboundapiclient.Client
}

// Metadata returns the data source type name.
func (d *local_dataDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
  resp.TypeName = req.ProviderTypeName + "_local_data"
}

// Configure adds the provider configured client to the data source.
func (d *local_dataDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// local_dataDataSourceModel maps the data source schema data.
type local_dataDataSourceModel struct {
  Id     types.Int64    `tfsdk:"id"`
  Domain types.String   `tfsdk:"domain"`
  Type   types.String   `tfsdk:"type"`
  Value  types.String   `tfsdk:"value"`
}

// Schema defines the schema for the data source.
func (d *local_dataDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
  resp.Schema = schema.Schema{
    // This description is used by the documentation generator and the language server.
    MarkdownDescription: "Configured local-data data source",

    Attributes: map[string]schema.Attribute{
      "id": schema.Int64Attribute{
        MarkdownDescription: "ID of the entry",
        Required:            true,
      },
      "domain": schema.StringAttribute{
        MarkdownDescription: "Domain of the entry",
        Computed:            true,
      },
      "type": schema.StringAttribute{
        MarkdownDescription: "Type of the entry",
        Computed:            true,
      },
      "value": schema.StringAttribute{
        MarkdownDescription: "Value of the entry",
        Computed:            true,
      },
    },
  }
}


// Read refreshes the Terraform state with the latest data.
func (d *local_dataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
  // Initialize state
  var state local_dataDataSourceModel
  diags := resp.State.Get(ctx, &state)
  resp.Diagnostics.Append(diags...)
  if resp.Diagnostics.HasError() {
    return
  }
  
  // Arguments for the API call
	valueId := int(state.Id.ValueInt64())
  
  localdata, err := d.client.GetLocalData(ctx, valueId)
  if err != nil {
    resp.Diagnostics.AddError(
      "Unable to Get Local-Data",
      err.Error(),
    )
    return
  }


  // Set state values based on parsed response
  state.Id     = types.Int64Value(int64(localdata.Id))
  state.Domain = types.StringValue(localdata.Domain)
  state.Type   = types.StringValue(localdata.Type)
  state.Value  = types.StringValue(localdata.Value)
  
  // Set the state
  diags = resp.State.Set(ctx, state)
  resp.Diagnostics.Append(diags...)
  if resp.Diagnostics.HasError() {
    return
  }

}
