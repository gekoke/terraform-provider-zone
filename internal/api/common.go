package api

import (
	"errors"
	"fmt"
	"io"

	"resty.dev/v3"
)

type Identificator = string

type ServiceName = string

type Domain = ServiceName

func getResponseError(res *resty.Response) error {
	if res.IsSuccess() {
		return nil
	}
	body, _ := io.ReadAll(res.Body)
	return errors.New(fmt.Sprintf("status: %s, body: %s", res.Status(), string(body)))
}
