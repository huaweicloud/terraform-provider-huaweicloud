package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHssPolicySwitchStatus_basic(t *testing.T) {
	dataSource := "data.huaweicloud_hss_policy_switch_status.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceHssPolicySwitchStatus_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "enable"),
				),
			},
		},
	})
}

const testDataSourceDataSourceHssPolicySwitchStatus_basic = `
data "huaweicloud_hss_policy_switch_status" "test" {
  enterprise_project_id = "all_granted_eps"
  policy_name           = "sp_feature"
}`
