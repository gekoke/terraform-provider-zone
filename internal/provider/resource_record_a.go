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
	_ resource.Resource                = &recordAResource{}
	_ resource.ResourceWithConfigure   = &recordAResource{}
	_ resource.ResourceWithImportState = &recordAResource{}
)

func NewRecordAResource() resource.Resource {
	return &recordAResource{}
}

type recordAResource struct {
	client api.Client
}

type recordAResourceModel struct {
	ID          types.String `tfsdk:"id"`
	LastUpdated types.String `tfsdk:"last_updated"`
	Domain      types.String `tfsdk:"domain"`
	Name        types.String `tfsdk:"name"`
	Destination types.String `tfsdk:"destination"`
	ResourceURL types.String `tfsdk:"resource_url"`
	Modify      types.Bool   `tfsdk:"modify"`
	Delete      types.Bool   `tfsdk:"delete"`
}

func (*recordAResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_record_a"
}

func (*recordAResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
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

func (resource *recordAResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (*recordAResource) ImportState(context context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(context, path.Root("id"), request, response)
}

func (resource *recordAResource) Create(context context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	plan, err := getCreatePlan[recordAResourceModel](context, request, response)
	if err != nil {
		return
	}

	destination, err := netip.ParseAddr(plan.Destination.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Error creating A record", "Error parsing address from destination: "+err.Error())
		return
	}
	if !destination.Is4() {
		response.Diagnostics.AddError("Error creating A record", "Destination must be an IPV4 address")
		return
	}
	domain := plan.Domain.ValueString()
	name := plan.Name.ValueString()
	record := api.ARecord{Name: name, Destination: destination}

	recordInfo, err := resource.client.CreateARecord(domain, record)
	if err != nil {
		response.Diagnostics.AddError("Error creating A record", "Request failed: "+err.Error())
		return
	}

	var newState recordAResourceModel
	newState.ID = types.StringValue(recordInfo.ID)
	newState.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	newState.Name = types.StringValue(recordInfo.Name)
	newState.Destination = types.StringValue(recordInfo.Destination.String())
	newState.ResourceURL = types.StringValue(recordInfo.ResourceURL.String())
	newState.Modify = types.BoolValue(recordInfo.Modify)
	newState.Delete = types.BoolValue(recordInfo.Delete)

	diagnostics := response.State.Set(context, &newState)
	response.Diagnostics.Append(diagnostics...)
}

func (resource *recordAResource) Read(context context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	state, err := getReadState[recordAResourceModel](context, request, response)
	if err != nil {
		return
	}

	domain := api.Domain(state.Domain.ValueString())
	id := api.Identificator(state.ID.ValueString())

	recordInfo, err := resource.client.GetARecord(domain, id)
	if err != nil {
		response.Diagnostics.AddError("Error reading A record", "Request failed: "+err.Error())
		return
	}

	var newState recordAResourceModel
	newState.ID = types.StringValue(recordInfo.ID)
	newState.Name = types.StringValue(recordInfo.Name)
	newState.Destination = types.StringValue(recordInfo.Destination.String())
	newState.ResourceURL = types.StringValue(recordInfo.ResourceURL.String())
	newState.Modify = types.BoolValue(recordInfo.Modify)
	newState.Delete = types.BoolValue(recordInfo.Delete)

	diagnostics := response.State.Set(context, &newState)
	response.Diagnostics.Append(diagnostics...)
}

func (resource *recordAResource) Update(context context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	plan, err := getUpdatePlan[recordAResourceModel](context, request, response)
	if err != nil {
		return
	}
	state, err := getUpdateState[recordAResourceModel](context, request, response)
	if err != nil {
		return
	}

	destination, err := netip.ParseAddr(plan.Destination.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Error creating A record", "Error parsing destination address: "+err.Error())
		return
	}
	if !destination.Is4() {
		response.Diagnostics.AddError("Error creating A record", "Destination address must be IPV4")
		return
	}
	domain := plan.Domain.ValueString()
	name := plan.Name.ValueString()
	recordID := state.ID.ValueString()
	record := api.ARecord{Name: name, Destination: destination}

	recordInfo, err := resource.client.UpdateARecord(domain, recordID, record)
	if err != nil {
		response.Diagnostics.AddError("Error updating A record", "Request failed: "+err.Error())
		return
	}

	//exhaustruct:enforce
	var newState = recordAAAAResourceModel{
		ID:          types.StringValue(recordInfo.ID),
		LastUpdated: types.StringValue(time.Now().Format(time.RFC850)),
		Domain:      types.StringValue(domain),
		Name:        types.StringValue(recordInfo.Name),
		Destination: types.StringValue(recordInfo.Destination.String()),
		ResourceURL: types.StringValue(recordInfo.ResourceURL.String()),
		Modify:      types.BoolValue(recordInfo.Modify),
		Delete:      types.BoolValue(recordInfo.Delete),
	}

	diagnostics := response.State.Set(context, &newState)
	response.Diagnostics.Append(diagnostics...)
}

func (resource *recordAResource) Delete(context context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	state, err := getDeleteState[recordAResourceModel](context, request, response)
	if err != nil {
		return
	}

	domain := state.Domain.ValueString()
	id := state.ID.ValueString()

	err = resource.client.DeleteARecord(domain, id)
	if err != nil {
		response.Diagnostics.AddError("Error deleting A record", "Request failed: "+err.Error())
	}
}
