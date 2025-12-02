package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAppAgentStatistics_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_app_agent_statistics.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case need to a host that has enabled host protection.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppAgentStatistics_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "agent_num"),
				),
			},
		},
	})
}

const testAccDataSourceAppAgentStatistics_basic = `
data "huaweicloud_hss_app_agent_statistics" "test" {
  enterprise_project_id = "0"
}
`
