package rabbitmq

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDmsRabbitmqQueues_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dms_rabbitmq_queues.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDmsRabbitmqQueues_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "queues.#"),
					resource.TestCheckResourceAttrSet(dataSource, "queues.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "queues.0.auto_delete"),
					resource.TestCheckResourceAttrSet(dataSource, "queues.0.durable"),
					resource.TestCheckResourceAttrSet(dataSource, "queues.0.dead_letter_exchange"),
					resource.TestCheckResourceAttrSet(dataSource, "queues.0.dead_letter_routing_key"),
					resource.TestCheckResourceAttrSet(dataSource, "queues.0.message_ttl"),
					resource.TestCheckResourceAttrSet(dataSource, "queues.0.lazy_mode"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDmsRabbitmqQueues_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_rabbitmq_queues" "test" {
  depends_on = [huaweicloud_dms_rabbitmq_queue.test]

  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  vhost       = huaweicloud_dms_rabbitmq_vhost.test.name
}
`, testRabbitmqQueue_basic(name))
}
