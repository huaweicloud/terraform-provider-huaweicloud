package dns

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataPublicZoneDetectionStatus_basic(t *testing.T) {
	var (
		name = fmt.Sprintf("acpttest.zone.%s.com.", acctest.RandString(5))

		dcName = "data.huaweicloud_dns_public_zone_detection_status.test"
		dc     = acceptance.InitDataSourceCheck(dcName)

		byType   = "data.huaweicloud_dns_public_zone_detection_status.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPublicZoneDetectionStatus_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dcName, "status"),
					// Filter by 'type' parameter.
					dcByType.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(byType, "status"),
				),
			},
		},
	})
}

func testAccDataPublicZoneDetectionStatus_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_zone" "test" {
  name      = "%[1]s"
  zone_type = "public"
}

resource "huaweicloud_dns_recordset" "test" {
  zone_id = huaweicloud_dns_zone.test.id
  name    = "tf.${huaweicloud_dns_zone.test.name}"
  type    = "TXT"
  records = ["\"test records\""]
}

# Without any filter parameters.
data "huaweicloud_dns_public_zone_detection_status" "test" {
  zone_id     = huaweicloud_dns_zone.test.id
  domain_name = huaweicloud_dns_zone.test.name
}

# Filter by 'type' parameter.
data "huaweicloud_dns_public_zone_detection_status" "filter_by_type" {
  zone_id     = huaweicloud_dns_zone.test.id
  domain_name = huaweicloud_dns_recordset.test.name
  type        = huaweicloud_dns_recordset.test.type
}`, name)
}
