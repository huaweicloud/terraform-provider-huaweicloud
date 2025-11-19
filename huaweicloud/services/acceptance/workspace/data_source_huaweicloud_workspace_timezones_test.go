package workspace

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTimezones_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_workspace_timezones.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTimezones_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "time_zones.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "time_zones.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "time_zones.0.offset"),
					resource.TestCheckResourceAttrSet(dataSourceName, "time_zones.0.us_description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "time_zones.0.cn_description"),
				),
			},
		},
	})
}

const testAccDataSourceTimezones_basic = `
data "huaweicloud_workspace_timezones" "test" {}
`
