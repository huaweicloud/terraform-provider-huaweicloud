package dns

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSystemLines_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dns_system_lines.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byLocale   = "data.huaweicloud_dns_system_lines.filter_by_locale"
		dcByLocale = acceptance.InitDataSourceCheck(byLocale)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSystemLines_basic,
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "lines.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "lines.0.id"),
					resource.TestCheckResourceAttrSet(all, "lines.0.name"),
					resource.TestCheckResourceAttrSet(all, "lines.0.father_id"),
					// Filter by 'locale' parameter.
					dcByLocale.CheckResourceExists(),
					resource.TestCheckOutput("is_locale_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSystemLines_basic = `
# Without any filter parameters.
data "huaweicloud_dns_system_lines" "test" {}

# Filter by 'locale' parameter.
data "huaweicloud_dns_system_lines" "filter_by_locale" {
  locale = "en-us"
}

locals {
  filter_result_by_locale = [for v in data.huaweicloud_dns_system_lines.filter_by_locale.lines[*].name :
  length(try(regex("^[A-Za-z]", v), "")) > 0]
}

output "is_locale_filter_useful" {
  value = length(local.filter_result_by_locale) > 0 && alltrue(local.filter_result_by_locale)
}
`
