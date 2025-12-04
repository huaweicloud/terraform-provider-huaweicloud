package lb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceLbWhitelists_basic(t *testing.T) {
	dataSource := "data.huaweicloud_lb_whitelists.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceLbWhitelists_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "whitelists.#"),
					resource.TestCheckResourceAttrSet(dataSource, "whitelists.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "whitelists.0.listener_id"),
					resource.TestCheckResourceAttrSet(dataSource, "whitelists.0.enable_whitelist"),
					resource.TestCheckResourceAttrSet(dataSource, "whitelists.0.whitelist"),
					resource.TestCheckOutput("whitelist_id_filter_is_useful", "true"),
					resource.TestCheckOutput("enable_whitelist_filter_is_useful", "true"),
					resource.TestCheckOutput("listener_id_filter_is_useful", "true"),
					resource.TestCheckOutput("whitelist_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceLbWhitelists_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_lb_whitelists" "test" {
  depends_on = [huaweicloud_lb_whitelist.whitelist_1]
}

locals {
  whitelist_id = huaweicloud_lb_whitelist.whitelist_1.id
}
data "huaweicloud_lb_whitelists" "whitelist_id_filter" {
  whitelist_id = huaweicloud_lb_whitelist.whitelist_1.id
}
output "whitelist_id_filter_is_useful" {
  value = length(data.huaweicloud_lb_whitelists.whitelist_id_filter.whitelists) > 0 && alltrue(
    [for v in data.huaweicloud_lb_whitelists.whitelist_id_filter.whitelists[*].id : v == local.whitelist_id]
  )
}

locals {
  enable_whitelist = huaweicloud_lb_whitelist.whitelist_1.enable_whitelist
}
data "huaweicloud_lb_whitelists" "enable_whitelist_filter" {
  depends_on = [huaweicloud_lb_whitelist.whitelist_1]

  enable_whitelist = huaweicloud_lb_whitelist.whitelist_1.enable_whitelist
}
output "enable_whitelist_filter_is_useful" {
  value = length(data.huaweicloud_lb_whitelists.enable_whitelist_filter.whitelists) > 0 && alltrue(
    [for v in data.huaweicloud_lb_whitelists.enable_whitelist_filter.whitelists[*].enable_whitelist : v == local.enable_whitelist]
  )
}

locals {
  listener_id = huaweicloud_lb_whitelist.whitelist_1.listener_id
}
data "huaweicloud_lb_whitelists" "listener_id_filter" {
  depends_on = [huaweicloud_lb_whitelist.whitelist_1]

  listener_id = huaweicloud_lb_whitelist.whitelist_1.listener_id
}
output "listener_id_filter_is_useful" {
  value = length(data.huaweicloud_lb_whitelists.listener_id_filter.whitelists) > 0 && alltrue(
    [for v in data.huaweicloud_lb_whitelists.listener_id_filter.whitelists[*].listener_id : v == local.listener_id]
  )
}

locals {
  whitelist = huaweicloud_lb_whitelist.whitelist_1.whitelist
}
data "huaweicloud_lb_whitelists" "whitelist_filter" {
  depends_on = [huaweicloud_lb_whitelist.whitelist_1]

  whitelist = huaweicloud_lb_whitelist.whitelist_1.whitelist
}
output "whitelist_filter_is_useful" {
  value = length(data.huaweicloud_lb_whitelists.whitelist_filter.whitelists) > 0 && alltrue(
    [for v in data.huaweicloud_lb_whitelists.whitelist_filter.whitelists[*].whitelist : v == local.whitelist]
  )
}
`, testAccLBV2WhitelistConfig_basic(name))
}
