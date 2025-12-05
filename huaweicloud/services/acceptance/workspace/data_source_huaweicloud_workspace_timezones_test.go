package workspace

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataTimezones_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_workspace_timezones.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataTimezones_basic,
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "time_zones.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(all, "time_zones.0.name"),
					resource.TestCheckResourceAttrSet(all, "time_zones.0.offset"),
					resource.TestCheckResourceAttrSet(all, "time_zones.0.us_description"),
					resource.TestCheckResourceAttrSet(all, "time_zones.0.cn_description"),
				),
			},
		},
	})
}

const testAccDataTimezones_basic = `
data "huaweicloud_workspace_timezones" "all" {}
`
