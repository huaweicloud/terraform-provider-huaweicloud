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

func TestAccDataZoneNameservers_basic(t *testing.T) {
	var (
		name = fmt.Sprintf("acpttest-zone-%s.com.", acctest.RandString(5))

		dataSource = "data.huaweicloud_dns_zone_nameservers.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataZoneNameservers_notFound(),
				ExpectError: regexp.MustCompile(`This zone does not exist`),
			},
			{
				Config: testAccDataZoneNameservers_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "nameservers.#", regexp.MustCompile(`[1-9][0-9]*`)),
					resource.TestCheckResourceAttrSet(dataSource, "nameservers.0.hostname"),
					resource.TestCheckResourceAttrSet(dataSource, "nameservers.0.priority"),
					resource.TestMatchResourceAttr(dataSource, "nameservers.0.hostname", regexp.MustCompile(`.+\.`)),
				),
			},
		},
	})
}

func testAccDataZoneNameservers_notFound() string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dns_zone_nameservers" "invalid_zone_id" {
  zone_id = "%[1]s"
}
`, randomId)
}

func testAccDataZoneNameservers_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_zone" "test" {
  name      = "%[1]s"
  zone_type = "public"
}

data "huaweicloud_dns_zone_nameservers" "test" {
  zone_id = huaweicloud_dns_zone.test.id
}
`, name)
}
