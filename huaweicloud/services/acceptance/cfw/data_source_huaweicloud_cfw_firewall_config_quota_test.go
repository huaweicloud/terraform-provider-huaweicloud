package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceFirewallConfigQuota_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cfw_firewall_config_quota.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a firewall instance ID.
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceFirewallConfigQuota_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.item_info.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.item_info.0.max_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.item_info.0.used_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.item_info.0.extras_info.%"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.max_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.quota_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.used_quota"),
				),
			},
		},
	})
}

func testDataSourceFirewallConfigQuota_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_firewall_config_quota" "test" {
  fw_instance_id = "%s"
  config_type    = "DNS_DOMAIN_SET"
}
`, acceptance.HW_CFW_INSTANCE_ID)
}
