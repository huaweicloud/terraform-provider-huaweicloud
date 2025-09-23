package waf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAlarmOptionalEventTypes_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_waf_alarm_optional_event_types.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAlarmOptionalEventTypes_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "threats.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.0.cmdi"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.0.llm_prompt_injection"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.0.anticrawler"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.0.custom_custom"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.0.third_bot_river"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.0.robot"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.0.custom_idc_ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.0.antiscan_dir_traversal"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.0.advanced_bot"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.0.xss"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.0.antiscan_high_freq_scan"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.0.webshell"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.0.cc"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.0.botm"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.0.illegal"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.0.llm_prompt_sensitive"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.0.sqli"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.0.lfi"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.0.antitamper"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.0.custom_geoip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.0.rfi"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.0.vuln"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.0.llm_response_sensitive"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.0.custom_whiteblackip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.0.leakage"),
				),
			},
		},
	})
}

const testDataSourceAlarmOptionalEventTypes_basic = `
data "huaweicloud_waf_alarm_optional_event_types" "test" {}
`
