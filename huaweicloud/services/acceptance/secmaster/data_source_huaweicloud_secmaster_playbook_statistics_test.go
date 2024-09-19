package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecmasterPlaybookStatistics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_playbook_statistics.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSecmasterPlaybookStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
				),
			},
		},
	})
}

func testDataSourceSecmasterPlaybookStatistics_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_playbook_statistics" "test" {
  workspace_id = "%[1]s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
