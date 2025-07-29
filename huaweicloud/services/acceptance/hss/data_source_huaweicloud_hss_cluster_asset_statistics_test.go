package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceClusterAssetStatistics_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_cluster_asset_statistics.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceClusterAssetStatistics_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "cluster_num"),
					resource.TestCheckResourceAttrSet(dataSource, "work_load_num"),
					resource.TestCheckResourceAttrSet(dataSource, "service_num"),
					resource.TestCheckResourceAttrSet(dataSource, "pod_num"),
				),
			},
		},
	})
}

const testDataSourceClusterAssetStatistics_basic = `
data "huaweicloud_hss_cluster_asset_statistics" "test" {
  enterprise_project_id = "0"
}
`
