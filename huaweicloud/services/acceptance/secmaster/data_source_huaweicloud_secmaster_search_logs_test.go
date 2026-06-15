package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSearchLogs_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_search_logs.test"
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
				Config: testAccDataSourceSearchLogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "results.#"),
					resource.TestCheckResourceAttrSet(dataSource, "analysis_results.#"),
				),
			},
		},
	})
}

func testAccDataSourceSearchLogs_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_search_logs" "test" {
  workspace_id = "%[1]s"
  dataspace_id = "%[2]s"
  pipe_id      = "%[3]s"
  query        = "*"
  from         = 1781506114075
  to           = 1781507014075
  sort         = "desc"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_DATASPACE_ID, acceptance.HW_SECMASTER_PIPE_ID)
}
