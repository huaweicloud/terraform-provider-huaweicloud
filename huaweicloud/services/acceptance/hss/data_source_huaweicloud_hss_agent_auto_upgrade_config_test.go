package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAgentAutoUpgradeConfig_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_agent_auto_upgrade_config.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAgentAutoUpgradeConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "enabled"),
				),
			},
		},
	})
}

func testDataSourceAgentAutoUpgradeConfig_basic() string {
	return `
data "huaweicloud_hss_agent_auto_upgrade_config" "test" {
  enterprise_project_id = "0"
}
`
}
