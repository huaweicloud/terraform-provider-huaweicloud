package rocketmq

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataBrokers_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_dms_rocketmq_brokers.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataBrokers_instanceNotFound(),
				ExpectError: regexp.MustCompile("This DMS instance does not exist"),
			},
			{
				Config: testAccDataBrokers_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "brokers.0", "broker-0"),
				),
			},
		},
	})
}

func testAccDataBrokers_base(name string) string {
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

  flavor_id         = "c6.4u8g.cluster"
  storage_spec_code = "dms.physical.storage.high.v2"
  broker_num        = 1
}
`, common.TestBaseNetwork(name), name)
}

func testAccDataBrokers_instanceNotFound() string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dms_rocketmq_brokers" "test" {
  instance_id = "%s"
}
`, randomId)
}

func testAccDataBrokers_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_rocketmq_brokers" "test" {
  instance_id = huaweicloud_dms_rocketmq_instance.test.id
}
`, testAccDataBrokers_base(name))
}
