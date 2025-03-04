package api

import "resty.dev/v3"

type Client struct {
	resty resty.Client
}

func MakeClient(baseUrl string, username string, password string) Client {
	resty := resty.New()
	resty.SetBaseURL(baseUrl)
	resty.SetBasicAuth(username, password)
	return Client{
		resty: *resty,
	}
}
