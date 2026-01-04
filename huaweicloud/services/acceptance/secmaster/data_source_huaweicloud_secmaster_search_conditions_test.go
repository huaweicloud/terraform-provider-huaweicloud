package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSearchConditions_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_search_conditions.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterPipeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSearchConditions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.#"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.condition_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.condition_name"),
				),
			},
		},
	})
}

func testAccDataSourceSearchConditions_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_search_conditions" "test" {
  workspace_id = "%[1]s"
  pipe_id      = "%[2]s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_PIPE_ID)
}
