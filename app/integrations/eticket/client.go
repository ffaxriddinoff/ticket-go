package eticket

import (
	"net/http"
	"strings"
	"ticket-backend/app/facades"
	"time"
)

type Client struct {
	baseUrl    string
	httpClient *http.Client
}

func NewClient() *Client {
	config := facades.Config()

	baseUrl := config.GetString("eticket.base_url")

	timeout := config.GetInt("eticket.timeout", 10)

	return &Client{
		baseUrl: strings.TrimRight(baseUrl, "/"),
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}
