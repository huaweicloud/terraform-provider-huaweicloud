package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIECPublicIPsDataSource_basic(t *testing.T) {
	resourceName := "data.huaweicloud_iec_eips.site"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		CheckDestroy:              testAccCheckIecEIPDestroy,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccIECEipsDataSource_config,
			},
			{
				Config: testAccIECEipsDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccIECPublicIPsDataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "eips.0.ip_version", "4"),
					resource.TestCheckResourceAttr(resourceName, "eips.0.bandwidth_share_type", "WHOLE"),
					resource.TestCheckResourceAttrSet(resourceName, "site_info"),
					resource.TestCheckResourceAttrSet(resourceName, "eips.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "eips.1.id"),
					resource.TestCheckResourceAttrSet(resourceName, "eips.0.status"),
					resource.TestCheckResourceAttrSet(resourceName, "eips.0.public_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "eips.0.bandwidth_id"),
					resource.TestCheckResourceAttrSet(resourceName, "eips.0.bandwidth_size"),
				),
			},
		},
	})
}

func testAccIECPublicIPsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find IEC public IPs data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("IEC public IPs data source ID not set")
		}

		return nil
	}
}

var testAccIECEipsDataSource_config string = `
data "huaweicloud_iec_sites" "sites_test" {}

resource "huaweicloud_iec_eip" "eip_test1" {
  site_id = data.huaweicloud_iec_sites.sites_test.sites[0].id
}

resource "huaweicloud_iec_eip" "eip_test2" {
  site_id = data.huaweicloud_iec_sites.sites_test.sites[0].id
}
`

func testAccIECEipsDataSource_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_iec_eips" "site" {
  site_id = data.huaweicloud_iec_sites.sites_test.sites[0].id
}
`, testAccIECEipsDataSource_config)
}
