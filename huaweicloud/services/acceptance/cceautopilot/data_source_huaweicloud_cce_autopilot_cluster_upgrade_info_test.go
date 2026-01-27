package cceautopilot

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceClusterUpgradeInfo_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cce_autopilot_cluster_upgrade_info.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCceClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceClusterUpgradeInfo_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "kind"),
					resource.TestCheckResourceAttrSet(dataSource, "api_version"),
					resource.TestCheckResourceAttrSet(dataSource, "metadata.#"),
					resource.TestCheckResourceAttrSet(dataSource, "spec.#"),
					resource.TestCheckResourceAttrSet(dataSource, "spec.0.last_upgrade_info.#"),
					resource.TestCheckResourceAttrSet(dataSource, "spec.0.last_upgrade_info.0.phase"),
					resource.TestCheckResourceAttrSet(dataSource, "spec.0.last_upgrade_info.0.completion_time"),
					resource.TestCheckResourceAttrSet(dataSource, "spec.0.version_info.#"),
					resource.TestCheckResourceAttrSet(dataSource, "spec.0.version_info.0.release"),
					resource.TestCheckResourceAttrSet(dataSource, "spec.0.version_info.0.patch"),
					resource.TestCheckResourceAttrSet(dataSource, "spec.0.version_info.0.suggest_patch"),
					resource.TestCheckResourceAttrSet(dataSource, "spec.0.version_info.0.target_versions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "spec.0.upgrade_feature_gates.#"),
					resource.TestCheckResourceAttrSet(dataSource, "spec.0.upgrade_feature_gates.0.support_upgrade_page_v4"),
					resource.TestCheckResourceAttrSet(dataSource, "status.#"),
					resource.TestCheckResourceAttrSet(dataSource, "status.0.phase"),
					resource.TestCheckResourceAttrSet(dataSource, "status.0.completion_time"),
				),
			},
		},
	})
}

func testAccDataSourceClusterUpgradeInfo_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cce_autopilot_cluster_upgrade_info" "test" {
  cluster_id = "%s"
}
`, acceptance.HW_CCE_CLUSTER_ID)
}
