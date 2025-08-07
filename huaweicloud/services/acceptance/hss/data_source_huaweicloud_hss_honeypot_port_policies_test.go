package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHoneypotPortPolicies_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_honeypot_port_policies.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID that has enabled premium edition host protection.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHoneypotPortPolicies_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.policy_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.policy_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.is_default"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.port_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.os_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.status"),
				),
			},
		},
	})
}

const testAccDataSourceHoneypotPortPolicies_basic = `data "huaweicloud_hss_honeypot_port_policies" "test" {}`
