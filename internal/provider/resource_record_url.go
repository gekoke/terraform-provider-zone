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
	_ resource.Resource                = &recordURLResource{}
	_ resource.ResourceWithConfigure   = &recordURLResource{}
	_ resource.ResourceWithImportState = &recordURLResource{}
)

func NewRecordURLResource() resource.Resource {
	return &recordURLResource{}
}

type recordURLResource struct {
	client api.Client
}

type recordURLResourceModel struct {
	ID          types.String `tfsdk:"id"`
	LastUpdated types.String `tfsdk:"last_updated"`
	Domain      types.String `tfsdk:"domain"`
	Name        types.String `tfsdk:"name"`
	Destination types.String `tfsdk:"destination"`
	Type        types.String `tfsdk:"type"`
	ResourceURL types.String `tfsdk:"resource_url"`
	Modify      types.Bool   `tfsdk:"modify"`
	Delete      types.Bool   `tfsdk:"delete"`
}

func (*recordURLResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_record_url"
}

func (*recordURLResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
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
			"type": schema.StringAttribute{
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

func (resource *recordURLResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (*recordURLResource) ImportState(context context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(context, path.Root("id"), request, response)
}

func (resource *recordURLResource) Create(context context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	plan, err := getCreatePlan[recordURLResourceModel](context, request, response)
	if err != nil {
		return
	}

	type_ := plan.Type.ValueString()
	if !validateType(&response.Diagnostics, type_) {
		return
	}
	domain := plan.Domain.ValueString()
	name := plan.Name.ValueString()
	destination := plan.Destination.ValueString()
	record := api.URLRecord{Name: name, Destination: destination, Type: type_}

	recordInfo, err := resource.client.CreateURLRecord(domain, record)
	if err != nil {
		response.Diagnostics.AddError("Error creating URL record", "Request failed: "+err.Error())
		return
	}

	var newState recordURLResourceModel
	newState.ID = types.StringValue(recordInfo.ID)
	newState.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	newState.Name = types.StringValue(recordInfo.Name)
	newState.Destination = types.StringValue(recordInfo.Destination)
	newState.Type = types.StringValue(recordInfo.Type)
	newState.ResourceURL = types.StringValue(recordInfo.ResourceURL.String())
	newState.Modify = types.BoolValue(recordInfo.Modify)
	newState.Delete = types.BoolValue(recordInfo.Delete)

	diagnostics := response.State.Set(context, &newState)
	response.Diagnostics.Append(diagnostics...)
}

func (resource *recordURLResource) Read(context context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	state, err := getReadState[recordURLResourceModel](context, request, response)
	if err != nil {
		return
	}

	domain := api.Domain(state.Domain.ValueString())
	id := api.Identificator(state.ID.ValueString())

	recordInfo, err := resource.client.GetURLRecord(domain, id)
	if err != nil {
		response.Diagnostics.AddError("Error reading URL record", "Request failed: "+err.Error())
		return
	}

	var newState recordURLResourceModel
	newState.ID = types.StringValue(recordInfo.ID)
	newState.Name = types.StringValue(recordInfo.Name)
	newState.Destination = types.StringValue(recordInfo.Destination)
	newState.Type = types.StringValue(recordInfo.Type)
	newState.ResourceURL = types.StringValue(recordInfo.ResourceURL.String())
	newState.Modify = types.BoolValue(recordInfo.Modify)
	newState.Delete = types.BoolValue(recordInfo.Delete)

	diagnostics := response.State.Set(context, &newState)
	response.Diagnostics.Append(diagnostics...)
}

func (resource *recordURLResource) Update(context context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	plan, err := getUpdatePlan[recordURLResourceModel](context, request, response)
	if err != nil {
		return
	}
	state, err := getUpdateState[recordURLResourceModel](context, request, response)
	if err != nil {
		return
	}

	type_ := plan.Type.ValueString()
	if !validateType(&response.Diagnostics, type_) {
		return
	}
	domain := plan.Domain.ValueString()
	name := plan.Name.ValueString()
	destination := plan.Destination.ValueString()
	recordID := state.ID.ValueString()
	record := api.URLRecord{Name: name, Destination: destination, Type: type_}

	recordInfo, err := resource.client.UpdateURLRecord(domain, recordID, record)
	if err != nil {
		response.Diagnostics.AddError("Error updating URL record", "Request failed: "+err.Error())
		return
	}

	var newState recordURLResourceModel
	newState.ID = types.StringValue(recordInfo.ID)
	newState.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	newState.Name = types.StringValue(recordInfo.Name)
	newState.Destination = types.StringValue(recordInfo.Destination)
	newState.Type = types.StringValue(recordInfo.Type)
	newState.ResourceURL = types.StringValue(recordInfo.ResourceURL.String())
	newState.Modify = types.BoolValue(recordInfo.Modify)
	newState.Delete = types.BoolValue(recordInfo.Delete)

	diagnostics := response.State.Set(context, &newState)
	response.Diagnostics.Append(diagnostics...)
}

func (resource *recordURLResource) Delete(context context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	state, err := getDeleteState[recordURLResourceModel](context, request, response)
	if err != nil {
		return
	}

	domain := state.Domain.ValueString()
	id := state.ID.ValueString()

	err = resource.client.DeleteURLRecord(domain, id)
	if err != nil {
		response.Diagnostics.AddError("Error deleting URL record", "Request failed: "+err.Error())
	}
}

func validateType(diagnostics *diag.Diagnostics, type_ string) bool {
	if type_ == "301" || type_ == "302" {
		return true
	}
	diagnostics.AddError(fmt.Sprintf("Invalid URL record type: %s", type_), "Valid URL record types are \"301\" and \"302\"")
	return false
}
