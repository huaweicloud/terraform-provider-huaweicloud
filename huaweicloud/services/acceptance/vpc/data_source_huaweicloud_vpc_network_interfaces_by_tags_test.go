package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceNetworkInterfacesByTags_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpc_network_interfaces_by_tags.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceNetworkInterfacesByTags_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_detail"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.%"),
					resource.TestCheckOutput("filter_by_tags_is_useful", "true"),
					resource.TestCheckOutput("filter_by_matches_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceNetworkInterfacesByTags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc_network_interface" "test" {
  name      = "%[2]s"
  subnet_id = huaweicloud_vpc_subnet.test.id

  tags = {
    key   = "value"
    owner = "terraform"
  }
}

data "huaweicloud_vpc_network_interfaces_by_tags" "basic" {
  depends_on = [huaweicloud_vpc_network_interface.test]
}

data "huaweicloud_vpc_network_interfaces_by_tags" "filter_by_tags" {
  depends_on = [huaweicloud_vpc_network_interface.test]

  tags {
    key    = "key"
    values = ["value"]
  }
}
output "filter_by_tags_is_useful" {
  value = length(data.huaweicloud_vpc_network_interfaces_by_tags.filter_by_tags) > 0
}

data "huaweicloud_vpc_network_interfaces_by_tags" "filter_by_matches" {
  depends_on = [huaweicloud_vpc_network_interface.test]

  matches {
    key   = "resource_name"
    value = "terraform"
  }
}
output "filter_by_matches_is_useful" {
  value = length(data.huaweicloud_vpc_network_interfaces_by_tags.filter_by_matches) > 0
}
`, common.TestVpc(name), name)
}
