package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVpcRouteTableDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "data.huaweicloud_vpc_route_table.rtb"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRouteTable_base(rName),
			},
			{
				Config: testAccDataSourceRouteTable_default(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "default", "true"),
					resource.TestCheckResourceAttr(resourceName, "subnets.#", "1"),
				),
			},
			{
				Config: testAccDataSourceRouteTable_custom(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "default", "false"),
					resource.TestCheckResourceAttr(resourceName, "subnets.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceRouteTable_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "172.16.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "172.16.10.0/24"
  gateway_ip = "172.16.10.1"
  vpc_id     = huaweicloud_vpc.test.id
}
`, rName, rName)
}

func testAccDataSourceRouteTable_default(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc_route_table" "rtb" {
  vpc_id = huaweicloud_vpc.test.id
}
`, testAccDataSourceRouteTable_base(rName))
}

func testAccDataSourceRouteTable_custom(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_route_table" "test" {
  name        = "%s"
  vpc_id      = huaweicloud_vpc.test.id
  description = "created by terraform"
}

data "huaweicloud_vpc_route_table" "rtb" {
  vpc_id = huaweicloud_vpc.test.id
  name   = "%s"

  depends_on = [huaweicloud_vpc_route_table.test]
}
`, testAccDataSourceRouteTable_base(rName), rName, rName)
}
