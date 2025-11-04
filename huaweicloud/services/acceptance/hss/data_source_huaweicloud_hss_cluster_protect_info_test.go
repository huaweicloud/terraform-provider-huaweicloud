package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceClusterProtectInfo_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_cluster_protect_info.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires the preparation of a CCE cluster under the default enterprise project.
			acceptance.TestAccPreCheckHSSClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceClusterProtectInfo_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.cluster_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.cluster_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.cluster_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.protect_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.policy_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.cluster_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.cluster_type"),
				),
			},
		},
	})
}

func testAccDataSourceClusterProtectInfo_basic() string {
	return `
data "huaweicloud_hss_cluster_protect_info" "test" {
  enterprise_project_id = "0"
}
`
}
