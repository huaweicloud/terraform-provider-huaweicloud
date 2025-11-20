package waf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePolicyIpReputation_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_waf_policy_ip_reputation.test"
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
				Config: testAccDataSourcePolicyIpReputation_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "ip_reputation_map.#"),
					resource.TestCheckResourceAttrSet(dataSource, "locale.%"),
				),
			},
		},
	})
}

const testAccDataSourcePolicyIpReputation_basic = `
data "huaweicloud_waf_policy_ip_reputation" "test" {
  lang = "cn"
  type = "idc"
}
`
