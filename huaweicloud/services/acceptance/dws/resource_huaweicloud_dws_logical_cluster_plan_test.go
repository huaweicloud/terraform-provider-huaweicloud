package dws

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dws"
)

func getLogicalClusterPlanResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return nil, fmt.Errorf("error creating DWS client: %s", err)
	}

	return dws.GetLogicalClusterPlanById(client, state.Primary.Attributes["cluster_id"], state.Primary.ID)
}

func TestAccLogicalClusterPlan_basic(t *testing.T) {
	var (
		obj interface{}

		rName     = "huaweicloud_dws_logical_cluster_plan.test"
		rc        = acceptance.InitResourceCheck(rName, &obj, getLogicalClusterPlanResourceFunc)
		name      = acceptance.RandomAccResourceName()
		today8am  = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 8, 0, 0, 0, time.Now().Location())
		startTime = today8am.UnixMilli()
		endTime   = today8am.AddDate(0, 1, 0).UnixMilli()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			// Step 1: Verify invalid cluster_id returns error
			{
				Config:      testAccLogicalClusterPlan_invalidCluster(name),
				ExpectError: regexp.MustCompile("Cluster does not exist or has been deleted."),
			},
			// Step 2: Create resource successfully
			{
				Config: testAccLogicalClusterPlan_basic(name, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_DWS_CLUSTER_ID),
					resource.TestCheckResourceAttr(rName, "plan_type", "periodicity"),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
					resource.TestCheckResourceAttr(rName, "main_logical_cluster", "v3_logical"),
					resource.TestCheckResourceAttrSet(rName, "start_time"),
					resource.TestCheckResourceAttrSet(rName, "end_time"),
					resource.TestCheckResourceAttr(rName, "actions.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(rName, "actions.*", map[string]string{
						"type":     "create",
						"strategy": "0 00 0 ? * 1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(rName, "actions.*", map[string]string{
						"type":     "delete",
						"strategy": "0 04 0 ? * 6",
					}),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			// Step 3: Import test (enabled=false state)
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"actions", "main_logical_cluster"},
				ImportStateIdFunc:       testAccLogicalClusterPlanImportState(rName),
			},
			// Step 4: Enable plan
			{
				Config: testAccLogicalClusterPlan_enable(name, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			// Step 5: Update actions attribute
			{
				Config: testAccLogicalClusterPlan_updateAction(name, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "actions.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(rName, "actions.*", map[string]string{
						"type":     "create",
						"strategy": "0 00 0 ? * 1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(rName, "actions.*", map[string]string{
						"type":     "delete",
						"strategy": "0 06 0 ? * 6",
					}),
				),
			},
			// Step 6: Disable plan
			{
				Config: testAccLogicalClusterPlan_disable(name, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
					resource.TestCheckResourceAttr(rName, "status", "disabled"),
				),
			},
			// Step 7: Import test (enabled=true state after enable/disable operations)
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"actions", "main_logical_cluster"},
				ImportStateIdFunc:       testAccLogicalClusterPlanImportState(rName),
			},
		},
	})
}

func testAccLogicalClusterPlan_invalidCluster(name string) string {
	randomUUID, _ := uuid.NewRandom()
	return fmt.Sprintf(`
variable "logical_cluster_plan_actions" {
  type    = list(object({
    type     = string
    strategy = string
  }))
  default = [
    {
      type     = "create"
      strategy = "0 00 0 ? * 1"
    },
    {
      type     = "delete"
      strategy = "0 04 0 ? * 6"
    }
  ]
}

resource "huaweicloud_dws_logical_cluster_plan" "invalid" {
  cluster_id           = "%[1]s"
  plan_type            = "periodicity"
  node_num             = 3
  logical_cluster_name = "%[2]s"

  dynamic "actions" {
    for_each = var.logical_cluster_plan_actions

    content {
      type     = actions.value.type
      strategy = actions.value.strategy
    }
  }
}
`, randomUUID.String(), name)
}

func testAccLogicalClusterPlan_basic(name string, startTime, endTime int64) string {
	return fmt.Sprintf(`
variable "logical_cluster_plan_actions" {
  type    = list(object({
    type     = string
    strategy = string
  }))
  default = [
    {
      type     = "create"
      strategy = "0 00 0 ? * 1"
    },
    {
      type     = "delete"
      strategy = "0 04 0 ? * 6"
    }
  ]
}

resource "huaweicloud_dws_logical_cluster_plan" "test" {
  cluster_id           = "%[1]s"
  plan_type            = "periodicity"
  node_num             = 3
  logical_cluster_name = "%[2]s"
  main_logical_cluster = "v3_logical"
  start_time           = "%[3]d"
  end_time             = "%[4]d"
  enabled              = false

  dynamic "actions" {
    for_each = var.logical_cluster_plan_actions

    content {
      type     = actions.value.type
      strategy = actions.value.strategy
    }
  }
}
`, acceptance.HW_DWS_CLUSTER_ID, name, startTime, endTime)
}

func testAccLogicalClusterPlan_enable(name string, startTime, endTime int64) string {
	return fmt.Sprintf(`
variable "logical_cluster_plan_actions" {
  type    = list(object({
    type     = string
    strategy = string
  }))
  default = [
    {
      type     = "create"
      strategy = "0 00 0 ? * 1"
    },
    {
      type     = "delete"
      strategy = "0 04 0 ? * 5"
    }
  ]
}

resource "huaweicloud_dws_logical_cluster_plan" "test" {
  cluster_id           = "%[1]s"
  plan_type            = "periodicity"
  node_num             = 3
  logical_cluster_name = "%[2]s"
  main_logical_cluster = "v3_logical"
  start_time           = "%[3]d"
  end_time             = "%[4]d"
  enabled              = true

  dynamic "actions" {
    for_each = var.logical_cluster_plan_actions

    content {
      type     = actions.value.type
      strategy = actions.value.strategy
    }
  }
}
`, acceptance.HW_DWS_CLUSTER_ID, name, startTime, endTime)
}

func testAccLogicalClusterPlan_updateAction(name string, startTime, endTime int64) string {
	return fmt.Sprintf(`
variable "logical_cluster_plan_actions" {
  type    = list(object({
    type     = string
    strategy = string
  }))
  default = [
    {
      type     = "create"
      strategy = "0 00 0 ? * 1"
    },
    {
      type     = "delete"
      strategy = "0 06 0 ? * 6"
    }
  ]
}

resource "huaweicloud_dws_logical_cluster_plan" "test" {
  cluster_id           = "%[1]s"
  plan_type            = "periodicity"
  node_num             = 3
  logical_cluster_name = "%[2]s"
  main_logical_cluster = "v3_logical"
  start_time           = "%[3]d"
  end_time             = "%[4]d"
  enabled              = true

  dynamic "actions" {
    for_each = var.logical_cluster_plan_actions

    content {
      type     = actions.value.type
      strategy = actions.value.strategy
    }
  }
}
`, acceptance.HW_DWS_CLUSTER_ID, name, startTime, endTime)
}

func testAccLogicalClusterPlan_disable(name string, startTime, endTime int64) string {
	return fmt.Sprintf(`
variable "logical_cluster_plan_actions" {
  type    = list(object({
    type     = string
    strategy = string
  }))
  default = [
    {
      type     = "create"
      strategy = "0 00 0 ? * 1"
    },
    {
      type     = "delete"
      strategy = "0 06 0 ? * 6"
    }
  ]
}

resource "huaweicloud_dws_logical_cluster_plan" "test" {
  cluster_id           = "%[1]s"
  plan_type            = "periodicity"
  node_num             = 3
  logical_cluster_name = "%[2]s"
  main_logical_cluster = "v3_logical"
  start_time           = "%[3]d"
  end_time             = "%[4]d"
  enabled              = false

  dynamic "actions" {
    for_each = var.logical_cluster_plan_actions

    content {
      type     = actions.value.type
      strategy = actions.value.strategy
    }
  }
}
`, acceptance.HW_DWS_CLUSTER_ID, name, startTime, endTime)
}

func testAccLogicalClusterPlanImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		return acceptance.HW_DWS_CLUSTER_ID + "/" + rs.Primary.ID, nil
	}
}
