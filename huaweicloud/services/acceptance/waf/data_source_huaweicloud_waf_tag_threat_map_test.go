package waf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTagThreatMap_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_waf_tag_threat_map.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTagThreatMap_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "threats.#"),
					resource.TestCheckResourceAttrSet(dataSource, "locale.#"),
					resource.TestCheckResourceAttrSet(dataSource, "locale.0.cmdi"),
					resource.TestCheckResourceAttrSet(dataSource, "locale.0.anticrawler"),
					resource.TestCheckResourceAttrSet(dataSource, "locale.0.custom_custom"),
					resource.TestCheckResourceAttrSet(dataSource, "locale.0.antiscan_dir_traversal"),
					resource.TestCheckResourceAttrSet(dataSource, "locale.0.webshell"),
					resource.TestCheckResourceAttrSet(dataSource, "locale.0.cc"),
					resource.TestCheckResourceAttrSet(dataSource, "locale.0.illegal"),
					resource.TestCheckResourceAttrSet(dataSource, "locale.0.antitamper"),
					resource.TestCheckResourceAttrSet(dataSource, "locale.0.custom_geoip"),
					resource.TestCheckResourceAttrSet(dataSource, "locale.0.vuln"),
					resource.TestCheckResourceAttrSet(dataSource, "locale.0.llm_response_sensitive"),
					resource.TestCheckResourceAttrSet(dataSource, "locale.0.custom_whiteblackip"),
					resource.TestCheckResourceAttrSet(dataSource, "locale.0.leakage"),
				),
			},
		},
	})
}

const testAccDataSourceTagThreatMap_basic = `data "huaweicloud_waf_tag_threat_map" "test" {}`
