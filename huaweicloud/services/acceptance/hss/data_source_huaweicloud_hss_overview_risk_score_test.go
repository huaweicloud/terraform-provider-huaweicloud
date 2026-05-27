package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOverviewRiskScore_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_overview_risk_score.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceOverviewRiskScore_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "score"),
					resource.TestCheckResourceAttrSet(dataSource, "risk_num"),
					resource.TestCheckResourceAttrSet(dataSource, "risk_server_num"),
				),
			},
		},
	})
}

func testDataSourceOverviewRiskScore_basic() string {
	return `
data "huaweicloud_hss_overview_risk_score" "test" {
  enterprise_project_id = "all_granted_eps"
}
`
}
