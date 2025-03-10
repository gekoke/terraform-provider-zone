package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func getCreatePlan[T any](context context.Context, request resource.CreateRequest, response *resource.CreateResponse) (T, error) {
	var target T
	diagnostics := request.Plan.Get(context, &target)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return target, errors.New("Error encountered when getting plan")
	}
	return target, nil
}

func getReadState[T any](context context.Context, request resource.ReadRequest, response *resource.ReadResponse) (T, error) {
	var target T
	diagnostics := request.State.Get(context, &target)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return target, errors.New("Error encountered when getting state")
	}
	return target, nil
}

func getUpdatePlan[T any](context context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) (T, error) {
	var target T
	diagnostics := request.Plan.Get(context, &target)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return target, errors.New("Error encountered when getting plan")
	}
	return target, nil
}

func getUpdateState[T any](context context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) (T, error) {
	var target T
	diagnostics := request.State.Get(context, &target)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return target, errors.New("Error encountered when getting state")
	}
	return target, nil
}

func getDeleteState[T any](context context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) (T, error) {
	var target T
	diagnostics := request.State.Get(context, &target)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return target, errors.New("Error encountered when getting state")
	}
	return target, nil
}
