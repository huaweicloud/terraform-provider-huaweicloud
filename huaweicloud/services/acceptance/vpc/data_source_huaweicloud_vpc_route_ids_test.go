package vpc

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVpcRouteIdsDataSource_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()

	dataSourceName := "data.huaweicloud_vpc_route_ids.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckDeprecated(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteIdsDataSource_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "ids.#", "1"),
				),
			},
		},
	})
}

func testAccRouteIdsDataSource_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test1" {
  name = "%s_1"
  cidr = "172.16.0.0/20"
}

resource "huaweicloud_vpc" "test2" {
  name = "%s_2"
  cidr = "172.16.128.0/20"
}
  
resource "huaweicloud_vpc_peering_connection" "test" {
  name        = "%s_1"
  vpc_id      = huaweicloud_vpc.test1.id
  peer_vpc_id = huaweicloud_vpc.test2.id
}

resource "huaweicloud_vpc_route" "test" {
  type        = "peering"
  nexthop     = huaweicloud_vpc_peering_connection.test.id
  destination = huaweicloud_vpc.test2.cidr
  vpc_id      = huaweicloud_vpc.test1.id
}
`, rName, rName, rName)
}

func testAccRouteIdsDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc_route_ids" "test" {
  vpc_id = huaweicloud_vpc_route.test.vpc_id
}
`, testAccRouteIdsDataSource_base(rName))
}
