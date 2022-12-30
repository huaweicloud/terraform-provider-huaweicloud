package dms

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDmsRocketMQGroupTopics_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	resourceName := "data.huaweicloud_dms_rocketmq_group_topic.test"
	dc := acceptance.InitDataSourceCheck(resourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDmsRocketMQGroupTopics_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(resourceName, "topics.#", regexp.MustCompile("[0-9]\\d*")),
				),
			},
		},
	})
}

func testAccDatasourceDmsRocketMQGroupTopics_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name        = "%[1]s"
  cidr        = "192.168.0.0/24"
  description = "Test for DMS RocketMQ"
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

resource "huaweicloud_dms_rocketmq_consumer_group" "test" {
  instance_id = huaweicloud_dms_rocketmq_instance.test.id
  broadcast   = "true"

  brokers = [
    "broker-0"
  ]

  name           = "%[1]s"
  retry_max_time = "3"
}

data "huaweicloud_dms_rocketmq_group_topic" "test" {
  instance_id = huaweicloud_dms_rocketmq_instance.test.id
  group       = huaweicloud_dms_rocketmq_consumer_group.test.name
}
`, name)
}
