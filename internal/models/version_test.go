// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package models_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("GetVersion", func() {
	It("gets the version object from the version number", func() {
		product := test.CreateFakeProduct("", "My Product", "my-product", "PENDING")
		test.AddVerions(product, "1.2.3")

		version := product.GetVersion("1.2.3")
		Expect(version).ToNot(BeNil())
		Expect(version.Number).To(Equal("1.2.3"))
	})

	Context("version does not exist", func() {
		It("returns nil", func() {
			product := test.CreateFakeProduct("", "My Product", "my-product", "PENDING")
			test.AddVerions(product, "1.2.3")

			version := product.GetVersion("9.9.9")
			Expect(version).To(BeNil())
		})
	})

	Context("Argument is empty", func() {
		It("returns the latest version", func() {
			product := test.CreateFakeProduct("", "My Product", "my-product", "PENDING")
			test.AddVerions(product, "1.2.3", "2.3.4", "0.0.1")

			version := product.GetVersion("")
			Expect(version).ToNot(BeNil())
			Expect(version.Number).To(Equal("2.3.4"))
		})

		Context("The product has no versions", func() {
			It("returns nil", func() {
				product := test.CreateFakeProduct("", "My Product", "my-product", "PENDING")

				version := product.GetVersion("")
				Expect(version).To(BeNil())
			})
		})
	})
})