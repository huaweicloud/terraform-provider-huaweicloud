package rds

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsMarketplaceEngineProducts_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_marketplace_engine_products.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsMarketplaceEngineProducts_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "marketplace_engine_products.#"),
					resource.TestCheckResourceAttrSet(dataSource, "marketplace_engine_products.0.engine_id"),
					resource.TestCheckResourceAttrSet(dataSource, "marketplace_engine_products.0.engine_version"),
					resource.TestCheckResourceAttrSet(dataSource, "marketplace_engine_products.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "marketplace_engine_products.0.product_id"),
					resource.TestCheckResourceAttrSet(dataSource, "marketplace_engine_products.0.bp_name"),
					resource.TestCheckResourceAttrSet(dataSource, "marketplace_engine_products.0.bp_domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "marketplace_engine_products.0.instance_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "marketplace_engine_products.0.image_id"),
					resource.TestCheckResourceAttrSet(dataSource, "marketplace_engine_products.0.user_license_agreement"),
					resource.TestCheckResourceAttrSet(dataSource, "marketplace_engine_products.0.agreements.#"),
					resource.TestCheckOutput("engine_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceRdsMarketplaceEngineProducts_basic() string {
	return `
data "huaweicloud_rds_business_partners" "test" {}

data "huaweicloud_rds_marketplace_engine_products" "test" {
  bp_domain_id = data.huaweicloud_rds_business_partners.test.business_partners[0].bp_domain_id
}

data "huaweicloud_rds_marketplace_engine_products" "engine_id_filter" {
  bp_domain_id = data.huaweicloud_rds_business_partners.test.business_partners[0].bp_domain_id
  engine_id    = data.huaweicloud_rds_marketplace_engine_products.test.marketplace_engine_products.0.engine_id
}

locals {
  engine_id = data.huaweicloud_rds_marketplace_engine_products.test.marketplace_engine_products.0.engine_id
}

output "engine_id_filter_is_useful" {
  value = length(data.huaweicloud_rds_marketplace_engine_products.engine_id_filter.marketplace_engine_products) > 0 && alltrue(
  [for v in data.huaweicloud_rds_marketplace_engine_products.engine_id_filter.marketplace_engine_products[*] :
  v.engine_id == local.engine_id]
  )
}
`
}
