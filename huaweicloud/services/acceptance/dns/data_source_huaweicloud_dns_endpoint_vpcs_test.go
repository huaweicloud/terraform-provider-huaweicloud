package dns

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDNSEndpointVpcs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dns_endpoint_vpcs.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDNSEndpointVpcs_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "vpcs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "vpcs.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "vpcs.0.inbound_endpoint_count"),
					resource.TestCheckResourceAttrSet(dataSource, "vpcs.0.outbound_endpoint_count"),

					resource.TestCheckOutput("is_vpc_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDNSEndpointVpcs_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dns_endpoint_vpcs" "test" {
  depends_on = [huaweicloud_dns_endpoint.test]
}

// filter by vpc_id
data "huaweicloud_dns_endpoint_vpcs" "filter_by_vpc_id" {
  vpc_id = huaweicloud_dns_endpoint.test.vpc_id
}

locals {
  filter_result_by_vpc_id = [for v in data.huaweicloud_dns_endpoint_vpcs.filter_by_vpc_id.vpcs[*].id :
    v == huaweicloud_dns_endpoint.test.vpc_id]
}

output "is_vpc_id_filter_useful" {
  value = length(local.filter_result_by_vpc_id) > 0 && alltrue(local.filter_result_by_vpc_id)
}
`, testDNSEndpoint_basic(name))
}
