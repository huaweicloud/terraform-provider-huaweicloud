package cceautopilot

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCceAutopilotClusters_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_cce_autopilot_clusters.basic"
	dataSource2 := "data.huaweicloud_cce_autopilot_clusters.filter_by_detail"
	dataSource3 := "data.huaweicloud_cce_autopilot_clusters.filter_by_status"
	rName := acceptance.RandomAccResourceNameWithDash()
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)
	dc3 := acceptance.InitDataSourceCheck(dataSource3)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCceAutopilotClusters_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					dc3.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_detail_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCceAutopilotClusters_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cce_autopilot_clusters" "basic" {
  depends_on = [huaweicloud_cce_autopilot_cluster.test]
}

data "huaweicloud_cce_autopilot_clusters" "filter_by_detail" {
  detail = "true"

  depends_on = [huaweicloud_cce_autopilot_cluster.test]
}

data "huaweicloud_cce_autopilot_clusters" "filter_by_status" {
  status = "Available"

  depends_on = [huaweicloud_cce_autopilot_cluster.test]
}

locals {
  detail_result = [for v in data.huaweicloud_cce_autopilot_clusters.filter_by_detail.clusters[*].annotations.installedAddonInstances : v != ""]
  status_result = [for v in data.huaweicloud_cce_autopilot_clusters.filter_by_status.clusters[*].status[0].phase : v == "Available"]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_cce_autopilot_clusters.basic.clusters) > 0
}

output "is_detail_filter_useful" {
  value = alltrue(local.detail_result) && length(local.detail_result) > 0
}

output "is_status_filter_useful" {
  value = alltrue(local.status_result) && length(local.status_result) > 0
}
`, testAccCluster_basic(name))
}
