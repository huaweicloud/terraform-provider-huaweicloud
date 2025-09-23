package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAggregatorDiscoveredResources_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_rms_resource_aggregator_discovered_resources.basic"
	dataSource2 := "data.huaweicloud_rms_resource_aggregator_discovered_resources.filter_by_service_type"
	dataSource3 := "data.huaweicloud_rms_resource_aggregator_discovered_resources.filter_by_resource_type"
	dataSource4 := "data.huaweicloud_rms_resource_aggregator_discovered_resources.filter_by_resource_id"
	rName := acceptance.RandomAccResourceName()
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)
	dc3 := acceptance.InitDataSourceCheck(dataSource3)
	dc4 := acceptance.InitDataSourceCheck(dataSource4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAggregatorDiscoveredResources_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					dc3.CheckResourceExists(),
					dc4.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_service_type_filter_useful", "true"),
					resource.TestCheckOutput("is_resource_type_filter_useful", "true"),
					resource.TestCheckOutput("is_resource_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceAggregatorDiscoveredResources_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_rms_resource_aggregator" "test" {
  name        = "%[1]s"
  type        = "ACCOUNT"
  account_ids = ["%[2]s"]

  depends_on = [huaweicloud_vpc.test]

  # wait 30 seconds to let the aggregator discover resources
  provisioner "local-exec" {
    command = "sleep 30"
  }
}
`, name, acceptance.HW_DOMAIN_ID)
}

func testDataSourceAggregatorDiscoveredResources_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rms_resource_aggregator_discovered_resources" "basic" {
  aggregator_id = huaweicloud_rms_resource_aggregator.test.id
}

data "huaweicloud_rms_resource_aggregator_discovered_resources" "filter_by_service_type" {
  aggregator_id = huaweicloud_rms_resource_aggregator.test.id
  service_type  = "vpc"
}

data "huaweicloud_rms_resource_aggregator_discovered_resources" "filter_by_resource_type" {
  aggregator_id = huaweicloud_rms_resource_aggregator.test.id
  resource_type = "vpcs"
}

data "huaweicloud_rms_resource_aggregator_discovered_resources" "filter_by_resource_id" {
  aggregator_id = huaweicloud_rms_resource_aggregator.test.id

  filter {
    resource_id = huaweicloud_vpc.test.id
  }
}

locals {
  service_type_filter_result = [
    for v in data.huaweicloud_rms_resource_aggregator_discovered_resources.filter_by_service_type.resources[*].service : v == "vpc"
  ]
  resource_type_filter_result = [
    for v in data.huaweicloud_rms_resource_aggregator_discovered_resources.filter_by_resource_type.resources[*].type : v == "vpcs"
  ]
  resource_id_filter_result = [
    for v in data.huaweicloud_rms_resource_aggregator_discovered_resources.filter_by_resource_id.resources[*].resource_id :
    v == huaweicloud_vpc.test.id
  ]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_rms_resource_aggregator_discovered_resources.basic.resources) > 0
}

output "is_service_type_filter_useful" {
  value = alltrue(local.service_type_filter_result) && length(local.service_type_filter_result) > 0
}

output "is_resource_type_filter_useful" {
  value = alltrue(local.resource_type_filter_result) && length(local.resource_type_filter_result) > 0
}

output "is_resource_id_filter_useful" {
  value = alltrue(local.resource_id_filter_result) && length(local.resource_id_filter_result) > 0
}
`, testDataSourceAggregatorDiscoveredResources_base(name), name)
}
