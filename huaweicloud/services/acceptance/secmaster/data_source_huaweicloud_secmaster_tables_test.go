package secmaster

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTables_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testDataSourceTables_basic(),
				ExpectError: regexp.MustCompile(`版本未升级，请使用老版本.`),
			},
		},
	})
}

// Due to the lack of a test environment, only expected failure scenarios can be tested at present.
func testDataSourceTables_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_tables" "test" {
  workspace_id = "%[1]s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
