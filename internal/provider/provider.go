// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
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

  "github.com/platylab/terraform-provider-unbound/internal/unboundapiclient"
)

// Ensure the implementation satisfies the expected interfaces.
var (
  _ provider.Provider = &unboundProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
  return func() provider.Provider {
    return &unboundProvider{
      version: version,
    }
  }
}

// unboundProvider is the provider implementation.
type unboundProvider struct {
  // version is set to the provider version on release, "dev" when the
  // provider is built and ran locally, and "test" when running acceptance
  // testing.
  version string
}

// Metadata returns the provider type name.
func (p *unboundProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
  resp.TypeName = "unbound"
  resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *unboundProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
  resp.Schema = schema.Schema{
    Attributes: map[string]schema.Attribute{
      "host": schema.StringAttribute{
        Optional: true,
      },
    },
  }
}

// unboundProviderModel maps provider schema data to a Go type.
type unboundProviderModel struct {
  Host           types.String `tfsdk:"host"`
}

// Configure prepares a unbound API client for data sources and resources.
func (p *unboundProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
  tflog.Info(ctx, "Configuring Unbound API client")
  // Retrieve provider data from configuration
  var config unboundProviderModel
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
      "Unknown unbound Host",
      "The provider cannot create the unbound API client as there is an unknown configuration value for the unbound API host. "+
        "Either target apply the source of the value first, set the value statically in the configuration, or use the UNBOUND_HOST environment variable.",
    )
  }

  if resp.Diagnostics.HasError() {
    return
  }

  // Default values to environment variables, but override
  // with Terraform configuration value if set.

  host := os.Getenv("UNBOUND_HOST")

  if !config.Host.IsNull() {
    host = config.Host.ValueString()
  }

  // If any of the expected configurations are missing, return
  // errors with provider-specific guidance.

  if host == "" {
    resp.Diagnostics.AddAttributeError(
      path.Root("host"),
      "Missing unbound API Host",
      "The provider cannot create the unbound API client as there is a missing or empty value for the unbound API host. "+
        "Set the host value in the configuration or use the UNBOUND_HOST environment variable. "+
        "If either is already set, ensure the value is not empty.",
    )
  }

  if resp.Diagnostics.HasError() {
    return
  }

  ctx = tflog.SetField(ctx, "unbound_host", host)

  tflog.Debug(ctx, "Creating Unbound API client")


  // Create a new API client using the configuration values
  client, err := unboundapiclient.NewClient(host)
  if err != nil {
    resp.Diagnostics.AddError(
      "Unable to Create unbound API Client",
      "An unexpected error occurred when creating the unbound API client. "+
        "If the error is not clear, please contact the provider developers.\n\n"+
        "unbound Client Error: "+err.Error(),
    )
    return
  }

  // Make the unbound client available during DataSource and Resource
  // type Configure methods.
  resp.DataSourceData = client
  resp.ResourceData = client

  tflog.Info(ctx, "Configured Unbound API client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *unboundProvider) DataSources(_ context.Context) []func() datasource.DataSource {
  return []func() datasource.DataSource {
  }
}

// Resources defines the resources implemented in the provider.
func (p *unboundProvider) Resources(_ context.Context) []func() resource.Resource {
  return nil
}
