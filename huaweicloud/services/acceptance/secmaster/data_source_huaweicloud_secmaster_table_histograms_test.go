package secmaster

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTableHistograms_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testDataSourceSecmasterTableHistograms_basic(),
				ExpectError: regexp.MustCompile(`版本未升级，请使用老版本.`),
			},
		},
	})
}

// Due to the lack of a test environment, only expected failure scenarios can be tested at present.
// The value of `table_id` in the test script is mock data.
func testDataSourceSecmasterTableHistograms_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_table_histograms" "test" {
  workspace_id = "%s"
  table_id     = "7251f007-f723-477f-949d-5b0f697238bf"
  query        = "*"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
