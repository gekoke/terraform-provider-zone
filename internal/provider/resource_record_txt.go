package provider

import (
	"context"
	"time"

	"github.com/gekoke/terraform-provider-zone/internal/api"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &recordTXTResource{}
	_ resource.ResourceWithConfigure   = &recordTXTResource{}
	_ resource.ResourceWithImportState = &recordTXTResource{}
)

func NewRecordTXTResource() resource.Resource {
	return &recordTXTResource{}
}

type recordTXTResource struct {
	client api.Client
}

type recordTXTResourceModel struct {
	ID          types.String `tfsdk:"id"`
	LastUpdated types.String `tfsdk:"last_updated"`
	Domain      types.String `tfsdk:"domain"`
	Name        types.String `tfsdk:"name"`
	Destination types.String `tfsdk:"destination"`
	ResourceURL types.String `tfsdk:"resource_url"`
	Modify      types.Bool   `tfsdk:"modify"`
	Delete      types.Bool   `tfsdk:"delete"`
}

func (*recordTXTResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_record_txt"
}

func (*recordTXTResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"domain": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"destination": schema.StringAttribute{
				Required: true,
			},
			"resource_url": schema.StringAttribute{
				Computed: true,
			},
			"delete": schema.BoolAttribute{
				Computed: true,
			},
			"modify": schema.BoolAttribute{
				Computed: true,
			},
		},
	}
}

func (resource *recordTXTResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(api.Client)

	if !ok {
		response.Diagnostics.AddError("Unexpected Data Source Configure Type", "Please report this issue to the provider developers.")
		return
	}

	resource.client = client
}

func (*recordTXTResource) ImportState(context context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(context, path.Root("id"), request, response)
}

func (resource *recordTXTResource) Create(context context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan recordTXTResourceModel

	diagnostics := request.Plan.Get(context, &plan)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	destination := plan.Destination.ValueString()
	domain := plan.Domain.ValueString()
	name := plan.Name.ValueString()
	record := api.TXTRecord{Name: name, Destination: destination}

	recordInfo, err := resource.client.CreateTXTRecord(domain, record)

	if err != nil {
		response.Diagnostics.AddError("Error creating TXT record", "Request failed: "+err.Error())
		return
	}

	plan.ID = types.StringValue(recordInfo.ID)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.Name = types.StringValue(recordInfo.Name)
	plan.Destination = types.StringValue(recordInfo.Destination)
	plan.ResourceURL = types.StringValue(recordInfo.ResourceURL.String())
	plan.Modify = types.BoolValue(recordInfo.Modify)
	plan.Delete = types.BoolValue(recordInfo.Delete)

	diagnostics = response.State.Set(context, plan)
	response.Diagnostics.Append(diagnostics...)
}

func (resource *recordTXTResource) Read(context context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state recordTXTResourceModel
	diagnostics := request.State.Get(context, &state)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	domain := api.Domain(state.Domain.ValueString())
	id := api.Identificator(state.ID.ValueString())

	recordInfo, err := resource.client.GetTXTRecord(domain, id)

	if err != nil {
		response.Diagnostics.AddError("Error reading TXT record", "Request failed: "+err.Error())
		return
	}

	state.ID = types.StringValue(recordInfo.ID)
	state.Name = types.StringValue(recordInfo.Name)
	state.Destination = types.StringValue(recordInfo.Destination)
	state.ResourceURL = types.StringValue(recordInfo.ResourceURL.String())
	state.Modify = types.BoolValue(recordInfo.Modify)
	state.Delete = types.BoolValue(recordInfo.Delete)

	diagnostics = response.State.Set(context, &state)
	response.Diagnostics.Append(diagnostics...)
}

func (resource *recordTXTResource) Update(context context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan recordTXTResourceModel
	var state recordTXTResourceModel

	diagnostics := request.Plan.Get(context, &plan)
	response.Diagnostics.Append(diagnostics...)
	diagnostics = request.State.Get(context, &state)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	domain := plan.Domain.ValueString()
	destination := plan.Destination.ValueString()
	recordID := state.ID.ValueString()
	record := api.TXTRecord{Name: plan.Name.ValueString(), Destination: destination}

	recordInfo, err := resource.client.UpdateTXTRecord(domain, recordID, record)
	if err != nil {
		response.Diagnostics.AddError("Error updating TXT record", "Request failed: "+err.Error())
		return
	}

	plan.ID = types.StringValue(recordInfo.ID)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.Name = types.StringValue(recordInfo.Name)
	plan.Destination = types.StringValue(recordInfo.Destination)
	plan.ResourceURL = types.StringValue(recordInfo.ResourceURL.String())
	plan.Modify = types.BoolValue(recordInfo.Modify)
	plan.Delete = types.BoolValue(recordInfo.Delete)

	diagnostics = response.State.Set(context, plan)
	response.Diagnostics.Append(diagnostics...)
}

func (resource *recordTXTResource) Delete(context context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state recordTXTResourceModel
	diagnostics := request.State.Get(context, &state)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	domain := state.Domain.ValueString()
	id := state.ID.ValueString()

	err := resource.client.DeleteTXTRecord(domain, id)

	if err != nil {
		response.Diagnostics.AddError("Error deleting TXT record", "Request failed: "+err.Error())
	}
}
