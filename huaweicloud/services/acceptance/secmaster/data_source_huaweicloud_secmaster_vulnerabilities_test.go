package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVulnerabilities_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_vulnerabilities.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVulnerabilities_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
				),
			},
		},
	})
}

func testDataSourceVulnerabilities_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_vulnerabilities" "test" {
  workspace_id = "%[1]s"
  from_date    = "2023-02-20T00:00:00.000Z"
  to_date      = "2025-02-27T23:59:59.999Z"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
