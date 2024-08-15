package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWorkloadPlans_basic(t *testing.T) {
	var (
		rName      = acceptance.RandomAccResourceName()
		dataSource = "data.huaweicloud_dws_workload_plans.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byLogicalClusterName   = "data.huaweicloud_dws_workload_plans.filter_by_logical_cluster_name"
		dcByLogicalClusterName = acceptance.InitDataSourceCheck(byLogicalClusterName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testDataSourceDwsWorkloadPlans_clusterNotFound(),
				ExpectError: regexp.MustCompile("Cluster does not exist or has been deleted"),
			},
			{
				Config: testDataSourceWorkloadPlans_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "plans.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_exist_plan", "true"),
					dcByLogicalClusterName.CheckResourceExists(),
					resource.TestCheckOutput("is_userful_logical_cluster_name", "true"),
					resource.TestCheckResourceAttrSet(byLogicalClusterName, "plans.0.id"),
					resource.TestCheckResourceAttrSet(byLogicalClusterName, "plans.0.name"),
					resource.TestCheckResourceAttrSet(byLogicalClusterName, "plans.0.cluster_id"),
					resource.TestCheckResourceAttrSet(byLogicalClusterName, "plans.0.status"),
					resource.TestCheckResourceAttrSet(byLogicalClusterName, "plans.0.current_stage_name"),
				),
			},
		},
	})
}

func testDataSourceDwsWorkloadPlans_clusterNotFound() string {
	randUUID, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dws_workload_plans" "test" {
  cluster_id = "%s"
}
`, randUUID)
}

func testDataSourceDwsWorkloadPlans_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_dws_logical_cluster_rings" "test" {
  cluster_id = "%[1]s"
}

locals {
  ring_hosts = [for v in data.huaweicloud_dws_logical_cluster_rings.test.cluster_rings : v.ring_hosts if v.is_available]
}

resource "huaweicloud_dws_logical_cluster" "test" {
  cluster_id           = "%[1]s"
  logical_cluster_name = "%[2]s"

  cluster_rings {
    dynamic "ring_hosts" {
      for_each = local.ring_hosts[0]
      content {
        host_name = ring_hosts.value.host_name
        back_ip   = ring_hosts.value.back_ip
        cpu_cores = ring_hosts.value.cpu_cores
        memory    = ring_hosts.value.memory
        disk_size = ring_hosts.value.disk_size
      }
    }
  }

  lifecycle {
    ignore_changes = [ 
      cluster_rings,
    ]
  }
}

resource "huaweicloud_dws_workload_plan" "test" {
  cluster_id           = "%[1]s"
  name                 = "%[2]s"
  logical_cluster_name = huaweicloud_dws_logical_cluster.test.logical_cluster_name

  lifecycle {
    ignore_changes = [ 
      logical_cluster_name,
    ]
  }
}

resource "huaweicloud_dws_workload_queue" "test" {
  cluster_id           = "%[1]s"
  name                 = "%[2]s"
  logical_cluster_name = huaweicloud_dws_logical_cluster.test.logical_cluster_name

  configuration {
    resource_name  = "cpu_limit"
    resource_value = 10
  }
  configuration {
    resource_name  = "memory"
    resource_value = 10
  }
  configuration {
    resource_name  = "tablespace"
    resource_value = -1
  }
  configuration {
    resource_name  = "activestatements"
    resource_value = -1
  }
}

resource "huaweicloud_dws_workload_plan_stage" "test" {
  count      = 2
  cluster_id = "%[1]s"
  plan_id    = huaweicloud_dws_workload_plan.test.id
  name       = "%[2]s${count.index}"
  start_time = "01:00:00"
  end_time   = "02:00:00"
  month      = format("%%s", count.index + 1)
  day        = "1"

  queues {
    name = huaweicloud_dws_workload_queue.test.name

    configuration {
      resource_name  = "cpu"
      resource_value = 1
    }
    configuration {
      resource_name  = "cpu_limit"
      resource_value = 0
    }
    configuration {
      resource_name  = "memory"
      resource_value = 0
    }
    configuration {
      resource_name  = "concurrency"
      resource_value = 10
    }
    configuration {
      resource_name  = "shortQueryConcurrencyNum"
      resource_value = -1
    }
  }
}


resource "huaweicloud_dws_workload_plan_execution" "test" {
  depends_on = [
    huaweicloud_dws_workload_plan_stage.test,
  ]

  cluster_id = "%[1]s"
  plan_id    = huaweicloud_dws_workload_plan.test.id
  stage_id   = huaweicloud_dws_workload_plan_stage.test[0].id
}
`, acceptance.HW_DWS_CLUSTER_ID, name)
}

func testDataSourceWorkloadPlans_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dws_workload_plans" "test" {
  depends_on = [
    huaweicloud_dws_workload_plan_execution.test
  ]

  cluster_id = "%[2]s"
}

output "is_exist_plan" {
  value  = contains(data.huaweicloud_dws_workload_plans.test.plans[*].id, huaweicloud_dws_workload_plan.test.id)
}

data "huaweicloud_dws_workload_plans" "filter_by_logical_cluster_name" {
  depends_on = [
    huaweicloud_dws_workload_plan_execution.test
  ]

  cluster_id           = "%[2]s"
  logical_cluster_name = huaweicloud_dws_logical_cluster.test.logical_cluster_name
}

# The logical_cluster_name parameter is queried through exact match.
output "is_userful_logical_cluster_name" {
  value = length(data.huaweicloud_dws_workload_plans.filter_by_logical_cluster_name.plans) == 1
}
 `, testDataSourceDwsWorkloadPlans_base(name), acceptance.HW_DWS_CLUSTER_ID)
}
