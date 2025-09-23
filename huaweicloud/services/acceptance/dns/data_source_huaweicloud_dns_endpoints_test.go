package dns

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDNSEndpoints_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dns_endpoints.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDNSEndpoints_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "endpoints.#"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoints.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoints.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoints.0.vpc_id"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoints.0.direction"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoints.0.ipaddress_count"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoints.0.resolver_rule_count"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoints.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoints.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoints.0.update_time"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_vpc_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDNSEndpoints_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dns_endpoints" "test" {
  direction = huaweicloud_dns_endpoint.test.direction
}

// filter by vpc_id
data "huaweicloud_dns_endpoints" "filter_by_vpc_id" {
  direction = huaweicloud_dns_endpoint.test.direction
  vpc_id    = huaweicloud_dns_endpoint.test.vpc_id
}

locals {
  filter_result_by_vpc_id = [for v in data.huaweicloud_dns_endpoints.filter_by_vpc_id.endpoints[*].vpc_id :
    v == huaweicloud_dns_endpoint.test.vpc_id]
}

output "is_vpc_id_filter_useful" {
  value = length(local.filter_result_by_vpc_id) > 0 && alltrue(local.filter_result_by_vpc_id)
}

// filter by name
data "huaweicloud_dns_endpoints" "filter_by_name" {
  direction = huaweicloud_dns_endpoint.test.direction
  name      = huaweicloud_dns_endpoint.test.name
}

locals {
  filter_result_by_name = [for v in data.huaweicloud_dns_endpoints.filter_by_name.endpoints[*].name :
    v == huaweicloud_dns_endpoint.test.name]
}

output "is_name_filter_useful" {
  value = length(local.filter_result_by_name) > 0 && alltrue(local.filter_result_by_name)
}
`, testDNSEndpoint_basic(name))
}
