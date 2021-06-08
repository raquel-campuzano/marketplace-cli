// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
)

type PaginationObject struct {
	Page     int32 `json:"page"`
	PageSize int32 `json:"pagesize"`
}

func (p PaginationObject) ToUrlValue() []string {
	data, _ := json.Marshal(p)
	return []string{string(data)}
}

func Pagination(page, pageSize int32) []string {
	return PaginationObject{
		Page:     page,
		PageSize: pageSize,
	}.ToUrlValue()
}

func MakeRequest(method, path string, params url.Values, header map[string]string, content io.Reader) (*http.Request, error) {
	marketplaceURL := &url.URL{
		Scheme:   "https",
		Host:     viper.GetString("marketplace.host"),
		Path:     path,
		RawQuery: params.Encode(),
	}

	req, err := http.NewRequest(method, marketplaceURL.String(), content)
	if err != nil {
		return nil, fmt.Errorf("failed to build %s request: %w", path, err)
	}

	if header != nil {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("csp-auth-token", viper.GetString("csp.refresh-token"))
	return req, nil
}

func MakeGetRequest(path string, params url.Values) (*http.Request, error) {
	return MakeRequest("GET", path, params, map[string]string{}, nil)
}

//go:generate counterfeiter . HTTPClient
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var Client HTTPClient

func init() {
	Client = &http.Client{}
}
