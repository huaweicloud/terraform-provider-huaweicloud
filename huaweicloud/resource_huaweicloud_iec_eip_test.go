package huaweicloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	iec_common "github.com/chnsz/golangsdk/openstack/iec/v1/common"
	"github.com/chnsz/golangsdk/openstack/iec/v1/publicips"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
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
				Config: testAccIecEIP_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIecEIPExists(resourceName, &iecEip),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "4"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_share_type", "WHOLE"),
					resource.TestMatchResourceAttr(resourceName, "public_ip", regexp.MustCompile("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$")),
					resource.TestCheckResourceAttrSet(resourceName, "site_info"),
					resource.TestCheckResourceAttrSet(resourceName, "site_id"),
					resource.TestCheckResourceAttrSet(resourceName, "line_id"),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidth_id"),
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
	config := testAccProvider.Meta().(*config.Config)
	iecV1Client, err := config.IECV1Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud IEC client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_iec_eip" {
			continue
		}

		_, err := publicips.Get(iecV1Client, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("IEC EIP still exists")
		}
	}

	return nil
}

func testAccCheckIecEIPExists(n string, resource *iec_common.PublicIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		iecV1Client, err := config.IECV1Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating Huaweicloud IEC client: %s", err)
		}

		found, err := publicips.Get(iecV1Client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("IEC EIP not found")
		}

		*resource = *found

		return nil
	}
}

var testAccIecEIP_basic string = `
data "huaweicloud_iec_sites" "sites_test" {}

resource "huaweicloud_iec_eip" "eip_test" {
  site_id = data.huaweicloud_iec_sites.sites_test.sites[0].id
}
`
