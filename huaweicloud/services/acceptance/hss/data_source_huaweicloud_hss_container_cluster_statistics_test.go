package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceContainerClusterStatistics_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_container_cluster_statistics.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceContainerClusterStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "risk_cluster_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "app_vul_cluster_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "unscan_cluster_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "all_cluster_num"),
				),
			},
		},
	})
}

func testAccDataSourceContainerClusterStatistics_basic() string {
	return `
data "huaweicloud_hss_container_cluster_statistics" "test" {
  enterprise_project_id = "0"
}
`
}
