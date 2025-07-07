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

  "github.com/platylab/terraform-provider-unbound/internal/ssh"
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
      "username": schema.StringAttribute{
        Optional: true,
      },
      "private_key_path": schema.StringAttribute{
        Optional:  true,
      },
    },
  }
}

// unboundProviderModel maps provider schema data to a Go type.
type unboundProviderModel struct {
  Host           types.String `tfsdk:"host"`
  Username       types.String `tfsdk:"username"`
  PrivateKeyPath types.String `tfsdk:"private_key_path"`
}

// Configure prepares a unbound SSH client for data sources and resources.
func (p *unboundProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
  tflog.Info(ctx, "Configuring Unbound SSH client")
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
      "The provider cannot create the unbound SSH client as there is an unknown configuration value for the unbound SSH host. "+
        "Either target apply the source of the value first, set the value statically in the configuration, or use the UNBOUND_HOST environment variable.",
    )
  }

  if config.Username.IsUnknown() {
    resp.Diagnostics.AddAttributeError(
      path.Root("username"),
      "Unknown unbound SSH Username",
      "The provider cannot create the unbound SSH client as there is an unknown configuration value for the unbound SSH username. "+
        "Either target apply the source of the value first, set the value statically in the configuration, or use the UNBOUND_USERNAME environment variable.",
    )
  }

  if config.PrivateKeyPath.IsUnknown() {
    resp.Diagnostics.AddAttributeError(
      path.Root("private_key_path"),
      "Unknown unbound SSH Private Key",
      "The provider cannot create the unbound SSH client as there is an unknown configuration value for the unbound SSH private key. "+
        "Either target apply the source of the value first, set the value statically in the configuration, or use the UNBOUND_PRIVATE_KEY_PATH environment variable.",
    )
  }

  if resp.Diagnostics.HasError() {
    return
  }

  // Default values to environment variables, but override
  // with Terraform configuration value if set.

  host := os.Getenv("UNBOUND_HOST")
  username := os.Getenv("UNBOUND_USERNAME")
  private_key_path := os.Getenv("UNBOUND_PRIVATE_KEY_PATH")

  if !config.Host.IsNull() {
    host = config.Host.ValueString()
  }

  if !config.Username.IsNull() {
    username = config.Username.ValueString()
  }

  if !config.PrivateKeyPath.IsNull() {
    private_key_path = config.PrivateKeyPath.ValueString()
  }

  // If any of the expected configurations are missing, return
  // errors with provider-specific guidance.

  if host == "" {
    resp.Diagnostics.AddAttributeError(
      path.Root("host"),
      "Missing unbound SSH Host",
      "The provider cannot create the unbound SSH client as there is a missing or empty value for the unbound SSH host. "+
        "Set the host value in the configuration or use the UNBOUND_HOST environment variable. "+
        "If either is already set, ensure the value is not empty.",
    )
  }

  if username == "" {
    resp.Diagnostics.AddAttributeError(
      path.Root("username"),
      "Missing unbound SSH Username",
      "The provider cannot create the unbound SSH client as there is a missing or empty value for the unbound SSH username. "+
        "Set the username value in the configuration or use the UNBOUND_USERNAME environment variable. "+
        "If either is already set, ensure the value is not empty.",
    )
  }

  if private_key_path == "" {
    resp.Diagnostics.AddAttributeError(
      path.Root("private_key_path"),
      "Missing unbound SSH Private Key",
      "The provider cannot create the unbound SSH client as there is a missing or empty value for the unbound SSH Private Key. "+
        "Set the password value in the configuration or use the UNBOUND_PRIVATE_KEY_PATH environment variable. "+
        "If either is already set, ensure the value is not empty.",
    )
  }

  if resp.Diagnostics.HasError() {
    return
  }

  ctx = tflog.SetField(ctx, "unbound_host", host)
  ctx = tflog.SetField(ctx, "unbound_username", username)
  ctx = tflog.SetField(ctx, "unbound_private_key_path", private_key_path)

  tflog.Debug(ctx, "Creating Unbound SSH client")


  // Create a new ssh client using the configuration values
  client, err := ssh.NewSSHClient(host, username, private_key_path)
  if err != nil {
    resp.Diagnostics.AddError(
      "Unable to Create unbound SSH Client",
      "An unexpected error occurred when creating the unbound SSH client. "+
        "If the error is not clear, please contact the provider developers.\n\n"+
        "unbound Client Error: "+err.Error(),
    )
    return
  }

  // Make the unbound client available during DataSource and Resource
  // type Configure methods.
  resp.DataSourceData = client
  resp.ResourceData = client

  tflog.Info(ctx, "Configured Unbound SSH client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *unboundProvider) DataSources(_ context.Context) []func() datasource.DataSource {
  return []func() datasource.DataSource {
    NewLocalZoneDataSource,
  }
}

// Resources defines the resources implemented in the provider.
func (p *unboundProvider) Resources(_ context.Context) []func() resource.Resource {
  return nil
}
