package cae

import (
	"fmt"
	"log"
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

	return cae.GetVpcEgressById(
		client,
		state.Primary.Attributes["environment_id"],
		state.Primary.ID,
		state.Primary.Attributes["enterprise_project_id"],
	)
}

func TestAccVpcEgress_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		vpcEgress interface{}
		rName     = "huaweicloud_cae_vpc_egress.test.0"
		rc        = acceptance.InitResourceCheck(rName, &vpcEgress, getVpcEgressFunc)

		nameWithSpecifiedEps = "huaweicloud_cae_vpc_egress.test.1"
		rcWithSpecifiedEps   = acceptance.InitResourceCheck(nameWithSpecifiedEps, &vpcEgress, getVpcEgressFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCaeEnvironmentIds(t, 2)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rc.CheckResourceDestroy(),
			rcWithSpecifiedEps.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEgress_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "environment_id"),
					resource.TestCheckResourceAttrPair(rName, "route_table_id", "huaweicloud_vpc_route_table.test.0", "id"),
					resource.TestCheckResourceAttrPair(rName, "cidr", "huaweicloud_vpc_route_table.test.0", "route.0.destination"),
					resource.TestCheckNoResourceAttr(rName, "enterprise_project_id"),
					rcWithSpecifiedEps.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(nameWithSpecifiedEps, "environment_id"),
					resource.TestCheckResourceAttrPair(nameWithSpecifiedEps, "route_table_id", "huaweicloud_vpc_route_table.test.1", "id"),
					resource.TestCheckResourceAttrPair(nameWithSpecifiedEps, "cidr", "huaweicloud_vpc_route_table.test.1", "route.0.destination"),
					resource.TestCheckResourceAttrPair(nameWithSpecifiedEps, "enterprise_project_id",
						"data.huaweicloud_cae_environments.test.1", "environments.0.annotations.enterprise_project_id"),
				),
			},
			{
				ResourceName:      "huaweicloud_cae_vpc_egress.test[0]",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccVpcEgressImportStateFunc(rName, true),
			},
			{
				ResourceName:      "huaweicloud_cae_vpc_egress.test[1]",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccVpcEgressImportStateFunc(nameWithSpecifiedEps, false),
			},
		},
	})
}

func testAccVpcEgress_base(name string) string {
	return fmt.Sprintf(`
locals {
  env_ids     = split(",", "%[1]s")
}

data "huaweicloud_cae_environments" "test" {
  count          = 2
  environment_id = local.env_ids[count.index]
}

locals {
  env_vpc_ids = data.huaweicloud_cae_environments.test[*].environments[0].annotations.vpc_id
}

resource "huaweicloud_vpc" "peering" {
  count = 2

  name = "%[2]s${count.index}"
  cidr = cidrsubnet("172.16.0.0/16", 4, count.index)
}

resource "huaweicloud_vpc_peering_connection" "test" {
  count = 2

  name        = "%[2]s${count.index}"
  vpc_id      = local.env_vpc_ids[count.index]
  peer_vpc_id = huaweicloud_vpc.peering[count.index].id
}

resource "huaweicloud_vpc_route_table" "test" {
  count = 2

  name   = "%[2]s${count.index}"
  vpc_id = local.env_vpc_ids[count.index]

  subnets = [
    data.huaweicloud_cae_environments.test[count.index].environments[0].annotations.subnet_id,
  ]

  route {
    destination = cidrsubnet(huaweicloud_vpc.peering[count.index].cidr, 4, 1)
    type        = "peering"
    nexthop     = huaweicloud_vpc_peering_connection.test[count.index].id
  }
  route {
    destination = cidrsubnet(huaweicloud_vpc.peering[count.index].cidr, 4, 2)
    type        = "peering"
    nexthop     = huaweicloud_vpc_peering_connection.test[count.index].id
  }
}
`, acceptance.HW_CAE_ENVIRONMENT_IDS, name)
}

func testAccVpcEgress_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cae_vpc_egress" "test" {
  count = 2

  environment_id        = local.env_ids[count.index]
  route_table_id        = huaweicloud_vpc_route_table.test[count.index].id
  cidr                  = tolist(huaweicloud_vpc_route_table.test[count.index].route)[0].destination
  enterprise_project_id = count.index == 1 ? try(data.huaweicloud_cae_environments.test[1].environments[0].annotations.enterprise_project_id,
  null) : null
}
 `, testAccVpcEgress_base(name))
}

func testAccVpcEgressImportStateFunc(name string, isDefaultEps bool) resource.ImportStateIdFunc {
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
			return "", fmt.Errorf("some import IDs are missing, want '<environment_id>,<route_table_id>,<cidr>', "+
				"or '<environment_id>,<route_table_id>,<cidr>,<enterprise_project_id>', but got '%s,%s,%s'",
				environmentId, routeTableId, cidr)
		}

		if isDefaultEps {
			log.Printf("[DEBUG] isDefaultEps: %s", fmt.Sprintf("%s,%s,%s", environmentId, routeTableId, cidr))
			return fmt.Sprintf("%s,%s,%s", environmentId, routeTableId, cidr), nil
		}

		epsId := rs.Primary.Attributes["enterprise_project_id"]
		if epsId == "" {
			return "", fmt.Errorf("enterprise_project_id is missing, want '<environment_id>,<route_table_id>,<cidr>,<enterprise_project_id>', "+
				"but got '%s,%s,%s,%s'", environmentId, routeTableId, cidr, epsId)
		}

		return fmt.Sprintf("%s,%s,%s,%s", environmentId, routeTableId, cidr, epsId), nil
	}
}
