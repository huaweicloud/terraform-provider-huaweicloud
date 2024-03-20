package iotda

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceProducts_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_iotda_products.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProducts_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "products.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "products.0.space_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "products.0.space_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "products.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "products.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "products.0.device_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "products.0.data_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "products.0.manufacturer_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "products.0.industry"),
					resource.TestCheckResourceAttrSet(dataSourceName, "products.0.created_at"),

					resource.TestCheckOutput("is_product_id_filter_useful", "true"),
					resource.TestCheckOutput("is_product_name_filter_useful", "true"),
					resource.TestCheckOutput("is_space_id_filter_useful", "true"),
					resource.TestCheckOutput("is_space_name_filter_useful", "true"),
					resource.TestCheckOutput("is_device_type_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func TestAccDataSourceProducts_derived(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_iotda_products.test_derived"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProducts_derived(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "products.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "products.0.space_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "products.0.space_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "products.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "products.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "products.0.device_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "products.0.data_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "products.0.manufacturer_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "products.0.industry"),
					resource.TestCheckResourceAttrSet(dataSourceName, "products.0.created_at"),

					resource.TestCheckOutput("is_product_id_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccDataSourceProducts_basic() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%s

data "huaweicloud_iotda_products" "test" {
  depends_on = [huaweicloud_iotda_product.test]
}

# Filter using product ID.
locals {
  product_id = data.huaweicloud_iotda_products.test.products[0].id
}

data "huaweicloud_iotda_products" "product_id_filter" {
  product_id = local.product_id
}

output "is_product_id_filter_useful" {
  value = length(data.huaweicloud_iotda_products.product_id_filter.products) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_products.product_id_filter.products[*].id : v == local.product_id]
  )
}

# Filter using product name.
locals {
  product_name = data.huaweicloud_iotda_products.test.products[0].name
}

data "huaweicloud_iotda_products" "product_name_filter" {
  product_name = local.product_name
}

output "is_product_name_filter_useful" {
  value = length(data.huaweicloud_iotda_products.product_name_filter.products) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_products.product_name_filter.products[*].name : v == local.product_name]
  )
}

# Filter using space ID.
locals {
  space_id = data.huaweicloud_iotda_products.test.products[0].space_id
}

data "huaweicloud_iotda_products" "space_id_filter" {
  space_id = local.space_id
}

output "is_space_id_filter_useful" {
  value = length(data.huaweicloud_iotda_products.space_id_filter.products) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_products.space_id_filter.products[*].space_id : v == local.space_id]
  )
}

# Filter using space name.
locals {
  space_name = data.huaweicloud_iotda_products.test.products[0].space_name
}

data "huaweicloud_iotda_products" "space_name_filter" {
  space_name = local.space_name
}

output "is_space_name_filter_useful" {
  value = length(data.huaweicloud_iotda_products.space_name_filter.products) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_products.space_name_filter.products[*].space_name : v == local.space_name]
  )
}

# Filter using device type.
locals {
  device_type = data.huaweicloud_iotda_products.test.products[0].device_type
}

data "huaweicloud_iotda_products" "device_type_filter" {
  device_type = local.device_type
}

output "is_device_type_filter_useful" {
  value = length(data.huaweicloud_iotda_products.device_type_filter.products) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_products.device_type_filter.products[*].device_type : v == local.device_type]
  )
}

# Filter using non existent product name.
data "huaweicloud_iotda_products" "not_found" {
  product_name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_iotda_products.not_found.products) == 0
}
`, testProduct_basic(name))
}

func testAccDataSourceProducts_derived() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%[1]s
%[2]s

data "huaweicloud_iotda_products" "test_derived" {
  depends_on = [huaweicloud_iotda_product.test]
}

# Filter using product ID.
locals {
  product_id = data.huaweicloud_iotda_products.test_derived.products[0].id
}

data "huaweicloud_iotda_products" "product_id_filter" {
  product_id = local.product_id
}

output "is_product_id_filter_useful" {
  value = length(data.huaweicloud_iotda_products.product_id_filter.products) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_products.product_id_filter.products[*].id : v == local.product_id]
  )
}

# Filter using non existent product name.
data "huaweicloud_iotda_products" "not_found" {
  product_name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_iotda_products.not_found.products) == 0
}
`, testProduct_basic(name), buildIoTDAEndpoint())
}
