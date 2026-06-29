package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePolicyRecords_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_policy_records.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterPolicyID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePolicyRecords_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "workspace_id"),
					resource.TestCheckResourceAttrSet(dataSource, "policy_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.#"),
				),
			},
		},
	})
}

func testDataSourcePolicyRecords_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_policy_records" "test" {
  workspace_id = "%[1]s"
  policy_id    = "%[2]s"

  sort {
    sort_by = "create_time"
    order   = "desc"
  }

  group_by {
    group_by_fields = ["workspace_id"]
    group_by_hit {
      source = "defense_policy_object"
      dest   = "defense_policy_list"
    }
  }
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_POLICY_ID)
}
