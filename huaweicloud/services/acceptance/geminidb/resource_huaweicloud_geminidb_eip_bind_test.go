package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/geminidb"
)

func getResourceEipBindFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("geminidb", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating GeminiDB client: %s", err)
	}

	return geminidb.GetEipBindInfo(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

func TestAccEipBind_basic(t *testing.T) {
	var obj interface{}
	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_geminidb_eip_bind.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getResourceEipBindFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEipBind_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "huaweicloud_geminidb_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "public_ip_id", "huaweicloud_vpc_eip.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "public_ip", "huaweicloud_vpc_eip.test", "address"),
					resource.TestCheckResourceAttrSet(rName, "node_id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccEipBindImportStateFunc(rName),
			},
		},
	})
}

func testAccEipBind_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus             = 2
  engine            = "redis"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_vpc_eip" "test" {
  name = "%[2]s"

  publicip {
    type       = "5_sbgp"
    ip_version = 4
  }

  bandwidth {
    name        = "%[2]s"
    share_type  = "PER"
    size        = 1
    charge_mode = "bandwidth"
  }
}

resource "huaweicloud_geminidb_instance" "test" {
  name              = "%[2]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "test_1234"
  mode              = "Cluster"

  datastore {
    type           = "redis"
    version        = "5.0"
    storage_engine = "rocksDB"
  }

  flavor {
    num       = "3"
    size      = "16"
    storage   = "ULTRAHIGH"
    spec_code = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].name
  }
}
`, common.TestBaseNetwork(name), name)
}

func testAccEipBind_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_geminidb_eip_bind" "test" {
  instance_id  = huaweicloud_geminidb_instance.test.id
  node_id      = huaweicloud_geminidb_instance.test.groups[0].nodes[0].id
  public_ip    = huaweicloud_vpc_eip.test.address
  public_ip_id = huaweicloud_vpc_eip.test.id
}
`, testAccEipBind_base(name))
}

func testAccEipBindImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, nodeId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		instanceId = rs.Primary.Attributes["instance_id"]
		nodeId = rs.Primary.ID

		if instanceId == "" || nodeId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>', but got '%s/%s'",
				instanceId, nodeId)
		}

		return fmt.Sprintf("%s/%s", instanceId, nodeId), nil
	}
}
