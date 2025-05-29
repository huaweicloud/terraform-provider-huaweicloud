package asm

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

func getResourceAsmMeshFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		getMeshHttpUrl = "v1/{project_id}/meshes/{mesh_id}"
		getMeshProduct = "asm"
	)
	getMeshClient, err := cfg.NewServiceClient(getMeshProduct, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Mesh Client: %s", err)
	}

	getMeshPath := getMeshClient.Endpoint + getMeshHttpUrl
	getMeshPath = strings.ReplaceAll(getMeshPath, "{project_id}", getMeshClient.ProjectID)
	getMeshPath = strings.ReplaceAll(getMeshPath, "{mesh_id}", state.Primary.ID)

	getPotectionRulesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getMeshResp, err := getMeshClient.Request("GET", getMeshPath, &getPotectionRulesOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving ASM mesh: %s", err)
	}

	return utils.FlattenResponse(getMeshResp)
}

func TestAccResourceAsmMesh_basic(t *testing.T) {
	resourceName := "huaweicloud_asm_mesh.test"
	randName := acceptance.RandomAccResourceNameWithDash()
	var obj interface{}

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getResourceAsmMeshFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceResourceAsmMesh_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "type", "InCluster"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.18.7-r1"),
					resource.TestCheckResourceAttr(resourceName, "status", "Running"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "tags", "extend_params"},
			},
		},
	})
}

func testResourceResourceAsmMesh_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%[1]s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"

  //dns is required for cce node installing
  primary_dns   = "100.125.1.250"
  secondary_dns = "100.125.21.250"
  vpc_id        = huaweicloud_vpc.test.id
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 4
  memory_size       = 8
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[1]s"
  flavor_id              = "cce.s1.small"
  cluster_version        = "v1.28"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  container_network_cidr = "172.16.0.0/24"
}

resource "huaweicloud_cce_node" "test" {
  cluster_id        = huaweicloud_cce_cluster.test.id
  name              = "%[1]s"
  flavor_id         = data.huaweicloud_compute_flavors.test.ids[0]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  password          = "Test@1234"
  os                = "EulerOS 2.9"

  root_volume {
    size       = 40
    volumetype = "SSD"
  }

  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, name)
}

func testResourceResourceAsmMesh_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_asm_mesh" "test" {
  name    = "%[2]s"
  type    = "InCluster"
  version = "1.18.7-r1"

  tags = {
    foo = "bar"
    key = "value"
  }

  extend_params {
    clusters {
      cluster_id = huaweicloud_cce_cluster.test.id
      installation {
        nodes {
          field_selector {
            key      = "UID"
            operator = "In"
            values   = [
                huaweicloud_cce_node.test.id
            ]
          }
        }
      }
    }
  }
}
`, testResourceResourceAsmMesh_base(name), name)
}
