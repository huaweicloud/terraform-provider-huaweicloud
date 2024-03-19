package dws

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccWorkloadQueuesDataSource_basic(t *testing.T) {
	resourceName := "data.huaweicloud_dws_workload_queues.test"
	dc := acceptance.InitDataSourceCheck(resourceName)
	name := acceptance.RandomAccResourceName()
	// The cluster password requires a minimum length of 12 characters, and the string 'gap' is used to fill in the gap.
	password := acceptance.RandomPassword() + "gap"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWorkloadQueuesDataSourceBasic(name, password),
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

func testAccWorkloadQueuesDataSourceBasic(name, password string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dws_workload_queues" "test" {
  cluster_id = huaweicloud_dws_cluster.test.id

  depends_on = [huaweicloud_dws_workload_queue.test]
}

data "huaweicloud_dws_workload_queues" "name_filter" {
  cluster_id = huaweicloud_dws_cluster.test.id
  name       = huaweicloud_dws_workload_queue.test.name
}

data "huaweicloud_dws_workload_queues" "name_not_exist_filter" {
  cluster_id = huaweicloud_dws_cluster.test.id
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
`, testAccWorkloadQueue_basic(name, password))
}
