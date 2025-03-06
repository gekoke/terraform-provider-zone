package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func getCreatePlan(context context.Context, request resource.CreateRequest, response *resource.CreateResponse, target interface{}) bool {
	diagnostics := request.Plan.Get(context, target)
	response.Diagnostics.Append(diagnostics...)

	return !response.Diagnostics.HasError()
}

func getUpdatePlan(context context.Context, request resource.UpdateRequest, response *resource.UpdateResponse, target interface{}) bool {
	diagnostics := request.Plan.Get(context, target)
	response.Diagnostics.Append(diagnostics...)

	return !response.Diagnostics.HasError()
}
