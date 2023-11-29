package iec

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/iec/v1/common"
	"github.com/chnsz/golangsdk/openstack/iec/v1/publicips"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEIPResource_basic(t *testing.T) {
	var iecEip common.PublicIP
	resourceName := "huaweicloud_iec_eip.eip_test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIP_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(resourceName, &iecEip),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "4"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_share_type", "WHOLE"),
					resource.TestMatchResourceAttr(resourceName, "public_ip",
						regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)),
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

func testAccCheckEIPDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	iecV1Client, err := cfg.IECV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating IEC client: %s", err)
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

func testAccCheckEIPExists(n string, ipResource *common.PublicIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		iecV1Client, err := cfg.IECV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating IEC client: %s", err)
		}

		found, err := publicips.Get(iecV1Client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("IEC EIP not found")
		}

		*ipResource = *found

		return nil
	}
}

var testAccEIP_basic = `
data "huaweicloud_iec_sites" "sites_test" {}

resource "huaweicloud_iec_eip" "eip_test" {
  site_id = data.huaweicloud_iec_sites.sites_test.sites[0].id
}
`
