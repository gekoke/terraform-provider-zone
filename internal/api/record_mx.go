package api

import "net/url"

type MXRecordInfo struct {
	ID          Identificator `json:"id"`
	Name        string        `json:"name"`
	Destination string        `json:"destination"`
	Priority    uint16        `json:"priority"`
	ResourceURL url.URL       `json:"resource_url"`
	Delete      bool          `json:"delete"`
	Modify      bool          `json:"modify"`
}

type MXRecord struct {
	Name        string `json:"name"`
	Destination string `json:"destination"`
	Priority    uint16 `json:"priority"`
}

func (c Client) CreateMXRecord(domain Domain, mxRecord MXRecord) (MXRecordInfo, error) {
	var result []MXRecordInfo

	res, _ := c.resty.R().
		SetResult(&result).
		SetContentType("application/json").
		SetBody(mxRecord).
		SetPathParam("domain", domain).
		Post("/dns/{domain}/mx")

	if err := getResponseError(res); err != nil {
		return MXRecordInfo{}, err
	}
	record := result[0]
	return record, nil
}

func (c Client) GetMXRecords(domain Domain) ([]MXRecordInfo, error) {
	var result []MXRecordInfo

	res, _ := c.resty.R().
		SetResult(&result).
		SetPathParam("domain", domain).
		Get("/dns/{domain}/mx")

	if err := getResponseError(res); err != nil {
		return nil, err
	}
	return result, nil
}

func (c Client) GetMXRecord(domain Domain, id Identificator) (MXRecordInfo, error) {
	var result []MXRecordInfo

	res, _ := c.resty.R().
		SetResult(&result).
		SetPathParam("domain", domain).
		SetPathParam("id", id).
		Get("/dns/{domain}/mx/{id}")

	if err := getResponseError(res); err != nil {
		return MXRecordInfo{}, err
	}
	record := result[0]
	return record, nil
}

func (c Client) UpdateMXRecord(domain Domain, id Identificator, mxRecord MXRecord) (MXRecordInfo, error) {
	var result []MXRecordInfo

	res, _ := c.resty.R().
		SetResult(&result).
		SetContentType("application/json").
		SetBody(mxRecord).
		SetPathParam("domain", domain).
		SetPathParam("id", id).
		Put("/dns/{domain}/mx/{id}")

	if err := getResponseError(res); err != nil {
		return MXRecordInfo{}, err
	}
	record := result[0]
	return record, nil
}

func (c Client) DeleteMXRecord(domain Domain, id Identificator) error {
	res, _ := c.resty.R().
		SetPathParam("domain", domain).
		SetPathParam("id", id).
		Delete("/dns/{domain}/mx/{id}")

	return getResponseError(res)
}
