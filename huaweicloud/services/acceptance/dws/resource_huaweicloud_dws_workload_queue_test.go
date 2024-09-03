package dws

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dws"
)

func getWorkloadQueueResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "dws"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DWS client: %s", err)
	}

	clusterId := state.Primary.Attributes["cluster_id"]
	logicalClusterName := state.Primary.Attributes["logical_cluster_name"]
	return dws.GetWorkloadQueueByName(client, clusterId, state.Primary.ID, logicalClusterName)
}

func TestAccResourceWorkloadQueue_basic(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_dws_workload_queue.test"
		name         = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getWorkloadQueueResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkloadQueue_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "configuration.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "cluster_id", acceptance.HW_DWS_CLUSTER_ID),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testWorkloadQueueImportState(resourceName),
				ImportStateVerifyIgnore: []string{
					"configuration",
				},
			},
		},
	})
}

func testAccWorkloadQueue_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_workload_queue" "test" {
  cluster_id = "%[1]s"
  name       = "%[2]s"

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
`, acceptance.HW_DWS_CLUSTER_ID, name)
}

func TestAccResourceWorkloadQueue_logicalClusterName(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_dws_workload_queue.test"
		name         = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getWorkloadQueueResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsLogicalModeClusterId(t)
			acceptance.TestAccPreCheckDwsLogicalClusterName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkloadQueue_logicalClusterName(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cluster_id", acceptance.HW_DWS_LOGICAL_MODE_CLUSTER_ID),
					resource.TestCheckResourceAttr(resourceName, "logical_cluster_name", acceptance.HW_DWS_LOGICAL_CLUSTER_NAME),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testWorkloadQueueImportState(resourceName),
				ImportStateVerifyIgnore: []string{
					"configuration",
				},
			},
		},
	})
}

func testAccResourceWorkloadQueue_logicalClusterName(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_workload_queue" "test" {
  cluster_id           = "%[1]s"
  name                 = "%[2]s"
  logical_cluster_name = "%[3]s"

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
`, acceptance.HW_DWS_LOGICAL_MODE_CLUSTER_ID, name, acceptance.HW_DWS_LOGICAL_CLUSTER_NAME)
}

func testWorkloadQueueImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		clusterId := rs.Primary.Attributes["cluster_id"]
		id := rs.Primary.ID
		if clusterId == "" || id == "" {
			return "", fmt.Errorf("the workload queue is not exist or related cluster ID is missing")
		}

		logicalClusterName := rs.Primary.Attributes["logical_cluster_name"]
		if logicalClusterName != "" {
			return fmt.Sprintf("%s/%s/%s", clusterId, id, logicalClusterName), nil
		}
		return fmt.Sprintf("%s/%s", clusterId, id), nil
	}
}
