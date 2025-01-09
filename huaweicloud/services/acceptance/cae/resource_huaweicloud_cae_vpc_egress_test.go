package cae

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cae"
)

func getVpcEgressFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("cae", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CAE client: %s", err)
	}

	environmentId := state.Primary.Attributes["environment_id"]
	return cae.GetVpcEgressById(client, environmentId, state.Primary.ID)
}

func TestAccVpcEgress_basic(t *testing.T) {
	var (
		obj interface{}

		name = acceptance.RandomAccResourceName()

		rName = "huaweicloud_cae_vpc_egress.test"
		rc    = acceptance.InitResourceCheck(
			rName,
			&obj,
			getVpcEgressFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCaeEnvironment(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEgress_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "environment_id", acceptance.HW_CAE_ENVIRONMENT_ID),
					resource.TestCheckResourceAttrPair(rName, "route_table_id", "huaweicloud_vpc_route_table.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "cidr", "huaweicloud_vpc_route_table.test", "route.0.destination"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccVpcEgressImportStateFunc(rName),
			},
		},
	})
}

func testAccVpcEgress_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_cae_environments" "test" {
  environment_id = "%[1]s"
}

locals {
  vpc_id = data.huaweicloud_cae_environments.test.environments[0].annotations.vpc_id
}

resource "huaweicloud_vpc" "test" {
  name = "%[2]s"
  cidr = "192.168.10.0/24"
}

resource "huaweicloud_vpc_peering_connection" "test" {
  name        = "%[2]s"
  vpc_id      = local.vpc_id
  peer_vpc_id = huaweicloud_vpc.test.id
}

resource "huaweicloud_vpc_route_table" "test" {
  name   = "%[2]s"
  vpc_id = local.vpc_id

  subnets = [
    data.huaweicloud_cae_environments.test.environments[0].annotations.subnet_id,
  ]

  route {
    destination = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1)
    type        = "peering"
    nexthop     = huaweicloud_vpc_peering_connection.test.id
  }
  route {
    destination = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 2)
    type        = "peering"
    nexthop     = huaweicloud_vpc_peering_connection.test.id
  }
}
`, acceptance.HW_CAE_ENVIRONMENT_ID, name)
}

func testAccVpcEgress_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cae_vpc_egress" "test" {
  environment_id = "%[2]s"
  route_table_id = huaweicloud_vpc_route_table.test.id
  cidr           = tolist(huaweicloud_vpc_route_table.test.route)[0].destination
}
 `, testAccVpcEgress_base(name), acceptance.HW_CAE_ENVIRONMENT_ID)
}

func testAccVpcEgressImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		var (
			environmentId = rs.Primary.Attributes["environment_id"]
			routeTableId  = rs.Primary.Attributes["route_table_id"]
			cidr          = rs.Primary.Attributes["cidr"]
		)

		if environmentId == "" || routeTableId == "" || cidr == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<environment_id>,<route_table_id>,<cidr>', but got '%s,%s,%s'",
				environmentId, routeTableId, cidr)
		}

		return fmt.Sprintf("%s,%s,%s", environmentId, routeTableId, cidr), nil
	}
}
