package cbh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDatasourceCbhInstances_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_cbh_instances.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceCbhInstances_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "instances.0.name", name),
					resource.TestCheckResourceAttr(dataSourceName, "instances.0.bastion_type", "OEM"),
					resource.TestCheckResourceAttr(dataSourceName, "instances.0.flavor_id", "cbh.basic.50"),
				),
			},
		},
	})
}

func testAccDatasourceCbhInstances_base(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_cbh_instance" "test" {
  flavor_id          = "cbh.basic.50"
  name               = "%s"
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  region             = "%s"
  hx_password        = "test_123456"
  bastion_type       = "OEM"
  charging_mode      = "prePaid"
  period_unit        = "month"
  auto_renew         = "false"
  period             = "1"
  
  product_info {
    product_id         = "OFFI740586375358963717"
    resource_size      = "1"
  }
}
`, common.TestBaseNetwork(name), name, acceptance.HW_REGION_NAME)
}

func testAccDatasourceCbhInstances_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cbh_instances" "test" {
  name = huaweicloud_cbh_instance.test.name
}
`, testAccDatasourceCbhInstances_base(name))
}
