package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/compute/v2/servers"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccComputeV2EIPAssociate_basic(t *testing.T) {
	var instance servers.Server
	var eip eips.PublicIp

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_compute_eip_associate.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2EIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2EIPAssociate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("huaweicloud_compute_instance.test", &instance),
					testAccCheckVpcV1EIPExists("huaweicloud_vpc_eip.test", &eip),
					testAccCheckComputeV2EIPAssociateAssociated(&eip, &instance, 1),
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

func TestAccComputeV2EIPAssociate_fixedIP(t *testing.T) {
	var instance servers.Server
	var eip eips.PublicIp

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_compute_eip_associate.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2EIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2EIPAssociate_fixedIP(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("huaweicloud_compute_instance.test", &instance),
					testAccCheckVpcV1EIPExists("huaweicloud_vpc_eip.test", &eip),
					testAccCheckComputeV2EIPAssociateAssociated(&eip, &instance, 1),
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

func testAccCheckComputeV2EIPAssociateDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	computeClient, err := config.ComputeV2Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_compute_eip_associate" {
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
			if _, ok := err.(golangsdk.ErrDefault404); ok {
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
					return fmtp.Errorf("EIP %s is still attached to instance %s", floatingIP, instanceId)
				}
			}
		}
	}

	return nil
}

func testAccCheckComputeV2EIPAssociateAssociated(
	eip *eips.PublicIp, instance *servers.Server, n int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*config.Config)
		computeClient, err := config.ComputeV2Client(HW_REGION_NAME)

		newInstance, err := servers.Get(computeClient, instance.ID).Extract()
		if err != nil {
			return err
		}

		// Walk through the instance's addresses and find the match
		i := 0
		for _, networkAddresses := range newInstance.Addresses {
			i++
			if i != n {
				continue
			}
			for _, element := range networkAddresses.([]interface{}) {
				address := element.(map[string]interface{})
				if address["OS-EXT-IPS:type"] == "floating" && address["addr"] == eip.PublicAddress {
					return nil
				}
			}
		}
		return fmtp.Errorf("EIP %s was not attached to instance %s", eip.PublicAddress, instance.ID)
	}
}

func testAccComputeV2EIPAssociate_Base(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "test" {
  name               = "%s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%s"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}
`, testAccCompute_data, rName, rName)
}

func testAccComputeV2EIPAssociate_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_eip_associate" "test" {
  public_ip   = huaweicloud_vpc_eip.test.address
  instance_id = huaweicloud_compute_instance.test.id
}
`, testAccComputeV2EIPAssociate_Base(rName))
}

func testAccComputeV2EIPAssociate_fixedIP(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_eip_associate" "test" {
  public_ip   = huaweicloud_vpc_eip.test.address
  instance_id = huaweicloud_compute_instance.test.id
  fixed_ip    = huaweicloud_compute_instance.test.access_ip_v4
}
`, testAccComputeV2EIPAssociate_Base(rName))
}
