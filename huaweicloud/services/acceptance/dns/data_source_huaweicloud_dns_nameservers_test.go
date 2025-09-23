package dns

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataNameservers_basic(t *testing.T) {
	var (
		byType   = "data.huaweicloud_dns_nameservers.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byServerRegion           = "data.huaweicloud_dns_nameservers.filter_by_server_region"
		dcByServerRegion         = acceptance.InitDataSourceCheck(byServerRegion)
		byNotFoundServerRegion   = "data.huaweicloud_dns_nameservers.filter_by_not_found_server_region"
		dcByNotFoundServerRegion = acceptance.InitDataSourceCheck(byNotFoundServerRegion)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataNameservers_basic(),
				Check: resource.ComposeTestCheckFunc(
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					resource.TestCheckResourceAttr(byType, "nameservers.0.type", "public"),
					resource.TestCheckResourceAttrSet(byServerRegion, "nameservers.0.ns_records.0.priority"),
					// For the public name server, the 'address' value is empty.
					resource.TestCheckResourceAttr(byType, "nameservers.0.ns_records.0.address", ""),
					// Filter by server region.
					dcByServerRegion.CheckResourceExists(),
					resource.TestCheckOutput("server_region_filter_is_useful", "true"),
					resource.TestCheckResourceAttr(byServerRegion, "nameservers.0.type", "private"),
					dcByNotFoundServerRegion.CheckResourceExists(),
					resource.TestCheckOutput("server_region_not_found_validation_pass", "true"),
					resource.TestCheckResourceAttrSet(byServerRegion, "nameservers.0.ns_records.0.address"),
					// For the private name server, the 'hostname' value is empty.
					resource.TestCheckResourceAttr(byServerRegion, "nameservers.0.ns_records.0.hostname", ""),
				),
			},
		},
	})
}

func testAccDataNameservers_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dns_nameservers" "filter_by_type" {
  type = "public"
}

locals {
  type_filter_result = [for v in data.huaweicloud_dns_nameservers.filter_by_type.nameservers[*].type : v == "public"]
}

output "is_type_filter_useful" {
  value = alltrue(local.type_filter_result) && length(local.type_filter_result) > 0
}

data "huaweicloud_dns_nameservers" "filter_by_server_region" {
  server_region = "%[1]s"
}

locals {
  server_region_filter_result = [for v in data.huaweicloud_dns_nameservers.filter_by_server_region.nameservers[*].region : v == "%[1]s"]
}

output "server_region_filter_is_useful" {
  value = alltrue(local.server_region_filter_result) && length(local.server_region_filter_result) == 1
}

data "huaweicloud_dns_nameservers" "filter_by_not_found_server_region" {
  server_region = "not_fount"
}
  
output "server_region_not_found_validation_pass" {
  value = length(data.huaweicloud_dns_nameservers.filter_by_not_found_server_region.nameservers) == 0
}
`, acceptance.HW_REGION_NAME)
}
