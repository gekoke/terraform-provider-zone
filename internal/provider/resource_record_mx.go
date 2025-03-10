package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/gekoke/terraform-provider-zone/internal/api"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &recordMXResource{}
	_ resource.ResourceWithConfigure   = &recordMXResource{}
	_ resource.ResourceWithImportState = &recordMXResource{}
)

func NewRecordMXResource() resource.Resource {
	return &recordMXResource{}
}

type recordMXResource struct {
	client api.Client
}

type recordMXResourceModel struct {
	ID          types.String `tfsdk:"id"`
	LastUpdated types.String `tfsdk:"last_updated"`
	Domain      types.String `tfsdk:"domain"`
	Name        types.String `tfsdk:"name"`
	Destination types.String `tfsdk:"destination"`
	Priority    types.Int32  `tfsdk:"priority"`
	ResourceURL types.String `tfsdk:"resource_url"`
	Modify      types.Bool   `tfsdk:"modify"`
	Delete      types.Bool   `tfsdk:"delete"`
}

func (*recordMXResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_record_mx"
}

func (*recordMXResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
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
			"priority": schema.Int32Attribute{
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

func (resource *recordMXResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (recordMXResource) ImportState(context context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(context, path.Root("id"), request, response)
}

func (resource *recordMXResource) Create(context context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	plan, err := getCreatePlan[recordMXResourceModel](context, request, response)
	if err != nil {
		return
	}

	priority := plan.Priority.ValueInt32()
	if !validateMXPriority(&response.Diagnostics, priority) {
		return
	}
	destination := plan.Destination.ValueString()
	domain := plan.Domain.ValueString()
	name := plan.Name.ValueString()
	record := api.MXRecord{Name: name, Destination: destination, Priority: uint16(priority)}

	recordInfo, err := resource.client.CreateMXRecord(domain, record)
	if err != nil {
		response.Diagnostics.AddError("Error creating MX record", "Request failed: "+err.Error())
		return
	}

	var newState recordMXResourceModel
	newState.ID = types.StringValue(recordInfo.ID)
	newState.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	newState.Name = types.StringValue(recordInfo.Name)
	newState.Destination = types.StringValue(recordInfo.Destination)
	newState.Priority = types.Int32Value(int32(recordInfo.Priority))
	newState.ResourceURL = types.StringValue(recordInfo.ResourceURL.String())
	newState.Modify = types.BoolValue(recordInfo.Modify)
	newState.Delete = types.BoolValue(recordInfo.Delete)

	diagnostics := response.State.Set(context, &newState)
	response.Diagnostics.Append(diagnostics...)
}

func (resource *recordMXResource) Read(context context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	state, err := getReadState[recordMXResourceModel](context, request, response)
	if err != nil {
		return
	}

	domain := api.Domain(state.Domain.ValueString())
	id := api.Identificator(state.ID.ValueString())

	mxRecord, err := resource.client.GetMXRecord(domain, id)
	if err != nil {
		response.Diagnostics.AddError("Error reading MX record", "Request failed: "+err.Error())
		return
	}

	var newState recordMXResourceModel
	newState.ID = types.StringValue(mxRecord.ID)
	newState.Name = types.StringValue(mxRecord.Name)
	newState.Destination = types.StringValue(mxRecord.Destination)
	newState.Priority = types.Int32Value(int32(mxRecord.Priority))
	newState.ResourceURL = types.StringValue(mxRecord.ResourceURL.String())
	newState.Modify = types.BoolValue(mxRecord.Modify)
	newState.Delete = types.BoolValue(mxRecord.Delete)

	diagnostics := response.State.Set(context, &newState)
	response.Diagnostics.Append(diagnostics...)
}

func (resource *recordMXResource) Update(context context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	plan, err := getUpdatePlan[recordMXResourceModel](context, request, response)
	if err != nil {
		return
	}
	state, err := getUpdateState[recordMXResourceModel](context, request, response)
	if err != nil {
		return
	}

	priority := plan.Priority.ValueInt32()
	if !validateMXPriority(&response.Diagnostics, priority) {
		return
	}
	name := plan.Name.ValueString()
	destination := plan.Destination.ValueString()
	domain := plan.Domain.ValueString()
	recordID := state.ID.ValueString()
	record := api.MXRecord{Name: name, Destination: destination, Priority: uint16(priority)}

	recordInfo, err := resource.client.UpdateMXRecord(domain, recordID, record)
	if err != nil {
		response.Diagnostics.AddError("Error updating MX record", "Request failed: "+err.Error())
		return
	}

	var newState recordMXResourceModel
	newState.ID = types.StringValue(recordInfo.ID)
	newState.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	newState.Name = types.StringValue(recordInfo.Name)
	newState.Destination = types.StringValue(recordInfo.Destination)
	newState.Priority = types.Int32Value(int32(record.Priority))
	newState.ResourceURL = types.StringValue(recordInfo.ResourceURL.String())
	newState.Modify = types.BoolValue(recordInfo.Modify)
	newState.Delete = types.BoolValue(recordInfo.Delete)

	diagnostics := response.State.Set(context, &newState)
	response.Diagnostics.Append(diagnostics...)
}

func (resource *recordMXResource) Delete(context context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	state, err := getDeleteState[recordMXResourceModel](context, request, response)
	if err != nil {
		return
	}

	domain := state.Domain.ValueString()
	id := state.ID.ValueString()

	err = resource.client.DeleteMXRecord(domain, id)
	if err != nil {
		response.Diagnostics.AddError("Error deleting MX record", "Request failed: "+err.Error())
	}
}

func validateMXPriority(diagnostics *diag.Diagnostics, priority int32) bool {
	if priority >= 0 && priority <= 65535 {
		return true
	}
	diagnostics.AddError(fmt.Sprintf("Invalid MX priority: %d", priority), "valid MX priority values range from 0 to 65535 (RFC 974)")
	return false
}
