package rocketmq

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDatasourceDmsRocketMQInstances_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_dms_rocketmq_instances.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDmsRocketMQInstances_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "instances.0.name", name),
					resource.TestCheckResourceAttr(rName, "instances.0.engine_version", "4.8.0"),
					resource.TestCheckResourceAttr(rName, "instances.0.flavor_id", "c6.4u8g.cluster.small"),
					resource.TestCheckResourceAttr(rName, "instances.0.broker_num", "1"),

					resource.TestCheckResourceAttrSet(rName, "instances.0.name"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.engine_version"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.flavor_id"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.broker_num"),
				),
			},
		},
	})
}

func testAccDatasourceDmsRocketMQInstances_config(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dms_rocketmq_instance" "test" {
  name              = "%s"
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
`, common.TestBaseNetwork(name), name)
}

func testAccDatasourceDmsRocketMQInstances_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_rocketmq_instances" "test" {
  name = huaweicloud_dms_rocketmq_instance.test.name
}
`, testAccDatasourceDmsRocketMQInstances_config(name))
}
