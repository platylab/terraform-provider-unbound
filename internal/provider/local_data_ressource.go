package provider

import (
  "context"
	"fmt"
	"time"
	"github.com/hashicorp/terraform-plugin-framework/path"
  "github.com/hashicorp/terraform-plugin-framework/resource"
  "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
  "github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/platylab/terraform-provider-unbound/internal/unboundapiclient"
)

// Ensure the implementation satisfies the expected interfaces.
var (
  _ resource.Resource                = &local_dataResource{}
  _ resource.ResourceWithConfigure   = &local_dataResource{}
	_ resource.ResourceWithImportState = &local_dataResource{}

)

// NewLocalDataResource is a helper function to simplify the provider implementation.
func NewLocalDataResource() resource.Resource {
  return &local_dataResource{}
}

// local_dataResource is the resource implementation.
type local_dataResource struct{
	client *unboundapiclient.Client
}

// Metadata returns the resource type name.
func (r *local_dataResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
  resp.TypeName = req.ProviderTypeName + "_local_data"
}

// Configure adds the provider configured client to the resource.
func (r *local_dataResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}


// local_dataRessourceModel maps the data source schema data.
type local_dataRessourceModel struct {
  Id          types.Int64    `tfsdk:"id"`
  Domain      types.String   `tfsdk:"domain"`
  Type        types.String   `tfsdk:"type"`
  Value       types.String   `tfsdk:"value"`
  LastUpdated types.String   `tfsdk:"last_updated"`

}

// Schema defines the schema for the resource.
func (r *local_dataResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
  resp.Schema = schema.Schema{
    // This description is used by the documentation generator and the language server.
    MarkdownDescription: "Configured local-data ressource",

    Attributes: map[string]schema.Attribute{
      "id": schema.Int64Attribute{
        MarkdownDescription: "ID of the entry",
        Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
      },
      "domain": schema.StringAttribute{
        MarkdownDescription: "Domain of the entry",
        Required:            true,
      },
      "type": schema.StringAttribute{
        MarkdownDescription: "Type of the entry",
        Required:            true,
      },
      "value": schema.StringAttribute{
        MarkdownDescription: "Value of the entry",
        Required:            true,
      },
      "last_updated": schema.StringAttribute{
        Computed:            true,
      },
    },
  }
}

// Create creates the resource and sets the initial Terraform state.
func (r *local_dataResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan local_dataRessourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
  if resp.Diagnostics.HasError() {
    return
  }

	// Arguments for the API call
	body := unboundapiclient.PostConfigClauseAttributeValueIdJSONRequestBody{
		Value: strptr(fmt.Sprintf("\"%s IN %s %s\"", plan.Domain.ValueString(), plan.Type.ValueString(), plan.Value.ValueString())),
	}

	// Create new localdata
	localdata, err := r.client.CreateLocalData(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating LocalData",
			"Could not create LocalData, unexpected error: "+err.Error(),
		)
		return
	}


	// Map response body to schema and populate Computed attribute values
	plan.Id          = types.Int64Value(int64(localdata.Id))
	plan.Domain      = types.StringValue(localdata.Domain)
	plan.Type        = types.StringValue(localdata.Type)
	plan.Value       = types.StringValue(localdata.Value)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *local_dataResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
  var state local_dataRessourceModel
  diags := resp.State.Get(ctx, &state)
  resp.Diagnostics.Append(diags...)
  if resp.Diagnostics.HasError() {
    return
  }

	// Arguments for the API call
	valueId := int(state.Id.ValueInt64())
  
	// Read existing localdata
  localdata, err := r.client.GetLocalData(ctx, valueId)
  if err != nil {
    resp.Diagnostics.AddError(
      "Unable to Get LocalData",
      err.Error(),
    )
    return
  }

	// Map response body to schema
  state.Id          = types.Int64Value(int64(localdata.Id))
	state.Domain      = types.StringValue(localdata.Domain)
	state.Type        = types.StringValue(localdata.Type)
	state.Value       = types.StringValue(localdata.Value)
  
  // Set the state
  diags = resp.State.Set(ctx, state)
  resp.Diagnostics.Append(diags...)
  if resp.Diagnostics.HasError() {
    return
  }

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *local_dataResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan local_dataRessourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
  if resp.Diagnostics.HasError() {
    return
  }

	// Arguments for the API call
	valueId := int(plan.Id.ValueInt64())
	body := unboundapiclient.PutConfigClauseAttributeValueIdJSONRequestBody{
		Value: strptr(fmt.Sprintf("\"%s IN %s %s\"", plan.Domain.ValueString(), plan.Type.ValueString(), plan.Value.ValueString())),
	}

	// Update existing localdata
	localdata, err := r.client.UpdateLocalData(ctx, valueId, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating LocalData",
			"Could not update LocalData, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.Id          = types.Int64Value(int64(localdata.Id))
	plan.Domain      = types.StringValue(localdata.Domain)
	plan.Type        = types.StringValue(localdata.Type)
	plan.Value       = types.StringValue(localdata.Value)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *local_dataResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
  var state local_dataRessourceModel
  diags := resp.State.Get(ctx, &state)
  resp.Diagnostics.Append(diags...)
  if resp.Diagnostics.HasError() {
    return
  }

	// Arguments for the API call
	valueId := int(state.Id.ValueInt64())

	// Delete existing localdata
	_, err := r.client.DeleteLocalData(ctx, valueId)
  if err != nil {
    resp.Diagnostics.AddError(
      "Unable to Delete LocalData",
      "Could not delete LocalData, unexpected error: "+err.Error(),
    )
    return
  }

}


// ImportState imports an existing ressource into the state
func (r *local_dataResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
