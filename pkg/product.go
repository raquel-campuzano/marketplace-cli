// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/vmware-labs/marketplace-cli/v2/internal"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

type ListProductResponse struct {
	Response *ListProductResponsePayload `json:"response"`
}
type ListProductResponsePayload struct {
	Message    string            `json:"string"`
	StatusCode int               `json:"statuscode"`
	Products   []*models.Product `json:"dataList"`
	Params     struct {
		ProductCount int                  `json:"itemsnumber"`
		Pagination   *internal.Pagination `json:"pagination"`
	} `json:"params"`
}

func (m *Marketplace) ListProducts(allOrgs bool, searchTerm string) ([]*models.Product, error) {
	values := url.Values{
		"managed": []string{strconv.FormatBool(!allOrgs)},
	}
	if searchTerm != "" {
		values.Set("search", searchTerm)
	}

	var products []*models.Product
	totalProducts := 1
	pagination := &internal.Pagination{
		Page:     1,
		PageSize: 20,
	}

	for ; len(products) < totalProducts; pagination.Page++ {
		requestURL := m.MakeURL("/api/v1/products", values)
		requestURL = pagination.Apply(requestURL)
		resp, err := m.Get(requestURL)
		if err != nil {
			return nil, fmt.Errorf("sending the request for the list of products failed: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("getting the list of products failed: (%d) %s", resp.StatusCode, resp.Status)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read the list of products: %w", err)
		}

		response := &ListProductResponse{}
		err = json.Unmarshal(body, response)
		if err != nil {
			return nil, fmt.Errorf("failed to parse the list of products: %w", err)
		}
		totalProducts = response.Response.Params.ProductCount
		products = append(products, response.Response.Products...)
	}

	return products, nil
}

type GetProductResponse struct {
	Response *GetProductResponsePayload `json:"response"`
}
type GetProductResponsePayload struct {
	Message    string          `json:"message"`
	StatusCode int             `json:"statuscode"`
	Data       *models.Product `json:"data"`
}

func (m *Marketplace) GetProduct(slug string) (*models.Product, error) {
	requestURL := m.MakeURL(
		fmt.Sprintf("/api/v1/products/%s", slug),
		url.Values{
			"increaseViewCount": []string{"false"},
			"isSlug":            []string{"true"},
		},
	)

	resp, err := m.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("sending the request for product \"%s\" failed: %w", slug, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("product \"%s\" not found", slug)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("getting product \"%s\" failed: (%d)", slug, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response for product \"%s\": %w", slug, err)
	}

	response := &GetProductResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the response for product \"%s\": %w", slug, err)
	}
	return response.Response.Data, nil
}

func (m *Marketplace) GetProductWithVersion(slug, version string) (*models.Product, error) {
	product, err := m.GetProduct(slug)
	if err != nil {
		return nil, err
	}

	if version == "latest" {
		version = product.AllVersions[0].Number
	}
	if !product.HasVersion(version) {
		return nil, fmt.Errorf("product \"%s\" does not have a version %s", slug, version)
	}

	return product, nil
}

func (m *Marketplace) PutProduct(product *models.Product, versionUpdate bool) (*models.Product, error) {
	encoded, err := json.Marshal(product)
	if err != nil {
		return nil, err
	}

	requestURL := m.MakeURL(
		fmt.Sprintf("/api/v1/products/%s", product.ProductId),
		url.Values{
			"archivepreviousversion": []string{"false"},
			"isversionupdate":        []string{strconv.FormatBool(versionUpdate)},
		},
	)

	resp, err := m.Put(requestURL, bytes.NewReader(encoded), "application/json")
	if err != nil {
		return nil, fmt.Errorf("sending the update for product \"%s\" failed: %w", product.Slug, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the update response for product \"%s\": %w", product.Slug, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("updating product \"%s\" failed: (%d)\n%s", product.Slug, resp.StatusCode, body)
	}

	response := &GetProductResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the response for product \"%s\": %w", product.Slug, err)
	}
	return response.Response.Data, nil
}
