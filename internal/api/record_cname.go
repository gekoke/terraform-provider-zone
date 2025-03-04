package api

import "net/url"

type CNAMERecordInfo struct {
	ID          Identificator `json:"id"`
	Name        string        `json:"name"`
	Destination string        `json:"destination"`
	ResourceURL url.URL       `json:"resource_url"`
	Delete      bool          `json:"delete"`
	Modify      bool          `json:"modify"`
}

type CNAMERecord struct {
	Name        string `json:"name"`
	Destination string `json:"destination"`
}

func (c Client) CreateCNAMERecord(domain Domain, cnameRecord CNAMERecord) (CNAMERecordInfo, error) {
	var result []CNAMERecordInfo

	res, _ := c.resty.R().
		SetResult(&result).
		SetContentType("application/json").
		SetBody(cnameRecord).
		SetPathParam("domain", domain).
		Post("/dns/{domain}/cname")

	if err := getResponseError(res); err != nil {
		return CNAMERecordInfo{}, err
	}
	record := result[0]
	return record, nil
}

func (c Client) GetCNAMERecords(domain Domain) ([]CNAMERecordInfo, error) {
	var result []CNAMERecordInfo

	res, _ := c.resty.R().
		SetResult(&result).
		SetPathParam("domain", domain).
		Get("/dns/{domain}/cname")

	if err := getResponseError(res); err != nil {
		return nil, err
	}
	return result, nil
}

func (c Client) GetCNAMERecord(domain Domain, id Identificator) (CNAMERecordInfo, error) {
	var result []CNAMERecordInfo

	res, _ := c.resty.R().
		SetResult(&result).
		SetPathParam("domain", domain).
		SetPathParam("id", id).
		Get("/dns/{domain}/cname/{id}")

	if err := getResponseError(res); err != nil {
		return CNAMERecordInfo{}, err
	}
	record := result[0]
	return record, nil
}

func (c Client) UpdateCNAMERecord(domain Domain, id Identificator, cnameRecord CNAMERecord) (CNAMERecordInfo, error) {
	var result []CNAMERecordInfo

	res, _ := c.resty.R().
		SetResult(&result).
		SetContentType("application/json").
		SetBody(cnameRecord).
		SetPathParam("domain", domain).
		SetPathParam("id", id).
		Put("/dns/{domain}/cname/{id}")

	if err := getResponseError(res); err != nil {
		return CNAMERecordInfo{}, err
	}
	record := result[0]
	return record, nil
}

func (c Client) DeleteCNAMERecord(domain Domain, id Identificator) error {
	res, _ := c.resty.R().
		SetPathParam("domain", domain).
		SetPathParam("id", id).
		Delete("/dns/{domain}/cname/{id}")

	return getResponseError(res)
}
