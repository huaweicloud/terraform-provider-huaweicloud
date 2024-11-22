package smn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSmnSubscriptions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_smn_subscriptions.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSmnSubscriptions_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.0.protocol"),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.0.endpoint"),
					resource.TestCheckOutput("protocol_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("endpoint_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSmnSubscriptions_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_smn_subscriptions" "test" {
  depends_on = [ 
    huaweicloud_smn_subscription.subscription_1,
    huaweicloud_smn_subscription.subscription_2,
    huaweicloud_smn_subscription.subscription_3,
    huaweicloud_smn_subscription.subscription_4,
  ]
}

locals {
  protocol = data.huaweicloud_smn_subscriptions.test.subscriptions[0].protocol
  status   = tostring(data.huaweicloud_smn_subscriptions.test.subscriptions[0].status)
  endpoint = data.huaweicloud_smn_subscriptions.test.subscriptions[0].endpoint
}

data "huaweicloud_smn_subscriptions" "filter_by_protocol" {
  protocol = local.protocol
}

data "huaweicloud_smn_subscriptions" "filter_by_status" {
  status = local.status
}

data "huaweicloud_smn_subscriptions" "filter_by_endpoint" {
  endpoint = local.endpoint
}

locals {
  list_by_protocol = data.huaweicloud_smn_subscriptions.filter_by_protocol.subscriptions
  list_by_status   = data.huaweicloud_smn_subscriptions.filter_by_status.subscriptions
  list_by_endpoint = data.huaweicloud_smn_subscriptions.filter_by_endpoint.subscriptions
}

output "protocol_filter_is_useful" {
  value = length(local.list_by_protocol) > 0 && alltrue(
    [for v in local.list_by_protocol[*].protocol : v == local.protocol]
  )
}

output "status_filter_is_useful" {
  value = length(local.list_by_status) > 0 && alltrue(
    [for v in local.list_by_status[*].status : tostring(v) == local.status]
  )
}

output "endpoint_filter_is_useful" {
  value = length(local.list_by_endpoint) > 0 && alltrue(
    [for v in local.list_by_endpoint[*].endpoint : v == local.endpoint]
  )
}
`, testAccSMNV2SubscriptionConfig_basic(name))
}
