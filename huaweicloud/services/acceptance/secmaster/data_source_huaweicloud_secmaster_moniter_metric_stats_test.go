package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceMoniterMetricStats_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_moniter_metric_stats.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterDataspaceId(t)
			acceptance.TestAccPreCheckSecMasterPipeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMoniterMetricStats_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "results.#"),
				),
			},
		},
	})
}

func testAccDataSourceMoniterMetricStats_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_moniter_metric_stats" "test" {
  workspace_id    = "%[1]s"
  dataspace_id    = "%[2]s"
  pipe_id         = "%[3]s"
  start_timestamp = 1780313858887
  end_timestamp   = 1780918658887
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_DATASPACE_ID, acceptance.HW_SECMASTER_PIPE_ID)
}
