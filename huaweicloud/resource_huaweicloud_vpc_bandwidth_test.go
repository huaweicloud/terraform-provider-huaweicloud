package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v1/bandwidths"
)

func TestAccVpcBandWidthV2_basic(t *testing.T) {
	var bandwidth bandwidths.BandWidth

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_vpc_bandwidth.test"
	rNameUpdate := rName + "-updated"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcBandWidthV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcBandWidthV2_basic(rName, 5),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcBandWidthV2Exists(resourceName, &bandwidth),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "size", "5"),
				),
			},
			{
				Config: testAccVpcBandWidthV2_basic(rNameUpdate, 6),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcBandWidthV2Exists(resourceName, &bandwidth),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "size", "6"),
				),
			},
		},
	})
}

func TestAccVpcBandWidthV2_WithEpsId(t *testing.T) {
	var bandwidth bandwidths.BandWidth

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_vpc_bandwidth.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEpsID(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcBandWidthV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcBandWidthV2_epsId(rName, 5),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcBandWidthV2Exists(resourceName, &bandwidth),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", OS_ENTERPRISE_PROJECT_ID),
				),
			},
		},
	})
}

func testAccCheckVpcBandWidthV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating huaweicloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_vpc_bandwidth" {
			continue
		}

		_, err := bandwidths.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("BandWidth still exists")
		}
	}

	return nil
}

func testAccCheckVpcBandWidthV2Exists(n string, bandwidth *bandwidths.BandWidth) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.NetworkingV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating huaweicloud networking client: %s", err)
		}

		found, err := bandwidths.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("bandwidth not found")
		}

		*bandwidth = found

		return nil
	}
}

func testAccVpcBandWidthV2_basic(rName string, size int) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_bandwidth" "test" {
	name = "%s"
	size = "%d"
}
`, rName, size)
}

func testAccVpcBandWidthV2_epsId(rName string, size int) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_bandwidth" "test" {
	name = "%s"
	size = "%d"
	enterprise_project_id = "%s"
}
`, rName, size, OS_ENTERPRISE_PROJECT_ID)
}
