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
	var plan recordMXResourceModel

	diagnostics := request.Plan.Get(context, &plan)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
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

	plan.ID = types.StringValue(recordInfo.ID)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.Name = types.StringValue(recordInfo.Name)
	plan.Destination = types.StringValue(recordInfo.Destination)
	plan.Priority = types.Int32Value(int32(recordInfo.Priority))
	plan.ResourceURL = types.StringValue(recordInfo.ResourceURL.String())
	plan.Modify = types.BoolValue(recordInfo.Modify)
	plan.Delete = types.BoolValue(recordInfo.Delete)

	diagnostics = response.State.Set(context, plan)
	response.Diagnostics.Append(diagnostics...)
}

func (resource *recordMXResource) Read(context context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state recordMXResourceModel
	diagnostics := request.State.Get(context, &state)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	domain := api.Domain(state.Domain.ValueString())
	id := api.Identificator(state.ID.ValueString())

	mxRecord, err := resource.client.GetMXRecord(domain, id)

	if err != nil {
		response.Diagnostics.AddError("Error reading MX record", "Request failed: "+err.Error())
		return
	}

	state.ID = types.StringValue(mxRecord.ID)
	state.Name = types.StringValue(mxRecord.Name)
	state.Destination = types.StringValue(mxRecord.Destination)
	state.Priority = types.Int32Value(int32(mxRecord.Priority))
	state.ResourceURL = types.StringValue(mxRecord.ResourceURL.String())
	state.Modify = types.BoolValue(mxRecord.Modify)
	state.Delete = types.BoolValue(mxRecord.Delete)

	diagnostics = response.State.Set(context, &state)
	response.Diagnostics.Append(diagnostics...)
}

func (resource *recordMXResource) Update(context context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan recordMXResourceModel
	var state recordMXResourceModel

	diagnostics := request.Plan.Get(context, &plan)
	response.Diagnostics.Append(diagnostics...)
	diagnostics = request.State.Get(context, &state)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
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

	plan.ID = types.StringValue(recordInfo.ID)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.Name = types.StringValue(recordInfo.Name)
	plan.Destination = types.StringValue(recordInfo.Destination)
	plan.Priority = types.Int32Value(int32(record.Priority))
	plan.ResourceURL = types.StringValue(recordInfo.ResourceURL.String())
	plan.Modify = types.BoolValue(recordInfo.Modify)
	plan.Delete = types.BoolValue(recordInfo.Delete)

	diagnostics = response.State.Set(context, plan)
	response.Diagnostics.Append(diagnostics...)
}

func (resource *recordMXResource) Delete(context context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state recordMXResourceModel
	diagnostics := request.State.Get(context, &state)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	domain := state.Domain.ValueString()
	id := state.ID.ValueString()

	err := resource.client.DeleteMXRecord(domain, id)

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
