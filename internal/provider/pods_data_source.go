package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSourceWithConfigure = &podsDataSource{}
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
	Id        types.String `tfsdk:"id"`
	Namespace types.String `tfsdk:"namespace"`
	Pods      []podsModel  `tfsdk:"pods"`
}

type podsModel struct {
	Name types.String `tfsdk:"name"`
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
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Example identifier",
			},
			"pods": schema.ListNestedAttribute{
				MarkdownDescription: "k8s pods list",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *podsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(kubernetes.Interface)
}

func (d *podsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state podsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	namespace := state.Namespace.ValueString()

	pods, err := d.client.CoreV1().Pods(namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		tflog.Error(ctx, err.Error())
		return
	}

	for _, p := range pods.Items {
		state.Pods = append(state.Pods, podsModel{Name: types.StringValue(p.Name)})
	}
	state.Id = types.StringValue("example-id")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
