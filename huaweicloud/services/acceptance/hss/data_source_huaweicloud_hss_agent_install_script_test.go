package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAgentInstallScript_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_agent_install_script.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAgentInstallScript_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "install_script_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "install_script_list.0.package_type"),
					resource.TestCheckResourceAttrSet(dataSource, "install_script_list.0.cmd"),
				),
			},
		},
	})
}

const testDataSourceAgentInstallScript_basic string = `
data "huaweicloud_hss_agent_install_script" "test" {
  os_type       = "Linux"
  os_arch       = "x86_64"
  outside_host  = true
  batch_install = true
  type          = "password"
}
`
