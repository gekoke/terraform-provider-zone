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
	_ resource.Resource                = &recordCNAMEResource{}
	_ resource.ResourceWithConfigure   = &recordCNAMEResource{}
	_ resource.ResourceWithImportState = &recordCNAMEResource{}
)

func NewRecordCNAMEResource() resource.Resource {
	return &recordCNAMEResource{}
}

type recordCNAMEResource struct {
	client api.Client
}

type recordCNAMEResourceModel struct {
	ID          types.String `tfsdk:"id"`
	LastUpdated types.String `tfsdk:"last_updated"`
	Domain      types.String `tfsdk:"domain"`
	Name        types.String `tfsdk:"name"`
	Destination types.String `tfsdk:"destination"`
	ResourceURL types.String `tfsdk:"resource_url"`
	Modify      types.Bool   `tfsdk:"modify"`
	Delete      types.Bool   `tfsdk:"delete"`
}

func (*recordCNAMEResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_record_cname"
}

func (*recordCNAMEResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
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

func (resource *recordCNAMEResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (recordCNAMEResource) ImportState(context context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(context, path.Root("id"), request, response)
}

func (resource *recordCNAMEResource) Create(context context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan recordCNAMEResourceModel

	diagnostics := request.Plan.Get(context, &plan)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	destination := plan.Destination.ValueString()
	domain := plan.Domain.ValueString()
	name := plan.Name.ValueString()
	record := api.CNAMERecord{Name: name, Destination: destination}

	recordInfo, err := resource.client.CreateCNAMERecord(domain, record)

	if err != nil {
		response.Diagnostics.AddError("Error creating CNAME record", "Request failed: "+err.Error())
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

func (resource *recordCNAMEResource) Read(context context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state recordCNAMEResourceModel
	diagnostics := request.State.Get(context, &state)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	domain := api.Domain(state.Domain.ValueString())
	id := api.Identificator(state.ID.ValueString())

	cnameRecord, err := resource.client.GetCNAMERecord(domain, id)

	if err != nil {
		response.Diagnostics.AddError("Error reading CNAME record", "Request failed: "+err.Error())
		return
	}

	state.ID = types.StringValue(cnameRecord.ID)
	state.Name = types.StringValue(cnameRecord.Name)
	state.Destination = types.StringValue(cnameRecord.Destination)
	state.ResourceURL = types.StringValue(cnameRecord.ResourceURL.String())
	state.Modify = types.BoolValue(cnameRecord.Modify)
	state.Delete = types.BoolValue(cnameRecord.Delete)

	diagnostics = response.State.Set(context, &state)
	response.Diagnostics.Append(diagnostics...)
}

func (resource *recordCNAMEResource) Update(context context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan recordCNAMEResourceModel
	var state recordCNAMEResourceModel

	diagnostics := request.Plan.Get(context, &plan)
	response.Diagnostics.Append(diagnostics...)
	diagnostics = request.State.Get(context, &state)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	destination := plan.Destination.ValueString()
	domain := plan.Domain.ValueString()
	recordID := state.ID.ValueString()
	record := api.CNAMERecord{Name: plan.Name.ValueString(), Destination: destination}

	recordInfo, err := resource.client.UpdateCNAMERecord(domain, recordID, record)
	if err != nil {
		response.Diagnostics.AddError("Error updating CNAME record", "Request failed: "+err.Error())
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

func (resource *recordCNAMEResource) Delete(context context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state recordCNAMEResourceModel
	diagnostics := request.State.Get(context, &state)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	domain := state.Domain.ValueString()
	id := state.ID.ValueString()

	err := resource.client.DeleteCNAMERecord(domain, id)

	if err != nil {
		response.Diagnostics.AddError("Error deleting CNAME record", "Request failed: "+err.Error())
	}
}
