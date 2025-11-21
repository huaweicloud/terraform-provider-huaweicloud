package rgc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBestPracticeOverview_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_rgc_best_practice_overview.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBestPracticeOverview_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_score"),
					resource.TestCheckResourceAttrSet(dataSourceName, "detect_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "organization_account.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "identity_permission.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_planning.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "compliance_audit.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "financial_governance.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_boundary.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "om_monitor.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "security_management.#"),
				),
			},
		},
	})
}

const testAccDataSourceBestPracticeOverview_basic = `
data "huaweicloud_rgc_best_practice_overview" "test" {}
`
