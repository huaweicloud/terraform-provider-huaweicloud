package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRmsResources_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_rms_resources.basic"
	dataSource2 := "data.huaweicloud_rms_resources.filter_by_tracked"
	dataSource3 := "data.huaweicloud_rms_resources.filter_by_type"
	dataSource4 := "data.huaweicloud_rms_resources.filter_by_region_id"
	dataSource5 := "data.huaweicloud_rms_resources.filter_by_id"
	dataSource6 := "data.huaweicloud_rms_resources.filter_by_epsId"

	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)
	dc3 := acceptance.InitDataSourceCheck(dataSource3)
	dc4 := acceptance.InitDataSourceCheck(dataSource4)
	dc5 := acceptance.InitDataSourceCheck(dataSource5)
	dc6 := acceptance.InitDataSourceCheck(dataSource6)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRmsResources_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					dc3.CheckResourceExists(),
					dc4.CheckResourceExists(),
					dc5.CheckResourceExists(),
					dc6.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_tracked_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					resource.TestCheckOutput("is_region_id_filter_useful", "true"),
					resource.TestCheckOutput("is_resource_id_filter_useful", "true"),
					resource.TestCheckOutput("is_epsId_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceRmsResources_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rms_resources" "basic" {}

data "huaweicloud_rms_resources" "filter_by_tracked" {
  tracked = true
}

data "huaweicloud_rms_resources" "filter_by_type" {
  type = "vpc.vpcs"
}

data "huaweicloud_rms_resources" "filter_by_region_id" {
  region_id = "cn-north-4"
}

locals {
  name             = huaweicloud_vpc.test.name
  id               = huaweicloud_vpc.test.id
  service_result   = [for v in data.huaweicloud_rms_resources.filter_by_type.resources[*].service : v == "vpc"]
  type_result      = [for v in data.huaweicloud_rms_resources.filter_by_type.resources[*].type : v == "vpcs"]
  region_id_result = [for v in data.huaweicloud_rms_resources.filter_by_region_id.resources[*].region_id : v == "cn-north-4"]
}
  
data "huaweicloud_rms_resources" "filter_by_id" {
  resource_id = local.id
  
  depends_on = [huaweicloud_vpc.test]
}

output "is_resource_id_filter_useful" {
  value = length(data.huaweicloud_rms_resources.filter_by_id.resources) > 0 && alltrue([
    for v in data.huaweicloud_rms_resources.filter_by_id.resources[*] : v.id == local.id
  ])
}

data "huaweicloud_rms_resources" "filter_by_epsId" {
  enterprise_project_id = "%[2]s"
  
  depends_on = [huaweicloud_vpc.test]
}

output "is_epsId_filter_useful" {
  value = length(data.huaweicloud_rms_resources.filter_by_epsId.resources) > 0 && alltrue([
    for v in data.huaweicloud_rms_resources.filter_by_epsId.resources[*] : v.enterprise_project_id == "%[2]s"
  ])
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_rms_resources.basic.resources) > 0
}

output "is_tracked_filter_useful" {
  value = length(data.huaweicloud_rms_resources.filter_by_tracked.resources) > 0
}

output "is_type_filter_useful" {
  value = alltrue(local.service_result) && length(local.service_result) > 0 && alltrue(local.type_result) && length(local.type_result) > 0
}

output "is_region_id_filter_useful" {
  value = alltrue(local.region_id_result) && length(local.region_id_result) > 0
}
`, testDataSourceRmsResources_base(), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testDataSourceRmsResources_base() string {
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name                  = "%[1]s"
  cidr                  = "192.168.0.0/16"
  enterprise_project_id = "%[2]s"

  provisioner "local-exec" {
    command = "sleep 60"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
