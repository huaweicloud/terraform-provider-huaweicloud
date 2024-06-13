package dns

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceQuotas_basic(t *testing.T) {
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
				Config: testDataSourceDataSourceQuotas_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "quotas.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.max"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.used"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.unit"),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceQuotas_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dns_quotas" "test" {
  domain_id = "%[1]s"
}

locals {
  resource_type = data.huaweicloud_dns_quotas.test.quotas[0].type
}

data "huaweicloud_dns_quotas" "filter_by_type" {
  domain_id = "%[1]s"
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
`, acceptance.HW_DOMAIN_ID)
}

func TestAccDataSourceQuotas_expectError(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceQuotas_expectError(),
				// The domain ID not exist.
				ExpectError: regexp.MustCompile("Authentication required"),
			},
		},
	})
}

func testAccDataSourceQuotas_expectError() string {
	return `
data "huaweicloud_dns_quotas" "not_found_domain_id" {
  domain_id = "notfound"
}
`
}
