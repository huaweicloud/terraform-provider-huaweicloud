package vpcep

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceVpcepServices_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_vpcep_services.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceVpcepServices_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "endpoint_services.0.id"),
					resource.TestCheckResourceAttrSet(rName, "endpoint_services.0.service_name"),
					resource.TestCheckResourceAttrSet(rName, "endpoint_services.0.server_type"),
					resource.TestCheckResourceAttrSet(rName, "endpoint_services.0.vpc_id"),
					resource.TestCheckResourceAttrSet(rName, "endpoint_services.0.approval_enabled"),
					resource.TestCheckResourceAttrSet(rName, "endpoint_services.0.status"),
					resource.TestCheckResourceAttrSet(rName, "endpoint_services.0.service_type"),
					resource.TestCheckResourceAttrSet(rName, "endpoint_services.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "endpoint_services.0.updated_at"),
					resource.TestCheckResourceAttrSet(rName, "endpoint_services.0.connection_count"),
					resource.TestCheckResourceAttrSet(rName, "endpoint_services.0.tcp_proxy"),
					resource.TestCheckResourceAttrSet(rName, "endpoint_services.0.description"),
					resource.TestCheckResourceAttrSet(rName, "endpoint_services.0.public_border_group"),
					resource.TestCheckResourceAttrSet(rName, "endpoint_services.0.enable_policy"),

					resource.TestCheckOutput("service_name_filter_is_useful", "true"),

					resource.TestCheckOutput("server_type_filter_is_useful", "true"),

					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceVpcepServices_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpcep_services" "test" {
  depends_on = [huaweicloud_vpcep_service.test]
}

data "huaweicloud_vpcep_services" "service_name_filter" {
  service_name = data.huaweicloud_vpcep_services.test.endpoint_services.0.service_name
}

locals {
  service_name = data.huaweicloud_vpcep_services.test.endpoint_services.0.service_name
}

output "service_name_filter_is_useful" {
  value = length(data.huaweicloud_vpcep_services.service_name_filter.endpoint_services) > 0 && alltrue(
    [for v in data.huaweicloud_vpcep_services.service_name_filter.endpoint_services[*].service_name : v == local.service_name]
  )  
}

data "huaweicloud_vpcep_services" "server_type_filter" {
  server_type  = data.huaweicloud_vpcep_services.test.endpoint_services.0.server_type
}

locals {
  server_type = data.huaweicloud_vpcep_services.test.endpoint_services.0.server_type
}

output "server_type_filter_is_useful" {
  value = length(data.huaweicloud_vpcep_services.server_type_filter.endpoint_services) > 0 && alltrue(
    [for v in data.huaweicloud_vpcep_services.server_type_filter.endpoint_services[*].server_type : v == local.server_type]
  )  
}

data "huaweicloud_vpcep_services" "status_filter" {
  status  = data.huaweicloud_vpcep_services.test.endpoint_services.0.status
}
  
locals {
  status = data.huaweicloud_vpcep_services.test.endpoint_services.0.status
}
  
output "status_filter_is_useful" {
  value = length(data.huaweicloud_vpcep_services.status_filter.endpoint_services) > 0 && alltrue(
    [for v in data.huaweicloud_vpcep_services.status_filter.endpoint_services[*].status : v == local.status]
  )  
}
`, testAccVPCEPService_Basic(name))
}
