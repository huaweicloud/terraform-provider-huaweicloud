package waf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWafConfig_basic(t *testing.T) {
	var (
		datasource = "data.huaweicloud_waf_config.test"
		dc         = acceptance.InitDataSourceCheck(datasource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceWafConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(datasource, "eps"),
					resource.TestCheckResourceAttrSet(datasource, "tls"),
					resource.TestCheckResourceAttrSet(datasource, "ipv6"),
					resource.TestCheckResourceAttrSet(datasource, "alert"),
					resource.TestCheckResourceAttrSet(datasource, "cc_enhance"),
					resource.TestCheckResourceAttrSet(datasource, "geoip_enable"),
					resource.TestCheckResourceAttrSet(datasource, "ip_group"),
					resource.TestCheckResourceAttrSet(datasource, "custom"),
					resource.TestCheckResourceAttrSet(datasource, "js_crawler_enable"),
					resource.TestCheckResourceAttrSet(datasource, "robot_action_enable"),
				),
			},
		},
	})
}

const testDataSourceWafConfig_basic = `data "huaweicloud_waf_config" "test" {}`
