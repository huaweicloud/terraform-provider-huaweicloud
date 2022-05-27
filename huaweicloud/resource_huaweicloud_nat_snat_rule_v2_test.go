package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/hw_snatrules"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccNatSnatRule_basic(t *testing.T) {
	randSuffix := acctest.RandString(5)
	resourceName := "huaweicloud_nat_snat_rule.snat_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNatV2SnatRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatV2SnatRule_basic(randSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatV2GatewayExists("huaweicloud_nat_gateway.nat_1"),
					testAccCheckNatV2SnatRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform acc test"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(resourceName, "nat_gateway_id",
						"huaweicloud_nat_gateway.nat_1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "floating_ip_address",
						"huaweicloud_vpc_eip.eips.0", "address"),
				),
			},
			{
				Config: testAccNatV2SnatRule_update(randSuffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "updated by terraform acc test"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(resourceName, "nat_gateway_id",
						"huaweicloud_nat_gateway.nat_1", "id"),
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

func testAccCheckNatV2SnatRuleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	natClient, err := config.NatGatewayClient(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud nat client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_nat_snat_rule" {
			continue
		}

		_, err := hw_snatrules.Get(natClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Snat rule still exists")
		}
	}

	return nil
}

func testAccCheckNatV2SnatRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		natClient, err := config.NatGatewayClient(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud nat client: %s", err)
		}

		found, err := hw_snatrules.Get(natClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Snat rule not found")
		}

		return nil
	}
}

func testAccNatV2SnatRule_base(suffix string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_eip" "eips" {
  count = 2

  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_nat_gateway" "nat_1" {
  name        = "nat-gateway-basic-%s"
  description = "created by terraform acc test"
  spec        = "1"
  vpc_id      = huaweicloud_vpc.vpc_1.id
  subnet_id   = huaweicloud_vpc_subnet.subnet_1.id
}
`, testAccNatPreCondition(suffix), suffix)
}

func testAccNatV2SnatRule_basic(suffix string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_nat_snat_rule" "snat_1" {
  nat_gateway_id = huaweicloud_nat_gateway.nat_1.id
  subnet_id      = huaweicloud_vpc_subnet.subnet_1.id
  floating_ip_id = huaweicloud_vpc_eip.eips.0.id
  description    = "created by terraform acc test"
}
`, testAccNatV2SnatRule_base(suffix))
}

func testAccNatV2SnatRule_update(suffix string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_nat_snat_rule" "snat_1" {
  nat_gateway_id = huaweicloud_nat_gateway.nat_1.id
  subnet_id      = huaweicloud_vpc_subnet.subnet_1.id
  floating_ip_id = join(",", huaweicloud_vpc_eip.eips.*.id)
  description    = "updated by terraform acc test"
}
`, testAccNatV2SnatRule_base(suffix))
}
