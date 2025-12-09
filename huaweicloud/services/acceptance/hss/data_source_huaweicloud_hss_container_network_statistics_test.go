package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceContainerNetworkStatistics_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_container_network_statistics.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceContainerNetworkStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "protected_cluster_total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cluster_total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "namespace_total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_policy_total_num"),
				),
			},
		},
	})
}

func testDataSourceContainerNetworkStatistics_basic() string {
	return `
data "huaweicloud_hss_container_network_statistics" "test" {
  enterprise_project_id = "0"
}
`
}
