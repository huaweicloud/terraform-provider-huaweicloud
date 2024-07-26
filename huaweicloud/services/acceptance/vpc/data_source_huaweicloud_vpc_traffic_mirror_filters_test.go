package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpcTrafficMirrorFilters_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpc_traffic_mirror_filters.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceVpcTrafficMirrorFilters_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "name", rName),
				),
			},
		},
	})
}

func testDataSourceDataSourceVpcTrafficMirrorFilters_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_traffic_mirror_filter" "test" {
  name        = "%[1]s"
  description = "tf acc test filter"
}

resource "huaweicloud_vpc_traffic_mirror_filter_rule" "test1" {
  traffic_mirror_filter_id = huaweicloud_vpc_traffic_mirror_filter.test.id
  ethertype                = "IPv4"
  direction                = "ingress"
  protocol                 = "tcp"
  action                   = "accept"
  priority                 = 1
  description              = "create VPC traffic mirror filter rule"
}

resource "huaweicloud_vpc_traffic_mirror_filter_rule" "test2" {
  traffic_mirror_filter_id = huaweicloud_vpc_traffic_mirror_filter.test.id
  ethertype                = "IPv4"
  direction                = "egress"
  protocol                 = "all"
  action                   = "accept"
  priority                 = 20
  source_cidr_block        = "192.168.1.0/24"
}

data "huaweicloud_vpc_traffic_mirror_filters" "test" {
  name = "%[1]s"

  depends_on = [
    huaweicloud_vpc_traffic_mirror_filter.test,
    huaweicloud_vpc_traffic_mirror_filter_rule.test1,
    huaweicloud_vpc_traffic_mirror_filter_rule.test2
  ]
}
`, name)
}
