package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHoneypotPortSupportList_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_honeypot_port_support_list.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires the existence of a host with flagship version host protection enabled.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceHoneypotPortSupportList_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.host_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.host_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.private_ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.agent_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.os_type"),
				),
			},
		},
	})
}

func testDataSourceHoneypotPortSupportList_basic() string {
	return `
data "huaweicloud_hss_honeypot_port_support_list" "test" {
  os_type               = "Linux"
  enterprise_project_id = "all_granted_eps"
}
`
}
