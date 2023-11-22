package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIECPortDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "data.huaweicloud_iec_port.port_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIecVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIECNetworkConfig_base(rName),
			},
			{
				Config: testAccIECPortDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "mac_address"),
					resource.TestCheckResourceAttrSet(resourceName, "site_id"),
				),
			},
		},
	})
}

func testAccIECPortDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_iec_port" "port_1" {
  fixed_ip  = huaweicloud_iec_vpc_subnet.subnet_1.gateway_ip
  subnet_id = huaweicloud_iec_vpc_subnet.subnet_1.id
}
`, testAccIECNetworkConfig_base(rName))
}
