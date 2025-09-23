package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecmasterBaselineCheckResults_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_baseline_check_results.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSecmasterBaselineCheckResults_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					// There is currently no test data to support verification of response values.
					resource.TestCheckResourceAttrSet(dataSource, "baseline_check_results.#"),
				),
			},
		},
	})
}

func testDataSourceSecmasterBaselineCheckResults_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_baseline_check_results" "test" {
  workspace_id = "%s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
