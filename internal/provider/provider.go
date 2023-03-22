package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure K8sProvider satisfies various provider interfaces.
var _ provider.Provider = &K8sProvider{}

// K8sProvider defines the provider implementation.
type K8sProvider struct {
}

func (p *K8sProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "k8s"
}

func (p *K8sProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *K8sProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	client, err := CreateClusterClient(ctx)

	if err != nil {
		resp.Diagnostics.AddError(
			err.Error(),
			"couldn't initialize k8s client",
		)

		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
	tflog.Info(ctx, "Configured HashiCups client", map[string]any{"success": true})
}

func (p *K8sProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

func (p *K8sProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPodsDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &K8sProvider{}
	}
}
