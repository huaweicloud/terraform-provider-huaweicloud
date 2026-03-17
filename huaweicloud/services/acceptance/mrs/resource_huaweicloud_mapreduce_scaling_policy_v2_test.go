package mrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/mrs"
)

func getScalingPolicyV2ResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("mrs", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating MRS client: %s", err)
	}

	return mrs.GetScalingPolicyV2(client, state.Primary.Attributes["cluster_id"],
		state.Primary.Attributes["node_group_name"], state.Primary.Attributes["resource_pool_name"])
}

// Before running acceptance test, please ensure that the node name provided is that of the task node.
func TestAccScalingPolicyV2_basic(t *testing.T) {
	var (
		obj interface{}

		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_mapreduce_scaling_policy_v2.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getScalingPolicyV2ResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMrsClusterID(t)
			acceptance.TestAccPreCheckMrsClusterNodeGroupName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testScalingPolicyV2_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_MRS_CLUSTER_ID),
					resource.TestCheckResourceAttr(rName, "node_group_name", acceptance.HW_MRS_CLUSTER_NODE_GROUP_NAME),
					resource.TestCheckResourceAttr(rName, "resource_pool_name", "default"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.auto_scaling_enable", "false"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.min_capacity", "1"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.max_capacity", "10"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.resources_plans.#", "2"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.resources_plans.1.period_type", "daily"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.resources_plans.1.start_time", "05:00"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.resources_plans.1.end_time", "06:00"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.resources_plans.1.min_capacity", "0"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.resources_plans.1.max_capacity", "1"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.resources_plans.1.effective_days.#", "2"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.resources_plans.1.effective_days.1", "SATURDAY"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.rules.#", "2"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.rules.0.name", "default-expand-1"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.rules.0.adjustment_type", "scale_out"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.rules.0.description", "Created_by_TF_script"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.rules.1.name", "default-shrink-1"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.rules.1.adjustment_type", "scale_in"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.rules.1.cool_down_minutes", "30"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.rules.1.scaling_adjustment", "2"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.rules.1.trigger.0.metric_name", "YARNAppRunning"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.rules.1.trigger.0.metric_value", "25"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.rules.1.trigger.0.comparison_operator", "LT"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.rules.1.trigger.0.evaluation_periods", "2"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.tags.%", "1"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.tags.foo", "bar"),
				),
			},
			{
				Config: testScalingPolicyV2_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.auto_scaling_enable", "false"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.min_capacity", "0"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.max_capacity", "0"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.resources_plans.#", "0"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.rules.#", "1"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.rules.0.name", "default-shrink"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.rules.0.adjustment_type", "scale_out"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.rules.0.cool_down_minutes", "10"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.rules.0.scaling_adjustment", "2"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.rules.0.description", ""),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.rules.0.trigger.0.metric_name", "YARNAppPending"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.rules.0.trigger.0.metric_value", "80"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.rules.0.trigger.0.comparison_operator", "LTOE"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.rules.0.trigger.0.evaluation_periods", "5"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.tags.%", "0"),
				),
			},
			{
				Config: testScalingPolicyV2_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.auto_scaling_enable", "true"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.min_capacity", "2"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.max_capacity", "3"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.resources_plans.#", "1"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.resources_plans.0.period_type", "daily"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.resources_plans.0.start_time", "03:00"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.resources_plans.0.end_time", "03:30"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.resources_plans.0.effective_days.#", "1"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.resources_plans.0.effective_days.0", "MONDAY"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.rules.#", "0"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.tags.%", "2"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.tags.owner", "terraform"),
					resource.TestCheckResourceAttr(rName, "auto_scaling_policy.0.tags.foo", "barr"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testScalingPolicyV2ImportState(rName),
			},
		},
	})
}

func testScalingPolicyV2_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_mapreduce_scaling_policy_v2" "test" {
  cluster_id         = "%[1]s"
  node_group_name    = "%[2]s"
  resource_pool_name = "default"

  auto_scaling_policy {
    min_capacity = 1
    max_capacity = 10

    resources_plans {
      period_type  = "daily"
      start_time   = "01:00"
      end_time     = "03:00"
      min_capacity = 1
      max_capacity = 2
    }
    resources_plans {
      period_type    = "daily"
      start_time     = "05:00"
      end_time       = "06:00"
      min_capacity   = 0
      max_capacity   = 1
      effective_days = ["SUNDAY", "SATURDAY"]
    }

    rules {
      name               = "default-expand-1"
      adjustment_type    = "scale_out"
      cool_down_minutes  = 20
      scaling_adjustment = 1
      # Only letters, digits, underscores (_), and hyphens (-) are allowed.
      description        = "Created_by_TF_script"

      trigger {
        metric_name         = "YARNAppRunning"
        metric_value        = "75"
        comparison_operator = "GT"
        evaluation_periods  = 2
      }
    }
    rules {
      name               = "default-shrink-1"
      adjustment_type    = "scale_in"
      cool_down_minutes  = 30
      scaling_adjustment = 2

      trigger {
        metric_name         = "YARNAppRunning"
        metric_value        = "25"
        comparison_operator = "LT"
        evaluation_periods  = 2
      }
    }

    tags = {
      foo = "bar"
    }
  }
}
`, acceptance.HW_MRS_CLUSTER_ID, acceptance.HW_MRS_CLUSTER_NODE_GROUP_NAME, name)
}

func testScalingPolicyV2_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_mapreduce_scaling_policy_v2" "test" {
  cluster_id         = "%[1]s"
  node_group_name    = "%[2]s"
  resource_pool_name = "default"

  auto_scaling_policy {
    rules {
      name               = "default-shrink"
      adjustment_type    = "scale_out"
      cool_down_minutes  = 10
      scaling_adjustment = 2

      trigger {
        metric_name         = "YARNAppPending"
        metric_value        = "80"
        comparison_operator = "LTOE"
        evaluation_periods  = 5
      }
    }
  }
}
`, acceptance.HW_MRS_CLUSTER_ID, acceptance.HW_MRS_CLUSTER_NODE_GROUP_NAME, name)
}

func testScalingPolicyV2_basic_step3(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_mapreduce_scaling_policy_v2" "test" {
  cluster_id         = "%[1]s"
  node_group_name    = "%[2]s"
  resource_pool_name = "default"

  auto_scaling_policy {
    min_capacity        = 2
    max_capacity        = 3
    auto_scaling_enable = true

    resources_plans {
      period_type    = "daily"
      start_time     = "03:00"
      end_time       = "03:30"
      effective_days = ["MONDAY"]
    }

    tags = {
      owner = "terraform"
      foo   = "barr"
    }
  }
}
`, acceptance.HW_MRS_CLUSTER_ID, acceptance.HW_MRS_CLUSTER_NODE_GROUP_NAME, name)
}

func testScalingPolicyV2ImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		clusterId := rs.Primary.Attributes["cluster_id"]
		nodeGroupName := rs.Primary.Attributes["node_group_name"]
		resourcePoolName := rs.Primary.Attributes["resource_pool_name"]
		if clusterId == "" || nodeGroupName == "" || resourcePoolName == "" {
			return "", fmt.Errorf("some import IDs is valid, want `<cluster_id>/<node_group_name>/<resource_pool_name>`, "+
				"but got '%s/%s/%s'", clusterId, nodeGroupName, resourcePoolName)
		}

		return fmt.Sprintf("%s/%s/%s", clusterId, nodeGroupName, resourcePoolName), nil
	}
}
