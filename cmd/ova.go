// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vmware-labs/marketplace-cli/v2/internal"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

var pvaFile string

func init() {
	rootCmd.AddCommand(OVACmd)
	OVACmd.AddCommand(ListOVACmd)
	OVACmd.AddCommand(CreateOVACmd)
	OVACmd.PersistentFlags().StringVarP(&OutputFormat, "output-format", "f", FormatTable, "Output format")

	OVACmd.PersistentFlags().StringVarP(&ProductSlug, "product", "p", "", "Product slug")
	_ = OVACmd.MarkPersistentFlagRequired("product")
	OVACmd.PersistentFlags().StringVarP(&ProductVersion, "product-version", "v", "", "Product version")
	_ = OVACmd.MarkPersistentFlagRequired("product-version")

	CreateOVACmd.Flags().StringVar(&pvaFile, "ova-file", "", "OVA file to upload")
}

var OVACmd = &cobra.Command{
	Use:               "ova",
	Aliases:           []string{"ovas"},
	Short:             "ova",
	Long:              "",
	Args:              cobra.OnlyValidArgs,
	ValidArgs:         []string{"get", "list", "create"},
	PersistentPreRunE: GetRefreshToken,
}

var ListOVACmd = &cobra.Command{
	Use:   "list",
	Short: "list OVAs",
	Long:  "",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		product, err := Marketplace.GetProduct(ProductSlug)
		if err != nil {
			return err
		}

		if !product.HasVersion(ProductVersion) {
			return fmt.Errorf("product \"%s\" does not have a version %s", ProductSlug, ProductVersion)
		}

		return RenderOVAs(OutputFormat, ProductVersion, product, cmd.OutOrStdout())
	},
}

var CreateOVACmd = &cobra.Command{
	Use:     "create",
	Short:   "add an OVA to a product",
	Long:    "",
	Args:    cobra.NoArgs,
	PreRunE: GetUploadCredentials,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		product, err := Marketplace.GetProduct(ProductSlug)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		if !product.HasVersion(ProductVersion) {
			cmd.SilenceUsage = true
			return fmt.Errorf("product \"%s\" does not have a version %s, please add it first", ProductSlug, ProductVersion)
		}

		hashAlgo := internal.HashAlgoSHA1
		uploader := internal.NewS3Uploader(Marketplace.StorageRegion, hashAlgo, product.PublisherDetails.OrgId, UploadCredentials)
		fileURL, fileHash, err := uploader.Upload(Marketplace.StorageBucket, pvaFile)
		if err != nil {
			return err
		}

		product.ProductDeploymentFiles = []*models.ProductDeploymentFile{{
			Url:        fileURL,
			AppVersion: ProductVersion,
			HashDigest: fileHash,
			HashAlgo:   models.HashAlgoSHA1,
		}}

		_, err = Marketplace.PutProduct(product, false)
		if err != nil {
			return err
		}

		return nil // RenderOVAs(OutputFormat, ProductVersion, putResponse.Response.Data, cmd.OutOrStdout())
	},
}
