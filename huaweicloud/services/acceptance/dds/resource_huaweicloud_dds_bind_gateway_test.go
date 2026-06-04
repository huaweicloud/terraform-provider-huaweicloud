package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dds"
)

func getResourceBindGatewayFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dds", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DDS client: %s", err)
	}

	return dds.GetGatewayInfo(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

func TestAccResourceBindGateway_basic(t *testing.T) {
	var (
		rName  = "huaweicloud_dds_bind_gateway.test"
		name   = acceptance.RandomAccResourceName()
		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceBindGatewayFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccBindGateway_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "huaweicloud_dds_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "nat_gateway_id", "huaweicloud_nat_gateway.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "public_ip_id", "huaweicloud_vpc_eip.test", "id"),
					resource.TestCheckResourceAttr(rName, "external_service_port", "8080"),
					resource.TestCheckResourceAttrSet(rName, "node_id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccBindGatewayImportStateFunc(rName),
			},
		},
	})
}

func testAccBindGateway_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_nat_gateway" "test" {
  name      = "%[2]s"
  spec      = "1"
  vpc_id    = huaweicloud_vpc.test.id
  subnet_id = huaweicloud_vpc_subnet.test.id
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

resource "huaweicloud_dds_instance" "test" {
  name              = "%[2]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "Terraform@123"
  mode              = "ReplicaSet"

  datastore {
    type           = "DDS-Community"
    version        = "4.0"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "replica"
    storage   = "ULTRAHIGH"
    num       = 3
    size      = 10
    spec_code = "dds.mongodb.s6.large.2.repset"
  }
}

data "huaweicloud_dds_instances" "test" {
  name = huaweicloud_dds_instance.test.name
}

locals {
  node_id = try([for v in flatten(data.huaweicloud_dds_instances.test.instances[*].groups[*].nodes) : v if v.role == "Primary"][0].id, "")
}
`, common.TestBaseNetwork(name), name)
}

func testAccBindGateway_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dds_bind_gateway" "test"{
  instance_id           = huaweicloud_dds_instance.test.id
  node_id               = local.node_id
  nat_gateway_id        = huaweicloud_nat_gateway.test.id
  public_ip_id          = huaweicloud_vpc_eip.test.id
  external_service_port = 8080
}
`, testAccBindGateway_base(name))
}

func testAccBindGatewayImportStateFunc(rName string) resource.ImportStateIdFunc {
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
