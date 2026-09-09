package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOfflineKeyAnalysisNodes_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dcs_offline_key_analysis_nodes.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDCSInstanceID(t)
			acceptance.TestAccPreCheckDCSTaskId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOfflineKeyAnalysisNodes_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.group_name"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.node_ipv6"),
				),
			},
		},
	})
}

func testAccDataSourceOfflineKeyAnalysisNodes_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dcs_offline_key_analysis_nodes" "test" {
  instance_id = "%[1]s"
  task_id     = "%[2]s"
}
`, acceptance.HW_DCS_INSTANCE_ID, acceptance.HW_DCS_TASK_ID)
}
