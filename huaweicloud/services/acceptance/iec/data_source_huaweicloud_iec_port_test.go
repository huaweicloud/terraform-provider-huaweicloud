package iec

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccPortDataSource_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_iec_port.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIECNetworkConfig_base(rName),
			},
			{
				Config: testAccIECPortDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "mac_address"),
					resource.TestCheckResourceAttrSet(dataSourceName, "site_id"),
				),
			},
		},
	})
}

func testAccIECNetworkConfig_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_iec_sites" "test" {}

resource "huaweicloud_iec_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
  mode = "CUSTOMER"
}

resource "huaweicloud_iec_vpc_subnet" "test" {
  name       = "%[1]s"
  cidr       = "192.168.0.0/24"
  vpc_id     = huaweicloud_iec_vpc.test.id
  site_id    = data.huaweicloud_iec_sites.test.sites[0].id
  gateway_ip = "192.168.0.1"
}
`, rName)
}

func testAccIECPortDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_iec_port" "test" {
  fixed_ip  = huaweicloud_iec_vpc_subnet.test.gateway_ip
  subnet_id = huaweicloud_iec_vpc_subnet.test.id
}
`, testAccIECNetworkConfig_base(rName))
}
