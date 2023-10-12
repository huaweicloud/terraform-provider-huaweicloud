package mrs

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

func getScalingPolicyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getScalingPolicy: Query the scaling policy
	var (
		getScalingPolicyHttpUrl = "v2/{project_id}/autoscaling-policy/{cluster_id}"
		getScalingPolicyProduct = "mrs"
	)
	getScalingPolicyClient, err := cfg.NewServiceClient(getScalingPolicyProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating MRS client: %s", err)
	}

	getScalingPolicyPath := getScalingPolicyClient.Endpoint + getScalingPolicyHttpUrl
	getScalingPolicyPath = strings.ReplaceAll(getScalingPolicyPath, "{project_id}", getScalingPolicyClient.ProjectID)
	getScalingPolicyPath = strings.ReplaceAll(getScalingPolicyPath, "{cluster_id}", state.Primary.Attributes["cluster_id"])

	getScalingPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getScalingPolicyResp, err := getScalingPolicyClient.Request("GET", getScalingPolicyPath, &getScalingPolicyOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving scaling policy: %s", err)
	}

	getScalingPolicyRespBody, err := utils.FlattenResponse(getScalingPolicyResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving scaling policy: %s", err)
	}

	jsonPath := fmt.Sprintf("[?node_group_name =='%s']|[0]", state.Primary.Attributes["node_group"])
	scalingPolicy := utils.PathSearch(jsonPath, getScalingPolicyRespBody, nil)
	if scalingPolicy == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	if !utils.PathSearch("auto_scaling_policy.auto_scaling_enable", scalingPolicy, false).(bool) {
		return nil, golangsdk.ErrDefault404{}
	}

	return scalingPolicy, nil
}

func TestAccScalingPolicy_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_mapreduce_scaling_policy.test"
	pwd := acceptance.RandomPassword()

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getScalingPolicyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMrsBootstrapScript(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testScalingPolicy_basic(name, pwd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "cluster_id", "huaweicloud_mapreduce_cluster.test", "id"),
					resource.TestCheckResourceAttr(rName, "node_group", "task_node_streaming_group"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_enable", "true"),
					resource.TestCheckResourceAttr(rName, "min_capacity", "1"),
					resource.TestCheckResourceAttr(rName, "max_capacity", "10"),
					resource.TestCheckResourceAttr(rName, "resources_plans.0.period_type", "daily"),
					resource.TestCheckResourceAttr(rName, "resources_plans.0.start_time", "01:00"),
					resource.TestCheckResourceAttr(rName, "resources_plans.0.end_time", "03:00"),
					resource.TestCheckResourceAttr(rName, "resources_plans.0.min_capacity", "1"),
					resource.TestCheckResourceAttr(rName, "resources_plans.0.max_capacity", "10"),
					resource.TestCheckResourceAttr(rName, "rules.0.name", "default-expand-1"),
					resource.TestCheckResourceAttr(rName, "rules.0.adjustment_type", "scale_out"),
					resource.TestCheckResourceAttr(rName, "rules.0.cool_down_minutes", "20"),
					resource.TestCheckResourceAttr(rName, "rules.0.scaling_adjustment", "1"),
					resource.TestCheckResourceAttr(rName, "rules.0.trigger.0.metric_name", "StormSlotUsed"),
					resource.TestCheckResourceAttr(rName, "rules.0.trigger.0.metric_value", "75"),
					resource.TestCheckResourceAttr(rName, "rules.0.trigger.0.comparison_operator", "GT"),
					resource.TestCheckResourceAttr(rName, "rules.0.trigger.0.evaluation_periods", "2"),
					resource.TestCheckResourceAttr(rName, "rules.1.name", "default-shrink-1"),
					resource.TestCheckResourceAttr(rName, "rules.1.adjustment_type", "scale_in"),
					resource.TestCheckResourceAttr(rName, "rules.1.cool_down_minutes", "20"),
					resource.TestCheckResourceAttr(rName, "rules.1.scaling_adjustment", "1"),
					resource.TestCheckResourceAttr(rName, "rules.1.trigger.0.metric_name", "StormSlotUsed"),
					resource.TestCheckResourceAttr(rName, "rules.1.trigger.0.metric_value", "25"),
					resource.TestCheckResourceAttr(rName, "rules.1.trigger.0.comparison_operator", "LT"),
					resource.TestCheckResourceAttr(rName, "rules.1.trigger.0.evaluation_periods", "2"),
					resource.TestCheckResourceAttr(rName, "exec_scripts.0.name", "script_show_dir"),
					resource.TestCheckResourceAttr(rName, "exec_scripts.0.active_master", "false"),
					resource.TestCheckResourceAttr(rName, "exec_scripts.0.fail_action", "continue"),
					resource.TestCheckResourceAttr(rName, "exec_scripts.0.action_stage", "before_scale_out"),
				),
			},
			{
				Config: testScalingPolicy_basic_update(name, pwd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "cluster_id", "huaweicloud_mapreduce_cluster.test", "id"),
					resource.TestCheckResourceAttr(rName, "node_group", "task_node_streaming_group"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_enable", "true"),
					resource.TestCheckResourceAttr(rName, "min_capacity", "1"),
					resource.TestCheckResourceAttr(rName, "max_capacity", "11"),
					resource.TestCheckResourceAttr(rName, "resources_plans.0.period_type", "daily"),
					resource.TestCheckResourceAttr(rName, "resources_plans.0.start_time", "02:00"),
					resource.TestCheckResourceAttr(rName, "resources_plans.0.end_time", "04:00"),
					resource.TestCheckResourceAttr(rName, "resources_plans.0.min_capacity", "1"),
					resource.TestCheckResourceAttr(rName, "resources_plans.0.max_capacity", "11"),
					resource.TestCheckResourceAttr(rName, "rules.#", "1"),
					resource.TestCheckResourceAttr(rName, "rules.0.name", "default-shrink-1"),
					resource.TestCheckResourceAttr(rName, "rules.0.adjustment_type", "scale_in"),
					resource.TestCheckResourceAttr(rName, "rules.0.cool_down_minutes", "30"),
					resource.TestCheckResourceAttr(rName, "rules.0.scaling_adjustment", "1"),
					resource.TestCheckResourceAttr(rName, "rules.0.trigger.0.metric_name", "StormSlotUsed"),
					resource.TestCheckResourceAttr(rName, "rules.0.trigger.0.metric_value", "45"),
					resource.TestCheckResourceAttr(rName, "rules.0.trigger.0.comparison_operator", "LT"),
					resource.TestCheckResourceAttr(rName, "rules.0.trigger.0.evaluation_periods", "2"),
					resource.TestCheckResourceAttr(rName, "exec_scripts.#", "0"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testScalingPolicyImportState(rName),
			},
		},
	})
}

func testScalingPolicy_basic(name, pwd string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_mapreduce_scaling_policy" "test" {
  cluster_id = huaweicloud_mapreduce_cluster.test.id
  node_group = "task_node_streaming_group"

  auto_scaling_enable = true
  min_capacity        = 1
  max_capacity        = 10

  resources_plans {
    period_type  = "daily"
    start_time   = "01:00"
    end_time     = "03:00"
    min_capacity = 1
    max_capacity = 10
  }

  rules {
    name               = "default-expand-1"
    adjustment_type    = "scale_out"
    cool_down_minutes  = 20
    scaling_adjustment = 1
    trigger {
      metric_name         = "StormSlotUsed"
      metric_value        = "75"
      comparison_operator = "GT"
      evaluation_periods  = 2
    }
  }

  rules {
    name               = "default-shrink-1"
    adjustment_type    = "scale_in"
    cool_down_minutes  = 20
    scaling_adjustment = 1
    trigger {
      metric_name         = "StormSlotUsed"
      metric_value        = "25"
      comparison_operator = "LT"
      evaluation_periods  = 2
    }
  }

  exec_scripts {
    name          = "script_show_dir"
    uri           = "%s"
    parameters    = ""
    nodes         = ["task_node_streaming_group"]
    active_master = false
    fail_action   = "continue"
    action_stage  = "before_scale_out"
  }
}

`, testAccMrsMapReduceClusterConfig_basic(name, pwd), acceptance.HW_MAPREDUCE_BOOTSTRAP_SCRIPT)
}

func testScalingPolicy_basic_update(name, pwd string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_mapreduce_scaling_policy" "test" {
  cluster_id = huaweicloud_mapreduce_cluster.test.id
  node_group = "task_node_streaming_group"

  auto_scaling_enable = true
  min_capacity        = 1
  max_capacity        = 11

  resources_plans {
    period_type  = "daily"
    start_time   = "02:00"
    end_time     = "04:00"
    min_capacity = 1
    max_capacity = 11
  }

  rules {
    name               = "default-shrink-1"
    adjustment_type    = "scale_in"
    cool_down_minutes  = 30
    scaling_adjustment = 1
    trigger {
      metric_name         = "StormSlotUsed"
      metric_value        = "45"
      comparison_operator = "LT"
      evaluation_periods  = 2
    }
  }
}

`, testAccMrsMapReduceClusterConfig_basic(name, pwd))
}

func testScalingPolicyImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["cluster_id"] == "" {
			return "", fmt.Errorf("attribute (cluster_id) ofresource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["node_group"] == "" {
			return "", fmt.Errorf("attribute (node_group) of resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.Attributes["cluster_id"] + "/" +
			rs.Primary.Attributes["node_group"], nil
	}
}
