package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/peerings"
)

func TestAccVpcPeeringConnectionV2_basic(t *testing.T) {
	var peering peerings.Peering

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_vpc_peering_connection.test"
	rNameUpdate := rName + "updated"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcPeeringConnectionV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPeeringConnectionV2_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcPeeringConnectionV2Exists(resourceName, &peering),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccVpcPeeringConnectionV2_basic(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
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

func testAccCheckVpcPeeringConnectionV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	peeringClient, err := config.NetworkingHwV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating huaweicloud Peering client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_vpc_peering_connection_v2" {
			continue
		}

		_, err := peerings.Get(peeringClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Vpc Peering Connection still exists")
		}
	}

	return nil
}

func testAccCheckVpcPeeringConnectionV2Exists(n string, peering *peerings.Peering) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		peeringClient, err := config.NetworkingHwV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating huaweicloud Peering client: %s", err)
		}

		found, err := peerings.Get(peeringClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Vpc peering Connection not found")
		}

		*peering = *found

		return nil
	}
}

func testAccVpcPeeringConnectionV2_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc" "test2" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_peering_connection" "test" {
  name        = "%s"
  vpc_id      = huaweicloud_vpc.test.id
  peer_vpc_id = huaweicloud_vpc.test2.id
}
`, rName, rName+"2", rName)
}
