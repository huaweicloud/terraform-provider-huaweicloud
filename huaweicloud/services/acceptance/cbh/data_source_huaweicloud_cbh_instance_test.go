package cbh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
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
resource "huaweicloud_vpc" "test" {
  name        = "%[1]s"
  cidr        = "192.168.0.0/16"
  description = "Test for CBH instance"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%[1]s"
  cidr       = "192.168.0.0/20"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.test.id
}

resource "huaweicloud_networking_secgroup" "test" {
  name        = "%[1]s"
  description = "secgroup for CBH instance"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_cbh_instance" "test" {
  flavor_id = "cbh.basic.50"
  name      = "%[1]s"
  vpc_id    = huaweicloud_vpc.test.id
  nics {
    subnet_id = huaweicloud_vpc_subnet.test.id
  }
  security_groups {
    id = huaweicloud_networking_secgroup.test.id
  }
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  region            = "%[2]s"
  hx_password       = "test_123456"
  bastion_type      = "OEM"
  charging_mode      = "prePaid"
  period_unit        = "month"
  auto_renew         = "false"
  period             = "1"
  subscription_num   = "1"
  
  product_info {
    flavor_id          = "OFFI740586375358963717"
    resource_size      = "1"
  }
}
`, name, acceptance.HW_REGION_NAME)
}

func testAccDatasourceCbhInstances_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cbh_instances" "test" {
  name = huaweicloud_cbh_instance.test.name
}
`, testAccDatasourceCbhInstances_base(name))
}
