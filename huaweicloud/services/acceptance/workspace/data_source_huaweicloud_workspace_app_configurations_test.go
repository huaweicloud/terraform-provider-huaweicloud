package workspace

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAppConfigurations_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_workspace_app_configurations.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppConfigurations_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "configurations.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "configurations.0.config_key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "configurations.0.config_value"),
				),
			},
		},
	})
}

const testAccDataSourceAppConfigurations_basic = `
data "huaweicloud_workspace_app_configurations" "test" {
  items = ["help_center_origin", "IP_VIRTUAL_ENABLE"]
}
`
