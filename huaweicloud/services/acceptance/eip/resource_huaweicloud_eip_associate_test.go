package eip

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEIPAssociate_basic(t *testing.T) {
	var eip eips.PublicIp
	rName := acceptance.RandomAccResourceName()
	associateName := "huaweicloud_vpc_eip_associate.test"
	resourceName := "huaweicloud_vpc_eip.test"
	partten := `^((25[0-5]|2[0-4]\d|(1\d{2}|[1-9]?\d))\.){3}(25[0-5]|2[0-4]\d|(1\d{2}|[1-9]?\d))$`

	// huaweicloud_vpc_eip_associate and huaweicloud_vpc_eip have the same ID
	// and call the same API to get resource
	rc := acceptance.InitResourceCheck(
		associateName,
		&eip,
		getEipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEIPAssociate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(associateName, "status", "BOUND"),
					resource.TestCheckResourceAttrPair(
						associateName, "public_ip", resourceName, "address"),
					resource.TestMatchOutput("public_ip_address", regexp.MustCompile(partten)),
				),
			},
			{
				ResourceName:      associateName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccEIPAssociate_port(t *testing.T) {
	var eip eips.PublicIp
	rName := acceptance.RandomAccResourceName()
	associateName := "huaweicloud_vpc_eip_associate.test"
	resourceName := "huaweicloud_vpc_eip.test"

	// huaweicloud_vpc_eip_associate and huaweicloud_vpc_eip have the same ID
	// and call the same API to get resource
	rc := acceptance.InitResourceCheck(
		associateName,
		&eip,
		getEipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEIPAssociate_port(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(associateName, "status", "BOUND"),
					resource.TestCheckResourceAttrPtr(
						associateName, "port_id", &eip.PortID),
					resource.TestCheckResourceAttrPair(
						associateName, "public_ip", resourceName, "address"),
				),
			},
			{
				ResourceName:      associateName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccEIPAssociate_compatible(t *testing.T) {
	var eip eips.PublicIp
	rName := acceptance.RandomAccResourceName()
	associateName := "huaweicloud_networking_eip_associate.test"
	resourceName := "huaweicloud_vpc_eip.test"

	// huaweicloud_networking_eip_associate and huaweicloud_vpc_eip have the same ID
	// and call the same API to get resource
	rc := acceptance.InitResourceCheck(
		associateName,
		&eip,
		getEipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEIPAssociate_compatible(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPtr(
						associateName, "port_id", &eip.PortID),
					resource.TestCheckResourceAttrPair(
						associateName, "public_ip", resourceName, "address"),
				),
			},
			{
				ResourceName:      associateName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEIPAssociate_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%[1]s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.test.id
}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    size        = 5
    name        = "%[1]s"
    charge_mode = "traffic"
  }
}`, rName)
}

func testAccEIPAssociate_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 8
  memory_size       = 16
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = "%[2]s"
  delete_default_rules = true
}

resource "huaweicloud_kps_keypair" "test" {
  name = "%[2]s"
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]

  key_pair = huaweicloud_kps_keypair.test.name

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_vpc_eip_associate" "test" {
  public_ip  = huaweicloud_vpc_eip.test.address
  network_id = huaweicloud_compute_instance.test.network[0].uuid
  fixed_ip   = huaweicloud_compute_instance.test.network[0].fixed_ip_v4
}

data "huaweicloud_compute_instance" "test" {
  depends_on = [huaweicloud_vpc_eip_associate.test]

  name = "%[2]s"
}

output "public_ip_address" {
  value = data.huaweicloud_compute_instance.test.public_ip
}
`, testAccEIPAssociate_base(rName), rName)
}

func testAccEIPAssociate_port(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_vip" "test" {
  name       = "%s"
  network_id = huaweicloud_vpc_subnet.test.id
}

resource "huaweicloud_vpc_eip_associate" "test" {
  public_ip = huaweicloud_vpc_eip.test.address
  port_id   = huaweicloud_networking_vip.test.id
}
`, testAccEIPAssociate_base(rName), rName)
}

func testAccEIPAssociate_compatible(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_vip" "test" {
  name       = "%s"
  network_id = huaweicloud_vpc_subnet.test.id
}
  
resource "huaweicloud_networking_eip_associate" "test" {
  public_ip = huaweicloud_vpc_eip.test.address
  port_id   = huaweicloud_networking_vip.test.id
}
`, testAccEIPAssociate_base(rName), rName)
}
