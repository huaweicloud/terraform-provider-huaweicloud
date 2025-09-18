package waf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTagIpReputationMap_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_waf_tag_ip_reputation_map.test"
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
				Config: testAccDataSourceTagIpReputationMap_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "ip_reputation_map.#"),
					resource.TestCheckResourceAttrSet(dataSource, "ip_reputation_map.0.idc.#"),
				),
			},
		},
	})
}

const testAccDataSourceTagIpReputationMap_basic = `
data "huaweicloud_waf_tag_ip_reputation_map" "test" {
  lang = "en"
  type = "idc"
}`
