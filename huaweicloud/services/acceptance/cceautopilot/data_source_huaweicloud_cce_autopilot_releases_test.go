package cceautopilot

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCCEAutopilotReleases_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cce_autopilot_releases.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCCEAutopilotReleases_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "id"),
					resource.TestCheckResourceAttrSet(dataSource, "region"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.name"),
				),
			},
		},
	})
}

const testDataSourceCCEAutopilotReleases_basic = `
data "huaweicloud_cce_autopilot_releases" "test" {
  cluster_id = "83a085fe-d4f1-11f0-80d0-0255ac10178d"
}
`
