package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceConfigs_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_configs.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConfigs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.config_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.config_value"),
				),
			},
		},
	})
}

func testAccDataSourceConfigs_basic() string {
	return `
data "huaweicloud_hss_configs" "test" {
  config_name_list = ["password_min_len"]
}
`
}
