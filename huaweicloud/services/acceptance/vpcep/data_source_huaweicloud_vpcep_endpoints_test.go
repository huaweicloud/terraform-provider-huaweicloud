package vpcep

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceVPCEPEndpoints_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_vpcep_endpoints.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceVpcepEndpoints_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "endpoints.0.service_id"),
					resource.TestCheckResourceAttrSet(rName, "endpoints.0.vpc_id"),
					resource.TestCheckResourceAttrSet(rName, "endpoints.0.subnet_id"),
					resource.TestCheckResourceAttrSet(rName, "endpoints.0.enable_dns"),
					resource.TestCheckResourceAttrSet(rName, "endpoints.0.description"),
					resource.TestCheckResourceAttrSet(rName, "endpoints.0.enable_whitelist"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.tags.owner", "tf-acc"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.whitelist.0", "192.168.0.0/24"),

					resource.TestCheckOutput("service_name_filter_is_useful", "true"),

					resource.TestCheckOutput("vpc_id_filter_is_useful", "true"),

					resource.TestCheckOutput("endpoint_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceVpcepEndpoints_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpcep_endpoints" "test" {
  depends_on = [huaweicloud_vpcep_endpoint.test]
}

data "huaweicloud_vpcep_endpoints" "service_name_filter" {
  service_name = data.huaweicloud_vpcep_endpoints.test.endpoints.0.service_name
}

locals {
  service_name = data.huaweicloud_vpcep_endpoints.test.endpoints.0.service_name
}

output "service_name_filter_is_useful" {
  value = length(data.huaweicloud_vpcep_endpoints.service_name_filter.endpoints) > 0 && alltrue(
    [for v in data.huaweicloud_vpcep_endpoints.service_name_filter.endpoints[*].service_name : v == local.service_name]
  )  
}

data "huaweicloud_vpcep_endpoints" "vpc_id_filter" {
  vpc_id = data.huaweicloud_vpcep_endpoints.test.endpoints.0.vpc_id
}
  
locals {
  vpc_id = data.huaweicloud_vpc.myvpc.id
}
  
output "vpc_id_filter_is_useful" {
  value = length(data.huaweicloud_vpcep_endpoints.vpc_id_filter.endpoints) > 0 && alltrue(
	[for v in data.huaweicloud_vpcep_endpoints.vpc_id_filter.endpoints[*].vpc_id : v == local.vpc_id]
  )  
}

data "huaweicloud_vpcep_endpoints" "endpoint_id_filter" {
  endpoint_id = data.huaweicloud_vpcep_endpoints.test.endpoints.0.id
}
  
locals {
  endpoint_id = huaweicloud_vpcep_endpoint.test.id
}
  
output "endpoint_id_filter_is_useful" {
  value = length(data.huaweicloud_vpcep_endpoints.endpoint_id_filter.endpoints) > 0 && alltrue(
    [for v in data.huaweicloud_vpcep_endpoints.endpoint_id_filter.endpoints[*].id : v == local.endpoint_id]
  )  
}
`, testAccVPCEndpoint_Basic(name))
}
