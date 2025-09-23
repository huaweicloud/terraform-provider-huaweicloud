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
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
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
				Config: testAccCluster_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "alias", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key_update", "value_update"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"updated_at"},
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

  tags = {
    "foo" = "bar"
    "key" = "value"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_autopilot_cluster" "test" {
  name        = "%[2]s"
  alias       = "%[2]s-update"
  flavor      = "cce.autopilot.cluster"
  description = ""

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

  tags = {
    "foo"        = "bar_update"
    "key_update" = "value_update"
  }
}
`, common.TestVpc(rName), rName)
}

func TestAccAutopilotCluster_withEip(t *testing.T) {
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
				Config: testAccCluster_withEip(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					// if the length of endpoints is 3, there is an external endpoint
					// which means the EIP is bind to the cluster
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status.0.phase", "Available"),
					resource.TestCheckResourceAttr(resourceName, "status.0.endpoints.#", "3"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccCluster_withEipUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status.0.endpoints.#", "2"),
				),
			},
		},
	})
}

func testAccCluster_withEip(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_cce_autopilot_cluster" "test" {
  name   = "%[2]s"
  flavor = "cce.autopilot.cluster"
  eip_id = huaweicloud_vpc_eip.test.id

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

func testAccCluster_withEipUpdate(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_cce_autopilot_cluster" "test" {
  name   = "%[2]s"
  flavor = "cce.autopilot.cluster"
  eip_id = ""

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
