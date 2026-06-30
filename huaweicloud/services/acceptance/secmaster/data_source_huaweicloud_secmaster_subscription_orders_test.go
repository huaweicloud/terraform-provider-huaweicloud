package secmaster

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSubscriptionOrders_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_subscription_orders.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSubscriptionOrders_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "csb_version"),
					resource.TestCheckResourceAttrSet(dataSource, "ecs_count"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.#"),
				),
			},
		},
	})
}

func testDataSourceSubscriptionOrders_basic() string {
	return `
data "huaweicloud_secmaster_subscription_orders" "test" {
  page = "PURCHASE"
}
`
}
