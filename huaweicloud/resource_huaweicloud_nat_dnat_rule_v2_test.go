// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file at
//     https://www.github.com/huaweicloud/magic-modules
//
// ----------------------------------------------------------------------------

package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccNatDnat_basic(t *testing.T) {
	randSuffix := acctest.RandString(5)
	resourceName := "huaweicloud_nat_dnat_rule.dnat"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNatDnatDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatV2DnatRule_basic(randSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatDnatExists(),
					resource.TestCheckResourceAttr(resourceName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
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

func TestAccNatDnat_protocol(t *testing.T) {
	randSuffix := acctest.RandString(5)
	resourceName := "huaweicloud_nat_dnat_rule.dnat"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNatDnatDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatV2DnatRule_protocol(randSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatDnatExists(),
					resource.TestCheckResourceAttr(resourceName, "protocol", "any"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
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

func testAccCheckNatDnatDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.NatGatewayClient(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating sdk client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_nat_dnat_rule" {
			continue
		}

		url, err := replaceVarsForTest(rs, "dnat_rules/{id}")
		if err != nil {
			return err
		}
		url = client.ServiceURL(url)

		_, err = client.Get(
			url, nil,
			&golangsdk.RequestOpts{MoreHeaders: map[string]string{"Accept": "application/json"}})
		if err == nil {
			return fmtp.Errorf("huaweicloud dnat rule still exists at %s", url)
		}
	}

	return nil
}

func testAccCheckNatDnatExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*config.Config)
		client, err := config.NatGatewayClient(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating sdk client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources["huaweicloud_nat_dnat_rule.dnat"]
		if !ok {
			return fmtp.Errorf("Error checking huaweicloud_nat_dnat_rule.dnat exist, err=not found huaweicloud_nat_dnat_rule.dnat")
		}

		url, err := replaceVarsForTest(rs, "dnat_rules/{id}")
		if err != nil {
			return fmtp.Errorf("Error checking huaweicloud_nat_dnat_rule.dnat exist, err=building url failed: %s", err)
		}
		url = client.ServiceURL(url)

		_, err = client.Get(
			url, nil,
			&golangsdk.RequestOpts{MoreHeaders: map[string]string{"Accept": "application/json"}})
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return fmtp.Errorf("huaweicloud_nat_dnat_rule.dnat is not exist")
			}
			return fmtp.Errorf("Error checking huaweicloud_nat_dnat_rule.dnat exist, err=send request failed: %s", err)
		}
		return nil
	}
}

func testAccNatV2DnatRule_base(suffix string) string {
	return fmt.Sprintf(`
data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

resource "huaweicloud_vpc_eip" "eip_1" {
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

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "huaweicloud_compute_instance" "instance_1" {
  name               = "instance-acc-test-%s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = huaweicloud_vpc_subnet.subnet_1.id
  }
  
  tags = {
    foo = "bar"
  }
}
`, suffix)
}

func testAccNatV2DnatRule_basic(suffix string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_nat_dnat_rule" "dnat" {
  nat_gateway_id = huaweicloud_nat_gateway.nat_1.id
  floating_ip_id = huaweicloud_vpc_eip.eip_1.id
  private_ip     = huaweicloud_compute_instance.instance_1.network.0.fixed_ip_v4
  protocol       = "tcp"
  internal_service_port = 993
  external_service_port = 242
}
`, testAccNatV2Gateway_basic(suffix), testAccNatV2DnatRule_base(suffix))
}

func testAccNatV2DnatRule_protocol(suffix string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_nat_dnat_rule" "dnat" {
  nat_gateway_id = huaweicloud_nat_gateway.nat_1.id
  floating_ip_id = huaweicloud_vpc_eip.eip_1.id
  private_ip     = huaweicloud_compute_instance.instance_1.network.0.fixed_ip_v4
  protocol       = "any"
  internal_service_port = 0
  external_service_port = 0
}
`, testAccNatV2Gateway_basic(suffix), testAccNatV2DnatRule_base(suffix))
}
