package cts

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCtsOperations_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cts_operations.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCtsOperations_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "operations.0.operation_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "operations.0.service_type"),
					resource.TestCheckResourceAttrSet(dataSource, "operations.0.resource_type"),
					resource.TestCheckOutput("is_default_filter_useful", "true"),
					resource.TestCheckOutput("is_service_type_filter_useful", "true"),
					resource.TestCheckOutput("is_resource_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCtsOperations_basic() string {
	return `
locals {
  service_type  = "CTS"
  resource_type = "tracker"
}

data "huaweicloud_cts_operations" "test" {}

output "is_default_filter_useful" {
  value = length(data.huaweicloud_cts_operations.test.operations) >= 1
}

data "huaweicloud_cts_operations" "filter_by_service_type" {
  service_type = local.service_type
}

output "is_service_type_filter_useful" {
  value = length(data.huaweicloud_cts_operations.filter_by_service_type.operations) >= 1 && alltrue(
    [for op in data.huaweicloud_cts_operations.filter_by_service_type.operations[*] : op.service_type == local.service_type]
  )
}

data "huaweicloud_cts_operations" "filter_by_resource_type" {
  service_type  = local.service_type
  resource_type = local.resource_type
}

output "is_resource_type_filter_useful" {
  value = length(data.huaweicloud_cts_operations.filter_by_resource_type.operations) >= 1 && alltrue(
    [for op in data.huaweicloud_cts_operations.filter_by_resource_type.operations[*] : op.resource_type == local.resource_type]
  )
}
`
}
