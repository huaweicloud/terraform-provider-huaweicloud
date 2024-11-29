package live

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceLiveDomains_basic(t *testing.T) {
	var (
		domainName = fmt.Sprintf("%s.huaweicloud.com", acceptance.RandomAccResourceNameWithDash())

		dataSource = "data.huaweicloud_live_domains.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byName   = "data.huaweicloud_live_domains.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLiveDomains_basic(domainName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "domains.0.cname"),
					resource.TestCheckResourceAttrSet(dataSource, "domains.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "domains.0.is_ipv6"),
					resource.TestCheckResourceAttrSet(dataSource, "domains.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "domains.0.region"),
					resource.TestCheckResourceAttrSet(dataSource, "domains.0.service_area"),
					resource.TestCheckResourceAttrSet(dataSource, "domains.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "domains.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "domains.0.vendor"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceLiveDomains_basic(domainName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "test" {
  name = "%s"
  type = "push"
}

data "huaweicloud_live_domains" "test" {
  depends_on = [huaweicloud_live_domain.test]
}

# Filter by name
locals {
  name = data.huaweicloud_live_domains.test.domains[0].name
}

data "huaweicloud_live_domains" "filter_by_name" {
  name = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_live_domains.filter_by_name.domains[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}
`, domainName)
}
