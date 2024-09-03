package dws

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWorkloadQueues_basic(t *testing.T) {
	resourceName := "data.huaweicloud_dws_workload_queues.test"
	dc := acceptance.InitDataSourceCheck(resourceName)
	name := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceWorkloadQueues_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "queues.#"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("no_filter_is_useful", "true"),
					resource.TestCheckOutput("name_not_exist_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceWorkloadQueues_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dws_workload_queues" "test" {
  cluster_id = "%[2]s"

  depends_on = [huaweicloud_dws_workload_queue.test]
}

data "huaweicloud_dws_workload_queues" "name_filter" {
  cluster_id = "%[2]s"
  name       = huaweicloud_dws_workload_queue.test.name
}

data "huaweicloud_dws_workload_queues" "name_not_exist_filter" {
  cluster_id = "%[2]s"
  name       = "name_not_exist"

  depends_on = [huaweicloud_dws_workload_queue.test]
}

locals {
  no_filter = data.huaweicloud_dws_workload_queues.test.queues[*]

  name_filter = [for v in data.huaweicloud_dws_workload_queues.name_filter.queues[*] : 
  strcontains(v.name, huaweicloud_dws_workload_queue.test.name)]

  name_not_exist_filter = data.huaweicloud_dws_workload_queues.name_not_exist_filter.queues[*]
}

output "no_filter_is_useful" {
  value = length(local.name_filter) > 0
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter) && length(local.name_filter) > 0
}

output "name_not_exist_filter_is_useful" {
  value = length(local.name_not_exist_filter) == 0
}
`, testAccWorkloadQueue_basic(name), acceptance.HW_DWS_CLUSTER_ID)
}

func TestAccDataSourceWorkloadQueues_logicalCluster(t *testing.T) {
	resourceName := "data.huaweicloud_dws_workload_queues.filter_by_logical_cluster_name"
	dc := acceptance.InitDataSourceCheck(resourceName)
	name := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsLogicalModeClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceWorkloadQueues_logicalCluster(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "queues.#"),
					resource.TestCheckOutput("is_userful_logical_cluster_name", "true"),
				),
			},
		},
	})
}

func testAccDataSourceWorkloadQueues_logicalCluster(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dws_workload_queues" "filter_by_logical_cluster_name" {
  depends_on = [huaweicloud_dws_workload_queue.test]

  cluster_id           = "%[2]s"
  logical_cluster_name = "%[3]s"
}

locals {
  fliter_result = [for v in data.huaweicloud_dws_workload_queues.filter_by_logical_cluster_name.queues[*].logical_cluster_name : v == "%[3]s"]
}

output "is_userful_logical_cluster_name" {
  value = alltrue(local.fliter_result) && length(local.fliter_result) > 0
}
`, testAccResourceWorkloadQueue_logicalClusterName(name),
		acceptance.HW_DWS_LOGICAL_MODE_CLUSTER_ID,
		acceptance.HW_DWS_LOGICAL_CLUSTER_NAME)
}
