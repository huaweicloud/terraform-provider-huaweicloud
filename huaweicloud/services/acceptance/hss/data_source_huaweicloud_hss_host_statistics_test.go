package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHostStatistics_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_host_statistics.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceHostStatistics_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "version_enterprise_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "version_advanced_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "host_group_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "asset_value_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "asset_value_list.0.value_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "protected_num"),
				),
			},
		},
	})
}

const testDataSourceHostStatistics_basic = `
data "huaweicloud_hss_host_statistics" "test" {
  enterprise_project_id = "0"
}
`
