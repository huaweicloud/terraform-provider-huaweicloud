package cceautopilot

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceClusterLogConfigs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cce_autopilot_cluster_log_configs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCceClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceClusterLogConfigs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "ttl_in_days"),
					resource.TestCheckResourceAttrSet(dataSource, "log_configs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "log_configs.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "log_configs.0.enable"),
					resource.TestCheckResourceAttrSet(dataSource, "log_configs.0.name"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceClusterLogConfigs_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cce_autopilot_cluster_log_configs" "test" {
  cluster_id = "%[1]s"
}
data "huaweicloud_cce_autopilot_cluster_log_configs" "type_filter" {
  cluster_id = "%[1]s"
  type       = "control"
}
output "type_filter_is_useful" {
  value = length(data.huaweicloud_cce_autopilot_cluster_log_configs.type_filter.log_configs) > 0
}
`, acceptance.HW_CCE_CLUSTER_ID)
}
