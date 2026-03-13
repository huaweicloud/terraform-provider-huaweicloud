package dws

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataMonitorIndicators_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dws_monitor_indicators.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataMonitorIndicators_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "indicators.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "indicators.0.indicator_name"),
					resource.TestCheckResourceAttrSet(all, "indicators.0.plugin_name"),
					resource.TestCheckResourceAttrSet(all, "indicators.0.default_collect_rate"),
					resource.TestCheckResourceAttrSet(all, "indicators.0.support_datastore_version"),
				),
			},
		},
	})
}

const testAccDataMonitorIndicators_basic = `
data "huaweicloud_dws_monitor_indicators" "test" {}
`
