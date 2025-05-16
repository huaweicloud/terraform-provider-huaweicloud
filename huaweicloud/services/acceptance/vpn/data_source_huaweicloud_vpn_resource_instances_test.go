package vpn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpnInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpn_resource_instances.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVpnInstances_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.0.value"),

					resource.TestCheckOutput("without_any_tag_filter_is_useful", "true"),
					resource.TestCheckOutput("tags_filter_is_useful", "true"),
					resource.TestCheckOutput("matches_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceVpnInstances_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_vpn_resource_instances" "test" {
  depends_on = [huaweicloud_vpn_gateway.test]

  resource_type = "vpn-gateway"
}

data "huaweicloud_vpn_resource_instances" "without_any_tag_filter" {
  depends_on = [huaweicloud_vpn_gateway.test]

  resource_type   = "vpn-gateway"
  without_any_tag = true
}
	
output "without_any_tag_filter_is_useful" {
  value = length(data.huaweicloud_vpn_resource_instances.without_any_tag_filter.resources) >= 0 && alltrue(
    [for v in data.huaweicloud_vpn_resource_instances.without_any_tag_filter.resources[*].tags : length(v) == 0]
  )
}

data "huaweicloud_vpn_resource_instances" "tags_filter" {
  depends_on = [huaweicloud_vpn_gateway.test]

  resource_type = "vpn-gateway"

  tags {
    key    = "key"
    values = ["val"]
  }
}
	
output "tags_filter_is_useful" {
  value = length(data.huaweicloud_vpn_resource_instances.tags_filter.resources) > 0
}

data "huaweicloud_vpn_resource_instances" "matches_filter" {
  depends_on = [huaweicloud_vpn_gateway.test]

  resource_type = "vpn-gateway"

  matches {
    key   = "resource_name"
    value = "%[2]s"
  }
}
	
output "matches_filter_is_useful" {
  value = length(data.huaweicloud_vpn_resource_instances.matches_filter.resources) > 0 && alltrue(
    [for v in data.huaweicloud_vpn_resource_instances.matches_filter.resources[*].resource_name : v == "%[2]s"])
}
`, testGateway_basic(name), name)
}
