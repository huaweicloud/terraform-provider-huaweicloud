package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v1/bandwidths"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getBandwidthAssociateResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.NetworkingV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating bandwidth client: %s", err)
	}

	bwID := state.Primary.Attributes["bandwidth_id"]
	return bandwidths.Get(c, bwID).Extract()
}

func TestAccBandWidthAssociate_basic(t *testing.T) {
	var bandwidth bandwidths.BandWidth

	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_vpc_bandwidth_associate.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&bandwidth,
		getBandwidthAssociateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccBandWidthAssociate_basic(randName, 0),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_name", randName),
					resource.TestCheckResourceAttrPair(resourceName, "bandwidth_id", "huaweicloud_vpc_bandwidth.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "eip_id", "huaweicloud_vpc_eip.test.0", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip", "huaweicloud_vpc_eip.test.0", "address"),
				),
			},
			{
				Config: testAccBandWidthAssociate_basic(randName, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "eip_id", "huaweicloud_vpc_eip.test.1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip", "huaweicloud_vpc_eip.test.1", "address"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBandWidthAssociate_ipv6Port(randName, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "eip_id", ""),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "6"),
					resource.TestCheckResourceAttr(resourceName, "public_ip_type", "5_dualStack"),
					resource.TestCheckResourceAttrPair(resourceName, "port_id", "huaweicloud_networking_vip.test.0", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "network_id", "huaweicloud_networking_vip.test.0", "network_id"),
					resource.TestCheckResourceAttrPair(resourceName, "fixed_ip", "huaweicloud_networking_vip.test.0", "ip_address"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ipv6", "huaweicloud_networking_vip.test.0", "ip_address"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBandWidthAssociate_basic(randName, 0),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_name", randName),
					resource.TestCheckResourceAttrPair(resourceName, "bandwidth_id", "huaweicloud_vpc_bandwidth.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "eip_id", "huaweicloud_vpc_eip.test.0", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip", "huaweicloud_vpc_eip.test.0", "address"),
					resource.TestCheckResourceAttr(resourceName, "port_id", ""),
					resource.TestCheckResourceAttr(resourceName, "network_id", ""),
					resource.TestCheckResourceAttr(resourceName, "fixed_ip", ""),
				),
			},
			{
				Config: testAccBandWidthAssociate_ipv6Ip(randName, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "eip_id", ""),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "6"),
					resource.TestCheckResourceAttr(resourceName, "public_ip_type", "5_dualStack"),
					resource.TestCheckResourceAttrPair(resourceName, "port_id", "huaweicloud_networking_vip.test.0", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "network_id", "huaweicloud_networking_vip.test.0", "network_id"),
					resource.TestCheckResourceAttrPair(resourceName, "fixed_ip", "huaweicloud_networking_vip.test.0", "ip_address"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ipv6", "huaweicloud_networking_vip.test.0", "ip_address"),
				),
			},
			{
				Config: testAccBandWidthAssociate_ipv6Ip(randName, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "eip_id", ""),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "6"),
					resource.TestCheckResourceAttr(resourceName, "public_ip_type", "5_dualStack"),
					resource.TestCheckResourceAttrPair(resourceName, "port_id", "huaweicloud_networking_vip.test.1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "network_id", "huaweicloud_networking_vip.test.1", "network_id"),
					resource.TestCheckResourceAttrPair(resourceName, "fixed_ip", "huaweicloud_networking_vip.test.1", "ip_address"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ipv6", "huaweicloud_networking_vip.test.1", "ip_address"),
				),
			},
			{
				Config: testAccBandWidthAssociate_ipv6Port(randName, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "eip_id", ""),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "6"),
					resource.TestCheckResourceAttr(resourceName, "public_ip_type", "5_dualStack"),
					resource.TestCheckResourceAttrPair(resourceName, "port_id", "huaweicloud_networking_vip.test.0", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "network_id", "huaweicloud_networking_vip.test.0", "network_id"),
					resource.TestCheckResourceAttrPair(resourceName, "fixed_ip", "huaweicloud_networking_vip.test.0", "ip_address"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ipv6", "huaweicloud_networking_vip.test.0", "ip_address"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccBandWidthAssociateIpv6ImportStateFunc(resourceName),
			},
		},
	})
}

func TestAccBandWidthAssociate_migrate(t *testing.T) {
	var bandwidth bandwidths.BandWidth

	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_vpc_bandwidth_associate.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&bandwidth,
		getBandwidthAssociateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccBandWidthAssociate_migrate(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_name", randName),
					resource.TestCheckResourceAttrPair(resourceName, "bandwidth_id", "huaweicloud_vpc_bandwidth.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "eip_id", "huaweicloud_vpc_eip.source", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip", "huaweicloud_vpc_eip.source", "address"),
				),
			},
			{
				Config: testAccBandWidthAssociate_owner(randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "eip_id", "huaweicloud_vpc_eip.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip", "huaweicloud_vpc_eip.test", "address"),
				),
			},
		},
	})
}

func TestAccBandWidthAssociate_migrate_ipv6(t *testing.T) {
	var bandwidth bandwidths.BandWidth

	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_vpc_bandwidth_associate.test_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&bandwidth,
		getBandwidthAssociateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccBandWidthAssociate_ipv6Ip_migrate(randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "bandwidth_id", "huaweicloud_vpc_bandwidth.test.1", "id"),
					resource.TestCheckResourceAttr(resourceName, "eip_id", ""),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "6"),
					resource.TestCheckResourceAttr(resourceName, "public_ip_type", "5_dualStack"),
					resource.TestCheckResourceAttrPair(resourceName, "port_id", "huaweicloud_networking_vip.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "network_id", "huaweicloud_networking_vip.test", "network_id"),
					resource.TestCheckResourceAttrPair(resourceName, "fixed_ip", "huaweicloud_networking_vip.test", "ip_address"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ipv6", "huaweicloud_networking_vip.test", "ip_address"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccBandWidthAssociate_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id      = huaweicloud_vpc.test.id
  name        = "%[1]s"
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
  ipv6_enable = true
}

resource "huaweicloud_networking_vip" "test" {
  count = 2

  name       = "%[1]s-${count.index}"
  network_id = huaweicloud_vpc_subnet.test.id
  ip_version = 6
}

resource "huaweicloud_vpc_bandwidth" "test" {
  name = "%[1]s"
  size = 5
}

resource "huaweicloud_vpc_eip" "test" {
  count = 2
  name  = "%[1]s-${count.index}"

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    name        = "%[1]s-${count.index}"
    size        = 5
    charge_mode = "traffic"
  }

  lifecycle {
    ignore_changes = [ bandwidth ]
  }
}
`, rName)
}

func testAccBandWidthAssociate_basic(rName string, index int) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_bandwidth_associate" "test" {
  bandwidth_id = huaweicloud_vpc_bandwidth.test.id
  eip_id       = huaweicloud_vpc_eip.test.%d.id
}
`, testAccBandWidthAssociate_base(rName), index)
}

func testAccBandWidthAssociate_ipv6Port(rName string, index int) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_bandwidth_associate" "test" {
  bandwidth_id = huaweicloud_vpc_bandwidth.test.id
  port_id      = huaweicloud_networking_vip.test.%d.id
}
`, testAccBandWidthAssociate_base(rName), index)
}

func testAccBandWidthAssociate_ipv6Ip(rName string, index int) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_bandwidth_associate" "test" {
  bandwidth_id = huaweicloud_vpc_bandwidth.test.id
  fixed_ip     = huaweicloud_networking_vip.test.%d.ip_address
  network_id   = huaweicloud_vpc_subnet.test.id
}
`, testAccBandWidthAssociate_base(rName), index)
}

func testAccBandWidthAssociate_migrate(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_bandwidth" "test" {
  name = "%[1]s"
  size = 5
}

resource "huaweicloud_vpc_bandwidth" "source" {
  name = "%[1]s-source"
  size = 5
}

resource "huaweicloud_vpc_eip" "source" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type = "WHOLE"
    id         = huaweicloud_vpc_bandwidth.source.id
  }

  lifecycle {
    ignore_changes = [ bandwidth ]
  }
}

resource "huaweicloud_vpc_bandwidth_associate" "test" {
  bandwidth_id = huaweicloud_vpc_bandwidth.test.id
  eip_id       = huaweicloud_vpc_eip.source.id
}
`, rName)
}

func testAccBandWidthAssociate_owner(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_bandwidth" "test" {
  name = "%[1]s"
  size = 5
}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type = "WHOLE"
    id         = huaweicloud_vpc_bandwidth.test.id
  }

  lifecycle {
    ignore_changes = [ bandwidth ]
  }
}

resource "huaweicloud_vpc_bandwidth_associate" "test" {
  bandwidth_id = huaweicloud_vpc_bandwidth.test.id
  eip_id       = huaweicloud_vpc_eip.test.id
}
`, rName)
}

func testAccBandWidthAssociate_ipv6Ip_migrate(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id      = huaweicloud_vpc.test.id
  name        = "%[1]s"
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
  ipv6_enable = true
}

resource "huaweicloud_networking_vip" "test" {
  name       = "%[1]s"
  network_id = huaweicloud_vpc_subnet.test.id
  ip_version = 6
}

resource "huaweicloud_vpc_bandwidth" "test" {
  count = 2

  name = "%[1]s"
  size = 5
}

resource "huaweicloud_vpc_bandwidth_associate" "test_0" {
  bandwidth_id = huaweicloud_vpc_bandwidth.test.0.id
  fixed_ip     = huaweicloud_networking_vip.test.ip_address
  network_id   = huaweicloud_vpc_subnet.test.id
}

resource "huaweicloud_vpc_bandwidth_associate" "test_1" {
  depends_on = [huaweicloud_vpc_bandwidth_associate.test_0]

  bandwidth_id = huaweicloud_vpc_bandwidth.test.1.id
  fixed_ip     = huaweicloud_networking_vip.test.ip_address
  network_id   = huaweicloud_vpc_subnet.test.id
}
`, rName)
}

func testAccBandWidthAssociateIpv6ImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		bwID := rs.Primary.Attributes["bandwidth_id"]
		networkID := rs.Primary.Attributes["network_id"]
		fixedIP := rs.Primary.Attributes["fixed_ip"]
		if bwID == "" || networkID == "" || fixedIP == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<bandwidth_id>/<network_id>/<fixed_ip>', but got '%s/%s/%s'",
				bwID, networkID, fixedIP)
		}

		return fmt.Sprintf("%s/%s/%s", bwID, networkID, fixedIP), nil
	}
}
