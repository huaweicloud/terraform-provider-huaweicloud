package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceContainerNodeStatistics_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_container_node_statistics.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceContainerNodeStatistics_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "unprotected_num"),
					resource.TestCheckResourceAttrSet(dataSource, "protected_num"),
					resource.TestCheckResourceAttrSet(dataSource, "protected_num_on_demand"),
					resource.TestCheckResourceAttrSet(dataSource, "protected_num_packet_cycle"),
					resource.TestCheckResourceAttrSet(dataSource, "cluster_node_not_installed_num"),
					resource.TestCheckResourceAttrSet(dataSource, "not_cluster_node_not_installed_num"),
				),
			},
		},
	})
}

const testDataSourceContainerNodeStatistics_basic string = `
data "huaweicloud_hss_container_node_statistics" "test" {
  enterprise_project_id = "0"
}
`
