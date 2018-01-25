package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var zoneName = fmt.Sprintf("ACPTTEST%s.com.", acctest.RandString(5))

func TestAccHuaweiCloudDNSZoneV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccHuaweiCloudDNSZoneV2DataSource_zone,
			},
			resource.TestStep{
				Config: testAccHuaweiCloudDNSZoneV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSZoneV2DataSourceID("data.huaweicloud_dns_zone_v2.z1"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_dns_zone_v2.z1", "name", zoneName),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_dns_zone_v2.z1", "type", "PRIMARY"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_dns_zone_v2.z1", "ttl", "7200"),
				),
			},
		},
	})
}

func testAccCheckDNSZoneV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find DNS Zone data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("DNS Zone data source ID not set")
		}

		return nil
	}
}

var testAccHuaweiCloudDNSZoneV2DataSource_zone = fmt.Sprintf(`
resource "huaweicloud_dns_zone_v2" "z1" {
  name = "%s"
  email = "terraform-dns-zone-v2-test-name@example.com"
  type = "PRIMARY"
  ttl = 7200
}`, zoneName)

var testAccHuaweiCloudDNSZoneV2DataSource_basic = fmt.Sprintf(`
%s
data "huaweicloud_dns_zone_v2" "z1" {
	name = "${huaweicloud_dns_zone_v2.z1.name}"
}
`, testAccHuaweiCloudDNSZoneV2DataSource_zone)
