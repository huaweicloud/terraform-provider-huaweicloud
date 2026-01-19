package secmaster

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSubscriptionProducts_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_subscription_products.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSubscriptionProducts_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "basic.#"),
					resource.TestCheckResourceAttrSet(dataSource, "basic.0.cloud_service_type"),
					resource.TestCheckResourceAttrSet(dataSource, "basic.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "basic.0.resource_spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "basic.0.resource_size_measure_id"),
					resource.TestCheckResourceAttrSet(dataSource, "basic.0.usage_factor"),
					resource.TestCheckResourceAttrSet(dataSource, "basic.0.usage_measure_id"),
					resource.TestCheckResourceAttrSet(dataSource, "basic.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "standard.#"),
					resource.TestCheckResourceAttrSet(dataSource, "standard.0.cloud_service_type"),
					resource.TestCheckResourceAttrSet(dataSource, "standard.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "standard.0.resource_spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "standard.0.resource_size_measure_id"),
					resource.TestCheckResourceAttrSet(dataSource, "standard.0.usage_factor"),
					resource.TestCheckResourceAttrSet(dataSource, "standard.0.usage_measure_id"),
					resource.TestCheckResourceAttrSet(dataSource, "standard.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "professional.#"),
					resource.TestCheckResourceAttrSet(dataSource, "large_screen.#"),
					resource.TestCheckResourceAttrSet(dataSource, "log_collection.#"),
					resource.TestCheckResourceAttrSet(dataSource, "log_retention.#"),
					resource.TestCheckResourceAttrSet(dataSource, "log_analysis.#"),
					resource.TestCheckResourceAttrSet(dataSource, "soar.#"),
				),
			},
		},
	})
}

const testAccDataSourceSubscriptionProducts_basic = `
data "huaweicloud_secmaster_subscription_products" "test" {
}
`
