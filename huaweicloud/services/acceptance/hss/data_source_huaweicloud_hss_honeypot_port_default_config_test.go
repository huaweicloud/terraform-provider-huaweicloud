package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHoneypotPortDefaultConfig_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_honeypot_port_default_config.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceHoneypotPortDefaultConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "auto_bind"),
					resource.TestCheckResourceAttrSet(dataSource, "windows_policy_id"),
					resource.TestCheckResourceAttrSet(dataSource, "linux_policy_id"),
				),
			},
		},
	})
}

func testDataSourceHoneypotPortDefaultConfig_basic() string {
	return `
data "huaweicloud_hss_honeypot_port_default_config" "test" {
  enterprise_project_id = "0"
}
`
}
