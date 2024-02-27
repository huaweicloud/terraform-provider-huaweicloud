package dws

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getWorkLoadPlanExecutionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v2/{project_id}/clusters/{cluster_id}/workload/plans/{plan_id}"
		product = "dws"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DWS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", state.Primary.Attributes["cluster_id"])
	getPath = strings.ReplaceAll(getPath, "{plan_id}", state.Primary.Attributes["plan_id"])
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error querying DWS workload plan: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	// When calling the API, there is a situation where the plan ID does not exist but still returns a 200 status code.
	// When the workload plan is successfully started, the value of the status attribute of the workload plan is 1.
	resCode := utils.PathSearch("workload_res_code", getRespBody, float64(0)).(float64)
	plan := utils.PathSearch("workload_plan", getRespBody, nil)
	status := utils.PathSearch("workload_plan.status", getRespBody, float64(0)).(float64)
	if resCode != 0 || status != 1 {
		return nil, golangsdk.ErrDefault404{}
	}

	return plan, nil
}

func TestAccResourceWorkLoadPlanExecution_basic(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_dws_workload_plan_execution.test"
		name         = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getWorkLoadPlanExecutionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkLoadPlanExecution_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "cluster_id", "huaweicloud_dws_cluster.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "plan_id", "huaweicloud_dws_workload_plan.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "stage_id", "huaweicloud_dws_workload_plan_stage.test1", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testAccWorkLoadPlanExecution_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "cluster_id", "huaweicloud_dws_cluster.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "plan_id", "huaweicloud_dws_workload_plan.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "stage_id", "huaweicloud_dws_workload_plan_stage.test2", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccWorkLoadPlanExecution_base(name string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "huaweicloud_dws_workload_plan_stage" "test1" {
  cluster_id = huaweicloud_dws_cluster.test.id
  plan_id    = huaweicloud_dws_workload_plan.test.id
  name       = "%[3]s_1"
  start_time = "07:00:00"
  end_time   = "08:00:00"
  month      = "1"
  day        = "1"

  queues {
    name = huaweicloud_dws_workload_queue.test.name

    configuration {
      resource_name  = "cpu"
      resource_value = 10
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

resource "huaweicloud_dws_workload_plan_stage" "test2" {
  cluster_id = huaweicloud_dws_cluster.test.id
  plan_id    = huaweicloud_dws_workload_plan.test.id
  name       = "%[3]s_2"
  start_time = "07:00:00"
  end_time   = "08:00:00"
  month      = "2"
  day        = "2"

  queues {
    name = huaweicloud_dws_workload_queue.test.name

    configuration {
      resource_name  = "cpu"
      resource_value = 10
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
`, testAccWorkLoadPlan_basic(name), testAccWorkloadPlanStage_base(name), name)
}

func testAccWorkLoadPlanExecution_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dws_workload_plan_execution" "test" {
  depends_on = [
    huaweicloud_dws_workload_plan_stage.test1,
    huaweicloud_dws_workload_plan_stage.test2,
  ]

  cluster_id = huaweicloud_dws_cluster.test.id
  plan_id    = huaweicloud_dws_workload_plan.test.id
  stage_id   = huaweicloud_dws_workload_plan_stage.test1.id
}
`, testAccWorkLoadPlanExecution_base(name))
}

func testAccWorkLoadPlanExecution_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dws_workload_plan_execution" "test" {
  depends_on = [
    huaweicloud_dws_workload_plan_stage.test1,
    huaweicloud_dws_workload_plan_stage.test2,
  ]

  cluster_id = huaweicloud_dws_cluster.test.id
  plan_id    = huaweicloud_dws_workload_plan.test.id
  stage_id   = huaweicloud_dws_workload_plan_stage.test2.id
}
`, testAccWorkLoadPlanExecution_base(name))
}
