package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"k8s.io/client-go/kubernetes"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &podsDataSource{}

func NewPodsDataSource() datasource.DataSource {
	return &podsDataSource{}
}

// podsDataSource defines the data source implementation.
type podsDataSource struct {
	client kubernetes.Interface
}

// podsDataSourceModel describes the data source data model.
type podsDataSourceModel struct {
	Namespace      types.String `tfsdk:"namespace"`
}

func (d *podsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pods"
}

func (d *podsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Getting Pods data source",

		Attributes: map[string]schema.Attribute{
			"namespace": schema.StringAttribute{
				MarkdownDescription: "k8s namespace",
				Required:            true,
			},
		},
	}
}

func (d *podsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, err := CreateClusterClient()


	if err != nil {
		resp.Diagnostics.AddError(
			err.Error(),
			"couldn't initialize k8s client",
		)

		return
	}

	d.client = client
}

func (d *podsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data podsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := d.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	// data.Id = types.StringValue("example-id")

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
