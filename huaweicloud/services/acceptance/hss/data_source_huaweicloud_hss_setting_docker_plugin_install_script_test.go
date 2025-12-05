package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSettingDockerPluginInstallScript_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_setting_docker_plugin_install_script.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{

			{
				Config: testAccDataSourceSettingDockerPluginInstallScript_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.package_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.cmd"),
				),
			},
		},
	})
}

const testAccDataSourceSettingDockerPluginInstallScript_basic = `
data "huaweicloud_hss_setting_docker_plugin_install_script" "test" {
  plugin                = "opa-docker-authz"
  operate_type          = "install"
  enterprise_project_id = "0"
}
`
