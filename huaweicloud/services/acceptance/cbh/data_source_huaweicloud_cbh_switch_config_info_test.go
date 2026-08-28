package cbh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwitchConfigInfo_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_cbh_switch_config_info.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSwitchConfigInfo_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "switch_info.0.is_support_unibuy"),
					resource.TestCheckResourceAttrSet(dataSourceName, "switch_info.0.is_support_iam_login"),
					resource.TestCheckResourceAttrSet(dataSourceName, "switch_info.0.is_support_ha"),
					resource.TestCheckResourceAttrSet(dataSourceName, "version_info.0.require_eip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "version_info.0.iam_login"),
				),
			},
		},
	})
}

const testAccDataSourceSwitchConfigInfo_basic string = `
data "huaweicloud_cbh_switch_config_info" "test" {}
`
