package cceautopilot

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

func getAutopilotClusterFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var (
		getClusterHttpUrl = "autopilot/v3/projects/{project_id}/clusters/{cluster_id}"
		getClusterProduct = "cce"
	)
	getClusterClient, err := cfg.NewServiceClient(getClusterProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CCE Client: %s", err)
	}

	getClusterPath := getClusterClient.Endpoint + getClusterHttpUrl
	getClusterPath = strings.ReplaceAll(getClusterPath, "{project_id}", getClusterClient.ProjectID)
	getClusterPath = strings.ReplaceAll(getClusterPath, "{cluster_id}", state.Primary.ID)

	getClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getClusterResp, err := getClusterClient.Request("GET", getClusterPath, &getClusterOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CCE autopolit cluster: %s", err)
	}

	return utils.FlattenResponse(getClusterResp)
}

func TestAccAutopilotCluster_basic(t *testing.T) {
	var (
		cluster      interface{}
		resourceName = "huaweicloud_cce_autopilot_cluster.test"
		rName        = acceptance.RandomAccResourceNameWithDash()

		rc = acceptance.InitResourceCheck(
			resourceName,
			&cluster,
			getAutopilotClusterFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "flavor", "cce.autopilot.cluster"),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(resourceName, "container_network.0.mode", "eni"),
					resource.TestCheckResourceAttrPair(resourceName, "host_network.0.vpc", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "host_network.0.subnet", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "eni_network.0.subnets.0.subnet_id",
						"huaweicloud_vpc_subnet.test", "ipv4_subnet_id"),
					resource.TestCheckResourceAttr(resourceName, "status.0.phase", "Available"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCluster_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_autopilot_cluster" "test" {
  name        = "%[2]s"
  flavor      = "cce.autopilot.cluster"
  description = "created by terraform"

  host_network {
    vpc    = huaweicloud_vpc.test.id
    subnet = huaweicloud_vpc_subnet.test.id
  }

  container_network {
    mode = "eni"
  }

  eni_network {
    subnets {
      subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id
    }
  }
}
`, common.TestVpc(rName), rName)
}
