package provider

import (
	"context"
	"net/netip"
	"time"

	"github.com/gekoke/terraform-provider-zone/internal/api"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &recordAAAAResource{}
	_ resource.ResourceWithConfigure   = &recordAAAAResource{}
	_ resource.ResourceWithImportState = &recordAAAAResource{}
)

func NewRecordAAAAResource() resource.Resource {
	return &recordAAAAResource{}
}

type recordAAAAResource struct {
	client api.Client
}

type recordAAAAResourceModel struct {
	ID          types.String `tfsdk:"id"`
	LastUpdated types.String `tfsdk:"last_updated"`
	Domain      types.String `tfsdk:"domain"`
	Name        types.String `tfsdk:"name"`
	Destination types.String `tfsdk:"destination"`
	ResourceURL types.String `tfsdk:"resource_url"`
	Modify      types.Bool   `tfsdk:"modify"`
	Delete      types.Bool   `tfsdk:"delete"`
}

func (*recordAAAAResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_record_aaaa"
}

func (*recordAAAAResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
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

func (resource *recordAAAAResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (*recordAAAAResource) ImportState(context context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(context, path.Root("id"), request, response)
}

func (resource *recordAAAAResource) Create(context context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	plan, err := getCreatePlan[recordAAAAResourceModel](context, request, response)
	if err != nil {
		return
	}

	destination, err := netip.ParseAddr(plan.Destination.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Error creating AAAA record", "Error parsing address from destination: "+err.Error())
		return
	}
	if !destination.Is6() {
		response.Diagnostics.AddError("Error creating AAAA record", "Destination must be an IPV6 address")
		return
	}
	domain := plan.Domain.ValueString()
	name := plan.Name.ValueString()
	record := api.AAAARecord{Name: name, Destination: destination}

	recordInfo, err := resource.client.CreateAAAARecord(domain, record)
	if err != nil {
		response.Diagnostics.AddError("Error creating AAAA record", "Request failed: "+err.Error())
		return
	}

	var newState recordAAAAResourceModel
	newState.ID = types.StringValue(recordInfo.ID)
	newState.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	newState.Name = types.StringValue(recordInfo.Name)
	newState.Destination = types.StringValue(recordInfo.Destination.StringExpanded())
	newState.ResourceURL = types.StringValue(recordInfo.ResourceURL.String())
	newState.Modify = types.BoolValue(recordInfo.Modify)
	newState.Delete = types.BoolValue(recordInfo.Delete)

	diagnostics := response.State.Set(context, &newState)
	response.Diagnostics.Append(diagnostics...)
}

func (resource *recordAAAAResource) Read(context context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	state, err := getReadState[recordAAAAResourceModel](context, request, response)
	if err != nil {
		return
	}

	domain := api.Domain(state.Domain.ValueString())
	id := api.Identificator(state.ID.ValueString())

	recordInfo, err := resource.client.GetAAAARecord(domain, id)
	if err != nil {
		response.Diagnostics.AddError("Error reading AAAA record", "Request failed: "+err.Error())
		return
	}

	var newState recordAAAAResourceModel
	newState.ID = types.StringValue(recordInfo.ID)
	newState.Name = types.StringValue(recordInfo.Name)
	newState.Destination = types.StringValue(recordInfo.Destination.String())
	newState.ResourceURL = types.StringValue(recordInfo.ResourceURL.String())
	newState.Modify = types.BoolValue(recordInfo.Modify)
	newState.Delete = types.BoolValue(recordInfo.Delete)

	diagnostics := response.State.Set(context, &newState)
	response.Diagnostics.Append(diagnostics...)
}

func (resource *recordAAAAResource) Update(context context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	plan, err := getUpdatePlan[recordAAAAResourceModel](context, request, response)
	if err != nil {
		return
	}
	state, err := getUpdateState[recordAAAAResourceModel](context, request, response)
	if err != nil {
		return
	}

	destination, err := netip.ParseAddr(plan.Destination.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Error creating AAAA record", "Error parsing destination address: "+err.Error())
		return
	}
	if !destination.Is6() {
		response.Diagnostics.AddError("Error creating AAAA record", "Destination address must be IPV6")
		return
	}
	domain := plan.Domain.ValueString()
	name := plan.Name.ValueString()
	recordID := state.ID.ValueString()
	record := api.AAAARecord{Name: name, Destination: destination}

	recordInfo, err := resource.client.UpdateAAAARecord(domain, recordID, record)
	if err != nil {
		response.Diagnostics.AddError("Error updating AAAA record", "Request failed: "+err.Error())
		return
	}

	var newState recordAAAAResourceModel
	newState.ID = types.StringValue(recordInfo.ID)
	newState.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	newState.Name = types.StringValue(recordInfo.Name)
	newState.Destination = types.StringValue(recordInfo.Destination.StringExpanded())
	newState.ResourceURL = types.StringValue(recordInfo.ResourceURL.String())
	newState.Modify = types.BoolValue(recordInfo.Modify)
	newState.Delete = types.BoolValue(recordInfo.Delete)

	diagnostics := response.State.Set(context, &newState)
	response.Diagnostics.Append(diagnostics...)
}

func (resource *recordAAAAResource) Delete(context context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	state, err := getDeleteState[recordAAAAResourceModel](context, request, response)
	if err != nil {
		return
	}

	domain := state.Domain.ValueString()
	id := state.ID.ValueString()

	err = resource.client.DeleteAAAARecord(domain, id)
	if err != nil {
		response.Diagnostics.AddError("Error deleting AAAA record", "Request failed: "+err.Error())
	}
}
