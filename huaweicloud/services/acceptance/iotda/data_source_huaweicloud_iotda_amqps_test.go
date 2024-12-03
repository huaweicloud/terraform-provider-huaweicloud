package iotda

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAMQPQueues_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_iotda_amqps.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		name           = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAMQPQueues_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "queues.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "queues.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "queues.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "queues.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "queues.0.updated_at"),

					resource.TestCheckOutput("queue_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAMQPQueues_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_iotda_amqps" "test" {
  depends_on = [
    huaweicloud_iotda_amqp.test
  ]
}

locals {
  queue_id = data.huaweicloud_iotda_amqps.test.queues[0].id
}

data "huaweicloud_iotda_amqps" "queue_id_filter" {
  queue_id = local.queue_id
}

output "queue_id_filter_is_useful" {
  value = length(data.huaweicloud_iotda_amqps.queue_id_filter.queues) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_amqps.queue_id_filter.queues[*].id : v == local.queue_id]
  )
}

locals {
  name = data.huaweicloud_iotda_amqps.test.queues[0].name
}

data "huaweicloud_iotda_amqps" "name_filter" {
  name = local.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_iotda_amqps.name_filter.queues) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_amqps.name_filter.queues[*].name : v == local.name]
  )
}

data "huaweicloud_iotda_amqps" "not_found" {
  name = "not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_iotda_amqps.not_found.queues) == 0
}
`, testAmqp_basic(name))
}
