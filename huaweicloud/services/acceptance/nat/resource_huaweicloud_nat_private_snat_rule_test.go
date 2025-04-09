package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/nat"
)

func getPrivateSnatRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return nil, fmt.Errorf("error creating NAT v3 client: %s", err)
	}

	return nat.GetPrivateSnatRule(client, state.Primary.ID)
}

func TestAccPrivateSnatRule_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_nat_private_snat_rule.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPrivateSnatRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateSnatRule_basic_step_1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "gateway_id", "huaweicloud_nat_private_gateway.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(rName, "transit_ip_ids.#", "2"),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttr(rName, "transit_ip_associations.#", "2"),
					resource.TestCheckResourceAttrSet(rName, "enterprise_project_id"),
				),
			},
			{
				Config: testAccPrivateSnatRule_basic_step_2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "transit_ip_ids.#", "1"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"transit_ip_ids"},
			},
		},
	})
}

func testAccPrivateSnatRule_transitIpConfig(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "transit_ip_used" {
  name = "%[1]s-transit-ip"
  cidr = "172.16.0.0/16"
}

resource "huaweicloud_vpc_subnet" "transit_ip_used" {
  vpc_id     = huaweicloud_vpc.transit_ip_used.id
  name       = "%[1]s-transit-ip"
  cidr       = cidrsubnet(huaweicloud_vpc.transit_ip_used.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.transit_ip_used.cidr, 4, 1), 1)
}

resource "huaweicloud_nat_private_transit_ip" "test1" {
  subnet_id             = huaweicloud_vpc_subnet.transit_ip_used.id
  enterprise_project_id = "0"
}

resource "huaweicloud_nat_private_transit_ip" "test2" {
  subnet_id             = huaweicloud_vpc_subnet.transit_ip_used.id
  enterprise_project_id = "0"
}
`, name)
}

func testAccPrivateSnatRule_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_private_gateway" "test" {
  subnet_id             = huaweicloud_vpc_subnet.test.id
  name                  = "%[2]s"
  enterprise_project_id = "0"
}
`, common.TestBaseNetwork(name), name)
}

func testAccPrivateSnatRule_basic_step_1(name string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "huaweicloud_nat_private_snat_rule" "test" {
  gateway_id     = huaweicloud_nat_private_gateway.test.id
  description    = "Created by acc test"
  transit_ip_ids = [huaweicloud_nat_private_transit_ip.test1.id,huaweicloud_nat_private_transit_ip.test2.id]
  subnet_id      = huaweicloud_vpc_subnet.test.id
}
`, testAccPrivateSnatRule_base(name), testAccPrivateSnatRule_transitIpConfig(name))
}

func testAccPrivateSnatRule_basic_step_2(name string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "huaweicloud_nat_private_snat_rule" "test" {
  gateway_id     = huaweicloud_nat_private_gateway.test.id
  transit_ip_ids = [huaweicloud_nat_private_transit_ip.test2.id]
  subnet_id      = huaweicloud_vpc_subnet.test.id
}
`, testAccPrivateSnatRule_base(name), testAccPrivateSnatRule_transitIpConfig(name))
}

func TestAccPrivateSnatRule_cidr(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_nat_private_snat_rule.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPrivateSnatRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateSnatRule_cidr_step_1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "gateway_id", "huaweicloud_nat_private_gateway.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttrPair(rName, "transit_ip_id", "huaweicloud_nat_private_transit_ip.test1", "id"),
					resource.TestCheckResourceAttrPair(rName, "cidr", "huaweicloud_vpc_subnet.test", "cidr"),
					resource.TestCheckResourceAttrSet(rName, "transit_ip_address"),
					resource.TestCheckResourceAttrSet(rName, "enterprise_project_id"),
				),
			},
			{
				Config: testAccPrivateSnatRule_cidr_step_2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttrPair(rName, "transit_ip_id", "huaweicloud_nat_private_transit_ip.test2", "id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPrivateSnatRule_cidr_step_1(name string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "huaweicloud_nat_private_snat_rule" "test" {
  gateway_id    = huaweicloud_nat_private_gateway.test.id
  description   = "Created by acc test"
  transit_ip_id = huaweicloud_nat_private_transit_ip.test1.id
  cidr          = huaweicloud_vpc_subnet.test.cidr
}
`, testAccPrivateSnatRule_base(name), testAccPrivateSnatRule_transitIpConfig(name))
}

func testAccPrivateSnatRule_cidr_step_2(name string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "huaweicloud_nat_private_snat_rule" "test" {
  gateway_id    = huaweicloud_nat_private_gateway.test.id
  transit_ip_id = huaweicloud_nat_private_transit_ip.test2.id
  cidr          = huaweicloud_vpc_subnet.test.cidr
}
`, testAccPrivateSnatRule_base(name), testAccPrivateSnatRule_transitIpConfig(name))
}
