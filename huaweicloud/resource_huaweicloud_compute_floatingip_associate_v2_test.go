package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccComputeV2FloatingIPAssociate_basic(t *testing.T) {
	var instance servers.Server
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2FloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeV2FloatingIPAssociate_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("huaweicloud_compute_instance_v2.instance_1", &instance),
					testAccCheckNetworkingV2FloatingIPExists("huaweicloud_networking_floatingip_v2.fip_1", &fip),
					testAccCheckComputeV2FloatingIPAssociateAssociated(&fip, &instance, 1),
				),
			},
		},
	})
}

func TestAccComputeV2FloatingIPAssociate_fixedIP(t *testing.T) {
	var instance servers.Server
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2FloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeV2FloatingIPAssociate_fixedIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("huaweicloud_compute_instance_v2.instance_1", &instance),
					testAccCheckNetworkingV2FloatingIPExists("huaweicloud_networking_floatingip_v2.fip_1", &fip),
					testAccCheckComputeV2FloatingIPAssociateAssociated(&fip, &instance, 1),
				),
			},
		},
	})
}

func TestAccComputeV2FloatingIPAssociate_attachToFirstNetwork(t *testing.T) {
	var instance servers.Server
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2FloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeV2FloatingIPAssociate_attachToFirstNetwork,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("huaweicloud_compute_instance_v2.instance_1", &instance),
					testAccCheckNetworkingV2FloatingIPExists("huaweicloud_networking_floatingip_v2.fip_1", &fip),
					testAccCheckComputeV2FloatingIPAssociateAssociated(&fip, &instance, 1),
				),
			},
		},
	})
}

func TestAccComputeV2FloatingIPAssociate_attachNew(t *testing.T) {
	var instance servers.Server
	var fip_1 floatingips.FloatingIP
	var fip_2 floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2FloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeV2FloatingIPAssociate_attachNew_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("huaweicloud_compute_instance_v2.instance_1", &instance),
					testAccCheckNetworkingV2FloatingIPExists("huaweicloud_networking_floatingip_v2.fip_1", &fip_1),
					testAccCheckNetworkingV2FloatingIPExists("huaweicloud_networking_floatingip_v2.fip_2", &fip_2),
					testAccCheckComputeV2FloatingIPAssociateAssociated(&fip_1, &instance, 1),
				),
			},
			resource.TestStep{
				Config: testAccComputeV2FloatingIPAssociate_attachNew_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("huaweicloud_compute_instance_v2.instance_1", &instance),
					testAccCheckNetworkingV2FloatingIPExists("huaweicloud_networking_floatingip_v2.fip_1", &fip_1),
					testAccCheckNetworkingV2FloatingIPExists("huaweicloud_networking_floatingip_v2.fip_2", &fip_2),
					testAccCheckComputeV2FloatingIPAssociateAssociated(&fip_2, &instance, 1),
				),
			},
		},
	})
}

func testAccCheckComputeV2FloatingIPAssociateDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	computeClient, err := config.computeV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_compute_floatingip_associate_v2" {
			continue
		}

		floatingIP, instanceId, _, err := parseComputeFloatingIPAssociateId(rs.Primary.ID)
		if err != nil {
			return err
		}

		instance, err := servers.Get(computeClient, instanceId).Extract()
		if err != nil {
			// If the error is a 404, then the instance does not exist,
			// and therefore the floating IP cannot be associated to it.
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return nil
			}
			return err
		}

		// But if the instance still exists, then walk through its known addresses
		// and see if there's a floating IP.
		for _, networkAddresses := range instance.Addresses {
			for _, element := range networkAddresses.([]interface{}) {
				address := element.(map[string]interface{})
				if address["OS-EXT-IPS:type"] == "floating" || address["OS-EXT-IPS:type"] == "fixed" {
					return fmt.Errorf("Floating IP %s is still attached to instance %s", floatingIP, instanceId)
				}
			}
		}
	}

	return nil
}

func testAccCheckComputeV2FloatingIPAssociateAssociated(
	fip *floatingips.FloatingIP, instance *servers.Server, n int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.computeV2Client(OS_REGION_NAME)

		newInstance, err := servers.Get(computeClient, instance.ID).Extract()
		if err != nil {
			return err
		}

		// Walk through the instance's addresses and find the match
		i := 0
		for _, networkAddresses := range newInstance.Addresses {
			i += 1
			if i != n {
				continue
			}
			for _, element := range networkAddresses.([]interface{}) {
				address := element.(map[string]interface{})
				if (address["OS-EXT-IPS:type"] == "floating" && address["addr"] == fip.FloatingIP) ||
					(address["OS-EXT-IPS:type"] == "fixed" && address["addr"] == fip.FixedIP) {
					return nil
				}
			}
		}
		return fmt.Errorf("Floating IP %s was not attached to instance %s", fip.FloatingIP, instance.ID)
	}
}

var testAccComputeV2FloatingIPAssociate_basic = fmt.Sprintf(`
resource "huaweicloud_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
  network {
    uuid = "%s"
  }
}

resource "huaweicloud_networking_floatingip_v2" "fip_1" {
}

resource "huaweicloud_compute_floatingip_associate_v2" "fip_1" {
  floating_ip = "${huaweicloud_networking_floatingip_v2.fip_1.address}"
  instance_id = "${huaweicloud_compute_instance_v2.instance_1.id}"
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)

var testAccComputeV2FloatingIPAssociate_fixedIP = fmt.Sprintf(`
resource "huaweicloud_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
  network {
    uuid = "%s"
  }
}

resource "huaweicloud_networking_floatingip_v2" "fip_1" {
}

resource "huaweicloud_compute_floatingip_associate_v2" "fip_1" {
  floating_ip = "${huaweicloud_networking_floatingip_v2.fip_1.address}"
  instance_id = "${huaweicloud_compute_instance_v2.instance_1.id}"
  fixed_ip = "${huaweicloud_compute_instance_v2.instance_1.access_ip_v4}"
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)

var testAccComputeV2FloatingIPAssociate_attachToFirstNetwork = fmt.Sprintf(`
resource "huaweicloud_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"

  network {
    uuid = "%s"
  }
}

resource "huaweicloud_networking_floatingip_v2" "fip_1" {
}

resource "huaweicloud_compute_floatingip_associate_v2" "fip_1" {
  floating_ip = "${huaweicloud_networking_floatingip_v2.fip_1.address}"
  instance_id = "${huaweicloud_compute_instance_v2.instance_1.id}"
  fixed_ip = "${huaweicloud_compute_instance_v2.instance_1.network.0.fixed_ip_v4}"
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)

var testAccComputeV2FloatingIPAssociate_attachToSecondNetwork = fmt.Sprintf(`
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
  cidr = "192.168.1.0/24"
  ip_version = 4
  enable_dhcp = true
  no_gateway = true
}

resource "huaweicloud_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"

  network {
    uuid = "${huaweicloud_networking_network_v2.network_1.id}"
  }

  network {
    uuid = "%s"
  }
}

resource "huaweicloud_networking_floatingip_v2" "fip_1" {
}

resource "huaweicloud_compute_floatingip_associate_v2" "fip_1" {
  floating_ip = "${huaweicloud_networking_floatingip_v2.fip_1.address}"
  instance_id = "${huaweicloud_compute_instance_v2.instance_1.id}"
  fixed_ip = "${huaweicloud_compute_instance_v2.instance_1.network.1.fixed_ip_v4}"
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)

var testAccComputeV2FloatingIPAssociate_attachNew_1 = fmt.Sprintf(`
resource "huaweicloud_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
  network {
    uuid = "%s"
  }
}

resource "huaweicloud_networking_floatingip_v2" "fip_1" {
}

resource "huaweicloud_networking_floatingip_v2" "fip_2" {
}

resource "huaweicloud_compute_floatingip_associate_v2" "fip_1" {
  floating_ip = "${huaweicloud_networking_floatingip_v2.fip_1.address}"
  instance_id = "${huaweicloud_compute_instance_v2.instance_1.id}"
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)

var testAccComputeV2FloatingIPAssociate_attachNew_2 = fmt.Sprintf(`
resource "huaweicloud_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
  network {
    uuid = "%s"
  }
}

resource "huaweicloud_networking_floatingip_v2" "fip_1" {
}

resource "huaweicloud_networking_floatingip_v2" "fip_2" {
}

resource "huaweicloud_compute_floatingip_associate_v2" "fip_1" {
  floating_ip = "${huaweicloud_networking_floatingip_v2.fip_2.address}"
  instance_id = "${huaweicloud_compute_instance_v2.instance_1.id}"
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)
