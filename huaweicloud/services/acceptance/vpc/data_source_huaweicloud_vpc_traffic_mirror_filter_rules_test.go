package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpcTrafficMirrorFilterRules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpc_traffic_mirror_filter_rules.test1"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceVpcTrafficMirrorFilterRules_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "direction", "ingress"),
					resource.TestCheckResourceAttr(dataSource, "protocol", "tcp"),
					resource.TestCheckResourceAttr(dataSource, "action", "accept"),
					resource.TestCheckResourceAttr(dataSource, "priority", "555"),
				),
			},
		},
	})
}

func testDataSourceDataSourceVpcTrafficMirrorFilterRules_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_traffic_mirror_filter" "test1" {
  name        = "%[1]s"
  description = "tf acc test"
}

resource "huaweicloud_vpc_traffic_mirror_filter_rule" "test1" {
  traffic_mirror_filter_id = huaweicloud_vpc_traffic_mirror_filter.test1.id
  direction                = "ingress"
  protocol                 = "tcp"
  ethertype                = "IPv4"
  action                   = "accept"
  priority                 = "333"
}

resource "huaweicloud_vpc_traffic_mirror_filter_rule" "test2" {
  traffic_mirror_filter_id = huaweicloud_vpc_traffic_mirror_filter.test1.id
  direction                = "ingress"
  protocol                 = "tcp"
  ethertype                = "IPv4"
  action                   = "accept"
  priority                 = "555"
}

resource "huaweicloud_vpc_traffic_mirror_filter_rule" "test3" {
  traffic_mirror_filter_id = huaweicloud_vpc_traffic_mirror_filter.test1.id
  direction                = "egress"
  protocol                 = "tcp"
  ethertype                = "IPv4"
  action                   = "reject"
  priority                 = "666"
}

data "huaweicloud_vpc_traffic_mirror_filter_rules" "test1" {
  direction = "ingress"
  protocol  = "tcp"
  action    = "accept"
  priority  = "555"
}
`, name)
}
