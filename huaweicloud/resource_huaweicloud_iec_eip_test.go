package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	iec_common "github.com/huaweicloud/golangsdk/openstack/iec/v1/common"
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/publicips"
)

func TestAccIecEIPResource_basic(t *testing.T) {
	var iecEip iec_common.PublicIP
	resourceName := "huaweicloud_iec_eip.eip_test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIecEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIecEIP_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIecEIPExists(resourceName, &iecEip),
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

func testAccCheckIecEIPDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	iecV1Client, err := config.IECV1Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud IEC client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_iec_eip" {
			continue
		}

		_, err := publicips.Get(iecV1Client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("IEC EIP still exists")
		}
	}

	return nil
}

func testAccCheckIecEIPExists(n string, resource *iec_common.PublicIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		iecV1Client, err := config.IECV1Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating Huaweicloud IEC client: %s", err)
		}

		found, err := publicips.Get(iecV1Client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("IEC EIP not found")
		}

		*resource = *found

		return nil
	}
}

func testAccIecEIP_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_iec_sites" "sites_test" {
  region = "%s"
  area   = "east"
  city   = "hangzhou"
}

resource "huaweicloud_iec_eip" "eip_test" {
  site_id    = data.huaweicloud_iec_sites.sites_test.sites[0].id
}
`, HW_REGION_NAME)
}
