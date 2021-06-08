// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/vmware-labs/marketplace-cli/v2/cmd"
	"github.com/vmware-labs/marketplace-cli/v2/lib"
	"github.com/vmware-labs/marketplace-cli/v2/lib/libfakes"
)

var _ = Describe("OVA", func() {
	var (
		stdout *Buffer
		stderr *Buffer

		originalHttpClient lib.HTTPClient
		httpClient         *libfakes.FakeHTTPClient
	)

	BeforeEach(func() {
		stdout = NewBuffer()
		stderr = NewBuffer()

		originalHttpClient = lib.Client
		httpClient = &libfakes.FakeHTTPClient{}
		lib.Client = httpClient
	})

	AfterEach(func() {
		lib.Client = originalHttpClient
	})

	Describe("ListOVACmd", func() {
		BeforeEach(func() {
			product := CreateFakeProduct(
				"",
				"My Super Product",
				"my-super-product",
				"PENDING")
			AddVerions(product, "1.2.3", "2.3.4")
			product.ProductDeploymentFiles = append(product.ProductDeploymentFiles, CreateFakeOVA("fake-ova", "1.2.3"))
			response := &cmd.GetProductResponse{
				Response: &cmd.GetProductResponsePayload{
					Data:       product,
					StatusCode: http.StatusOK,
					Message:    "testing",
				},
			}

			responseBytes, err := json.Marshal(response)
			Expect(err).ToNot(HaveOccurred())

			httpClient.DoReturns(&http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(responseBytes)),
			}, nil)

			cmd.ListOVACmd.SetOut(stdout)
			cmd.ListOVACmd.SetErr(stderr)
		})

		It("outputs the ovas", func() {
			cmd.ProductSlug = "my-super-product"
			cmd.ProductVersion = "1.2.3"
			err := cmd.ListOVACmd.RunE(cmd.ListOVACmd, []string{""})
			Expect(err).ToNot(HaveOccurred())

			By("sending the correct request", func() {
				Expect(httpClient.DoCallCount()).To(Equal(1))
				request := httpClient.DoArgsForCall(0)
				Expect(request.Method).To(Equal("GET"))
				Expect(request.URL.Path).To(Equal("/api/v1/products/my-super-product"))
				Expect(request.URL.Query().Get("increaseViewCount")).To(Equal("false"))
				Expect(request.URL.Query().Get("isSlug")).To(Equal("true"))
			})

			By("outputting the response", func() {
				Expect(stdout).To(Say("NAME      SIZE     TYPE      FILES"))
				Expect(stdout).To(Say("fake-ova  1000100  fake.ovf  2"))
			})
		})

		Context("No product found", func() {
			BeforeEach(func() {
				httpClient.DoReturns(&http.Response{
					StatusCode: http.StatusNotFound,
				}, nil)
			})

			It("says there are no products", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "1.2.3"
				err := cmd.ListOVACmd.RunE(cmd.ListOVACmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" not found"))
			})
		})

		Context("Error fetching products", func() {
			BeforeEach(func() {
				httpClient.DoReturnsOnCall(0, nil, fmt.Errorf("request failed"))
			})

			It("prints the error", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "1.2.3"
				err := cmd.ListOVACmd.RunE(cmd.ListOVACmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("sending the request for product \"my-super-product\" failed: request failed"))
			})
		})

		Context("No product version found", func() {
			It("says that the version does not exist", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "9.9.9"
				err := cmd.ListOVACmd.RunE(cmd.ListOVACmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" does not have a version 9.9.9"))
			})
		})

		Context("No ovas", func() {
			It("says there are no container images", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "2.3.4"
				err := cmd.ListOVACmd.RunE(cmd.ListOVACmd, []string{""})
				Expect(err).ToNot(HaveOccurred())
				Expect(stdout).To(Say("product \"my-super-product\" 2.3.4 does not have any OVAs"))
			})
		})
	})
})
