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
	var plan recordAResourceModel

	diagnostics := request.Plan.Get(context, &plan)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
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

	plan.ID = types.StringValue(recordInfo.ID)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.Name = types.StringValue(recordInfo.Name)
	plan.Destination = types.StringValue(recordInfo.Destination.String())
	plan.ResourceURL = types.StringValue(recordInfo.ResourceURL.String())
	plan.Modify = types.BoolValue(recordInfo.Modify)
	plan.Delete = types.BoolValue(recordInfo.Delete)

	diagnostics = response.State.Set(context, plan)
	response.Diagnostics.Append(diagnostics...)
}

func (resource *recordAResource) Read(context context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state recordAResourceModel
	diagnostics := request.State.Get(context, &state)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	domain := api.Domain(state.Domain.ValueString())
	id := api.Identificator(state.ID.ValueString())

	recordInfo, err := resource.client.GetARecord(domain, id)

	if err != nil {
		response.Diagnostics.AddError("Error reading A record", "Request failed: "+err.Error())
		return
	}

	state.ID = types.StringValue(recordInfo.ID)
	state.Name = types.StringValue(recordInfo.Name)
	state.Destination = types.StringValue(recordInfo.Destination.String())
	state.ResourceURL = types.StringValue(recordInfo.ResourceURL.String())
	state.Modify = types.BoolValue(recordInfo.Modify)
	state.Delete = types.BoolValue(recordInfo.Delete)

	diagnostics = response.State.Set(context, &state)
	response.Diagnostics.Append(diagnostics...)
}

func (resource *recordAResource) Update(context context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan recordAResourceModel
	var state recordAResourceModel

	diagnostics := request.Plan.Get(context, &plan)
	response.Diagnostics.Append(diagnostics...)
	diagnostics = request.State.Get(context, &state)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
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

	plan.ID = types.StringValue(recordInfo.ID)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.Name = types.StringValue(recordInfo.Name)
	plan.Destination = types.StringValue(recordInfo.Destination.String())
	plan.ResourceURL = types.StringValue(recordInfo.ResourceURL.String())
	plan.Modify = types.BoolValue(recordInfo.Modify)
	plan.Delete = types.BoolValue(recordInfo.Delete)

	diagnostics = response.State.Set(context, plan)
	response.Diagnostics.Append(diagnostics...)
}

func (resource *recordAResource) Delete(context context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state recordAResourceModel
	diagnostics := request.State.Get(context, &state)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	domain := state.Domain.ValueString()
	id := state.ID.ValueString()

	err := resource.client.DeleteARecord(domain, id)

	if err != nil {
		response.Diagnostics.AddError("Error deleting A record", "Request failed: "+err.Error())
	}
}
