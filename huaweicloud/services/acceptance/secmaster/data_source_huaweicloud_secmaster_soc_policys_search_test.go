package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSocPolicysSearch_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_secmaster_soc_policys_search.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSocPolicysSearch_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "data"),
				),
			},
		},
	})
}

func testDataSourceSocPolicysSearch_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_soc_policys_search" "test" {
  workspace_id = "%[1]s"
  sort {
    sort_by = "create_time"
    order   = "desc"
  }
  condition {
    conditions {
      name = "policy_name"
      data = ["test"]
    }
  }
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
