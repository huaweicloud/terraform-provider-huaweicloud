package dns

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataPublicZoneRecordsets_basic(t *testing.T) {
	var (
		name = fmt.Sprintf("tf.test.zone-%s.com.", acctest.RandString(5))

		emailName = "data.huaweicloud_dns_public_zone_recordsets.email"
		dcEmail   = acceptance.InitDataSourceCheck(emailName)

		websiteName = "data.huaweicloud_dns_public_zone_recordsets.website"
		dcWebsite   = acceptance.InitDataSourceCheck(websiteName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataPublicZoneRecordsets_notFound(),
				ExpectError: regexp.MustCompile(`This zone does not exist`),
			},
			{
				Config: testAccDataPublicZoneRecordsets_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Query email recordsets.
					dcEmail.CheckResourceExists(),
					resource.TestMatchResourceAttr(emailName, "recordsets.#", regexp.MustCompile(`[1-9][0-9]*`)),
					resource.TestCheckResourceAttrSet(emailName, "recordsets.0.id"),
					resource.TestCheckResourceAttrSet(emailName, "recordsets.0.name"),
					resource.TestCheckResourceAttrSet(emailName, "recordsets.0.zone_id"),
					resource.TestCheckResourceAttrSet(emailName, "recordsets.0.type"),
					resource.TestCheckResourceAttrSet(emailName, "recordsets.0.default"),
					resource.TestMatchResourceAttr(emailName, "recordsets.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(emailName, "recordsets.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Query website recordsets.
					dcWebsite.CheckResourceExists(),
					resource.TestMatchResourceAttr(websiteName, "recordsets.#", regexp.MustCompile(`[1-9][0-9]*`)),
				),
			},
		},
	})
}

func testAccDataPublicZoneRecordsets_notFound() string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dns_public_zone_recordsets" "test" {
  zone_id = replace("%[1]s", "-", "")
  type    = "email"
}
`, randomId)
}

func testAccDataPublicZoneRecordsets_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_zone" "test" {
  name      = "%[1]s"
  zone_type = "public"
}

resource "huaweicloud_dns_recordset" "test" {
  zone_id     = huaweicloud_dns_zone.test.id
  name        = "pop3.${huaweicloud_dns_zone.test.name}"
  type        = "CNAME"
  records     = ["pop3.sparkspace.huaweicloud.com."]
}

# Query email recordsets.
data "huaweicloud_dns_public_zone_recordsets" "email" {
  zone_id = huaweicloud_dns_zone.test.id
  type    = "email"

  depends_on = [huaweicloud_dns_recordset.test]
}

# Query website recordsets.
data "huaweicloud_dns_public_zone_recordsets" "website" {
  zone_id = huaweicloud_dns_zone.test.id
  type    = "website"
}`, name)
}
