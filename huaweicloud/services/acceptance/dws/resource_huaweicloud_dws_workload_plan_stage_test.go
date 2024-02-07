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

func getWorkloadPlanStageResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v2/{project_id}/clusters/{cluster_id}/workload/plans/{plan_id}/stages/{stage_id}"
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
	getPath = strings.ReplaceAll(getPath, "{stage_id}", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DWS workload plan stage: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	stage := utils.PathSearch("workload_plan_stage", respBody, nil)
	if stage == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return stage, nil
}

func TestAccResourceWorkLoadPlanStage_basic(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_dws_workload_plan_stage.test"
		name         = acceptance.RandomAccResourceName()
	)
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getWorkloadPlanStageResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkloadPlanStage_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "start_time", "07:08:00"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testWorkloadPlanStageImportState(resourceName),
			},
		},
	})
}

func testAccWorkloadPlanStage_base(name string) string {
	return fmt.Sprintf(`

resource "huaweicloud_dws_workload_queue" "test" {
  cluster_id = huaweicloud_dws_cluster.test.id
  name       = "%s"

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
`, name)
}

func testAccWorkloadPlanStage_basic(name string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_dws_workload_plan_stage" "test" {
  cluster_id = huaweicloud_dws_cluster.test.id
  plan_id    = huaweicloud_dws_workload_plan.test.id
  name       = "%s"
  start_time = "07:08:00"
  end_time   = "00:00:00"

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
`, testAccWorkLoadPlan_basic(name), testAccWorkloadPlanStage_base(name), name)
}

func testWorkloadPlanStageImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		clusterId := rs.Primary.Attributes["cluster_id"]
		planId := rs.Primary.Attributes["plan_id"]
		name := rs.Primary.Attributes["name"]
		if clusterId == "" || planId == "" || name == "" {
			return "", fmt.Errorf("the workload plan stage name, plan ID or cluster ID is missing")
		}

		return fmt.Sprintf("%s/%s/%s", clusterId, planId, name), nil
	}
}
