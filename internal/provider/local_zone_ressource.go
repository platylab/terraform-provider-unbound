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
  _ resource.Resource                = &local_zoneResource{}
  _ resource.ResourceWithConfigure   = &local_zoneResource{}
	_ resource.ResourceWithImportState = &local_zoneResource{}

)

// NewLocalZoneResource is a helper function to simplify the provider implementation.
func NewLocalZoneResource() resource.Resource {
  return &local_zoneResource{}
}

// local_zoneResource is the resource implementation.
type local_zoneResource struct{
	client *unboundapiclient.Client
}

// Metadata returns the resource type name.
func (r *local_zoneResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
  resp.TypeName = req.ProviderTypeName + "_local_zone"
}

// Configure adds the provider configured client to the resource.
func (r *local_zoneResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


// local_zoneRessourceModel maps the data source schema data.
type local_zoneRessourceModel struct {
  Id          types.Int64    `tfsdk:"id"`
  Name        types.String   `tfsdk:"name"`
  Type        types.String   `tfsdk:"type"`
  LastUpdated types.String   `tfsdk:"last_updated"`

}

// Schema defines the schema for the resource.
func (r *local_zoneResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
  resp.Schema = schema.Schema{
    // This description is used by the documentation generator and the language server.
    MarkdownDescription: "Configured local-zone ressource",

    Attributes: map[string]schema.Attribute{
      "id": schema.Int64Attribute{
        MarkdownDescription: "ID of the zone",
        Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
      },
      "name": schema.StringAttribute{
        MarkdownDescription: "Name of the zone",
        Required:            true,
      },
      "type": schema.StringAttribute{
        MarkdownDescription: "Type of the zone",
        Required:            true,
      },
      "last_updated": schema.StringAttribute{
        Computed:            true,
      },
    },
  }
}

// Create creates the resource and sets the initial Terraform state.
func (r *local_zoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan local_zoneRessourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
  if resp.Diagnostics.HasError() {
    return
  }

	// Arguments for the API call
	body := unboundapiclient.PostConfigClauseAttributeValueIdJSONRequestBody{
		Value: strptr(fmt.Sprintf("\"%s\" %s", plan.Name.ValueString(), plan.Type.ValueString())),
	}

	// Create new localzone
	localzone, err := r.client.CreateLocalZone(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating LocalZone",
			"Could not create LocalZone, unexpected error: "+err.Error(),
		)
		return
	}


	// Map response body to schema and populate Computed attribute values
	plan.Id          = types.Int64Value(int64(localzone.Id))
	plan.Name        = types.StringValue(localzone.Name)
	plan.Type        = types.StringValue(localzone.Type)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *local_zoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
  var state local_zoneRessourceModel
  diags := resp.State.Get(ctx, &state)
  resp.Diagnostics.Append(diags...)
  if resp.Diagnostics.HasError() {
    return
  }

	// Arguments for the API call
	valueId := int(state.Id.ValueInt64())
  
	// Read existing localzone
  localzone, err := r.client.GetLocalZone(ctx, valueId)
  if err != nil {
    resp.Diagnostics.AddError(
      "Unable to Get Local-Zone",
      err.Error(),
    )
    return
  }

	// Map response body to schema
  state.Id          = types.Int64Value(int64(localzone.Id))
  state.Name        = types.StringValue(localzone.Name)
  state.Type        = types.StringValue(localzone.Type)
  
  // Set the state
  diags = resp.State.Set(ctx, state)
  resp.Diagnostics.Append(diags...)
  if resp.Diagnostics.HasError() {
    return
  }

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *local_zoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan local_zoneRessourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
  if resp.Diagnostics.HasError() {
    return
  }

	// Arguments for the API call
	valueId := int(plan.Id.ValueInt64())
	body := unboundapiclient.PutConfigClauseAttributeValueIdJSONRequestBody{
		Value: strptr(fmt.Sprintf("\"%s\" %s", plan.Name.ValueString(), plan.Type.ValueString())),
	}

	// Update existing localzone
	localzone, err := r.client.UpdateLocalZone(ctx, valueId, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating LocalZone",
			"Could not update LocalZone, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.Id          = types.Int64Value(int64(localzone.Id))
	plan.Name        = types.StringValue(localzone.Name)
	plan.Type        = types.StringValue(localzone.Type)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *local_zoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
  var state local_zoneRessourceModel
  diags := resp.State.Get(ctx, &state)
  resp.Diagnostics.Append(diags...)
  if resp.Diagnostics.HasError() {
    return
  }

	// Arguments for the API call
	valueId := int(state.Id.ValueInt64())

	// Delete existing localzone
	_, err := r.client.DeleteLocalZone(ctx, valueId)
  if err != nil {
    resp.Diagnostics.AddError(
      "Unable to Delete LocalZone",
      "Could not delete LocalZone, unexpected error: "+err.Error(),
    )
    return
  }

}


// ImportState imports an existing ressource into the state
func (r *local_zoneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
