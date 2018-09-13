package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/huaweicloud/golangsdk/openstack/compute/v2/extensions/servergroups"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/servers"
)

func TestAccComputeV2ServerGroup_basic(t *testing.T) {
	var sg servergroups.ServerGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2ServerGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeV2ServerGroup_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2ServerGroupExists("huaweicloud_compute_servergroup_v2.sg_1", &sg),
				),
			},
		},
	})
}

func TestAccComputeV2ServerGroup_affinity(t *testing.T) {
	var instance servers.Server
	var sg servergroups.ServerGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2ServerGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeV2ServerGroup_affinity,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2ServerGroupExists("huaweicloud_compute_servergroup_v2.sg_1", &sg),
					testAccCheckComputeV2InstanceExists("huaweicloud_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceInServerGroup(&instance, &sg),
				),
			},
		},
	})
}

func testAccCheckComputeV2ServerGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	computeClient, err := config.computeV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_compute_servergroup_v2" {
			continue
		}

		_, err := servergroups.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("ServerGroup still exists")
		}
	}

	return nil
}

func testAccCheckComputeV2ServerGroupExists(n string, kp *servergroups.ServerGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.computeV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
		}

		found, err := servergroups.Get(computeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("ServerGroup not found")
		}

		*kp = *found

		return nil
	}
}

func testAccCheckComputeV2InstanceInServerGroup(instance *servers.Server, sg *servergroups.ServerGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(sg.Members) > 0 {
			for _, m := range sg.Members {
				if m == instance.ID {
					return nil
				}
			}
		}

		return fmt.Errorf("Instance %s is not part of Server Group %s", instance.ID, sg.ID)
	}
}

const testAccComputeV2ServerGroup_basic = `
resource "huaweicloud_compute_servergroup_v2" "sg_1" {
  name = "sg_1"
  policies = ["affinity"]
}
`

var testAccComputeV2ServerGroup_affinity = fmt.Sprintf(`
resource "huaweicloud_compute_servergroup_v2" "sg_1" {
  name = "sg_1"
  policies = ["affinity"]
}

resource "huaweicloud_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
  scheduler_hints {
    group = "${huaweicloud_compute_servergroup_v2.sg_1.id}"
  }
  network {
    uuid = "%s"
  }
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)
