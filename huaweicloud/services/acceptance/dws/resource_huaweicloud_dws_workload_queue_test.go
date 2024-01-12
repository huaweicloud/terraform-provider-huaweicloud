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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getWorkloadQueueResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v2/{project_id}/clusters/{cluster_id}/workload/queues"
		product = "dws"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DWS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", state.Primary.Attributes["cluster_id"])
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DWS workload queue: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DWS workload queue: %s", err)
	}

	expression := fmt.Sprintf("workload_queue_name_list[?@=='%s']|[0]", state.Primary.ID)
	resp := utils.PathSearch(expression, getRespBody, nil)
	if resp == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return resp, nil
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
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkloadQueue_basic(name, "Admin_user@123"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "configuration.#", "4"),
					resource.TestCheckResourceAttrPair(resourceName, "cluster_id", "huaweicloud_dws_cluster.test", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testWorkloadQueueImportState(resourceName),
				ImportStateVerifyIgnore: []string{
					"configuration", "logical_cluster_name",
				},
			},
		},
	})
}

func testAccWorkloadQueue_base(name, password string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dws_cluster" "test" {
  name              = "%[2]s"
  node_type         = "dwsk2.xlarge"
  number_of_node    = 3
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  user_name         = "admin_user"
  user_pwd          = "%[3]s"
}
`, common.TestBaseNetwork(name), name, password)
}

func testAccWorkloadQueue_basic(name, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dws_workload_queue" "test" {
  cluster_id = huaweicloud_dws_cluster.test.id
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
`, testAccWorkloadQueue_base(name, password), name)
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

		return fmt.Sprintf("%s/%s", clusterId, id), nil
	}
}
