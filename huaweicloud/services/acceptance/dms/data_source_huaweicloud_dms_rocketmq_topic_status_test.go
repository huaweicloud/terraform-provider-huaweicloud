package dms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDmsRocketMQTopicStatus_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	resourceName := "data.huaweicloud_dms_rocketmq_topic_status.test"
	dc := acceptance.InitDataSourceCheck(resourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDmsRocketMQTopicStatus_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "brokers.0.name", "broker-0"),
					resource.TestCheckResourceAttr(resourceName, "brokers.0.queues.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "brokers.0.queues.0.id", "0"),
					resource.TestCheckResourceAttr(resourceName, "brokers.0.queues.0.min_offset", "0"),
					resource.TestCheckResourceAttr(resourceName, "brokers.0.queues.0.max_offset", "0"),
					resource.TestCheckResourceAttr(resourceName, "brokers.0.queues.0.last_message_time", "0"),
					resource.TestCheckResourceAttr(resourceName, "brokers.0.queues.1.id", "1"),
					resource.TestCheckResourceAttr(resourceName, "brokers.0.queues.1.min_offset", "0"),
					resource.TestCheckResourceAttr(resourceName, "brokers.0.queues.1.max_offset", "0"),
					resource.TestCheckResourceAttr(resourceName, "brokers.0.queues.1.last_message_time", "0"),
					resource.TestCheckResourceAttr(resourceName, "brokers.0.queues.2.id", "2"),
					resource.TestCheckResourceAttr(resourceName, "brokers.0.queues.2.min_offset", "0"),
					resource.TestCheckResourceAttr(resourceName, "brokers.0.queues.2.max_offset", "0"),
					resource.TestCheckResourceAttr(resourceName, "brokers.0.queues.2.last_message_time", "0"),
				),
			},
		},
	})
}

func testAccDatasourceDmsRocketMQTopicStatus_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name        = "%[1]s"
  cidr        = "192.168.0.0/24"
  description = "test for DMS RocketMQ"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%[1]s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.test.id
}

resource "huaweicloud_networking_secgroup" "test" {
  name        = "%[1]s"
  description = "secgroup for rocketmq"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dms_rocketmq_instance" "test" {
  name              = "%[1]s"
  engine_version    = "4.8.0"
  storage_space     = 600
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  flavor_id         = "c6.4u8g.cluster.small"
  storage_spec_code = "dms.physical.storage.high.v2"
  broker_num        = 1
}


resource "huaweicloud_dms_rocketmq_topic" "test" {
  instance_id = huaweicloud_dms_rocketmq_instance.test.id
  name        = "%[1]s"
  permission  = "all"
  queue_num   = 3
  
  brokers {
    name = "broker-0"
  }
}

data "huaweicloud_dms_rocketmq_topic_status" "test" {
  instance_id = huaweicloud_dms_rocketmq_instance.test.id
  topic       = huaweicloud_dms_rocketmq_topic.test.name
}
`, name)
}
