package dns

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataQuotas_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dns_quotas.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byType   = "data.huaweicloud_dns_quotas.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataQuotas_domainIdNotFound(),
				// The domain ID not exist.
				ExpectError: regexp.MustCompile("Authentication required"),
			},
			{
				Config: testAccDataQuotas_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "quotas.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_used_set_and_valid", "true"),
					// Filter by resource type.
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					// Check attributes.
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.max"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.unit"),
				),
			},
		},
	})
}

func testAccDataQuotas_domainIdNotFound() string {
	return `
data "huaweicloud_dns_quotas" "not_found_domain_id" {
  domain_id = "notfound"
}
`
}

func testAccDataQuotas_basic() string {
	name := fmt.Sprintf("%s.com", acceptance.RandomAccResourceNameWithDash())
	return fmt.Sprintf(`
resource "huaweicloud_dns_zone" "test" {
  name = "%[1]s"
}

data "huaweicloud_dns_quotas" "test" {
  depends_on = [huaweicloud_dns_zone.test]
  domain_id  = "%[2]s"
}

output "is_used_set_and_valid" {
  value = length([for v in data.huaweicloud_dns_quotas.test.quotas[*].used : v if v > 0]) > 0
}

# Filter by resource type.
locals {
  resource_type = data.huaweicloud_dns_quotas.test.quotas[0].type
}

data "huaweicloud_dns_quotas" "filter_by_type" {
  domain_id = "%[2]s"
  type      = local.resource_type
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_dns_quotas.filter_by_type.quotas[*].type : v == local.resource_type
  ]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}
`, name, acceptance.HW_DOMAIN_ID)
}
