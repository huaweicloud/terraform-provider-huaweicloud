package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSettingDictionaries_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_setting_dictionaries.test"
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
				Config: testAccDataSourceSettingDictionaries_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.code"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.value"),

					resource.TestCheckOutput("code_filter_useful", "true"),
					resource.TestCheckOutput("scene_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceSettingDictionaries_basic = `
data "huaweicloud_hss_setting_dictionaries" "test" {
  group_code = "featureSwitch"
}

locals {
  code = data.huaweicloud_hss_setting_dictionaries.test.data_list[0].code
}

data "huaweicloud_hss_setting_dictionaries" "code_filter" {
  group_code = "featureSwitch"
  code       = local.code
}

output "code_filter_useful" {
  value = length(data.huaweicloud_hss_setting_dictionaries.code_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_setting_dictionaries.code_filter.data_list[*].code : v == local.code]
  )
}

data "huaweicloud_hss_setting_dictionaries" "scene_filter" {
  group_code = "featureSwitch"
  scene      = "hws"
}

output "scene_filter_useful" {
  value = length(data.huaweicloud_hss_setting_dictionaries.scene_filter.data_list) > 0
}
`
