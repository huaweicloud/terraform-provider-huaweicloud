package tms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceResourceTypes_basic(t *testing.T) {
	var (
		byServiceName     = "data.huaweicloud_tms_resource_types.filter_by_service_name"
		serviceNotFound   = "data.huaweicloud_tms_resource_types.not_found"
		dcByServiceName   = acceptance.InitDataSourceCheck(byServiceName)
		dcServiceNotFound = acceptance.InitDataSourceCheck(serviceNotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceResourceTypes_basic,
				Check: resource.ComposeTestCheckFunc(
					dcByServiceName.CheckResourceExists(),
					resource.TestCheckOutput("is_service_name_filter_useful", "true"),
					dcServiceNotFound.CheckResourceExists(),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

const testAccDataSourceResourceTypes_basic = `
data "huaweicloud_tms_resource_types" "filter_by_service_name" {
  service_name = "dli"
}

data "huaweicloud_tms_resource_types" "not_found" {
  service_name = "not_found"
}

locals {
  filter_result = [for v in data.huaweicloud_tms_resource_types.filter_by_service_name.types[*].service_name : v == "dli"]
}

output "is_service_name_filter_useful" {
  value = alltrue(local.filter_result) && length(local.filter_result) > 0
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_tms_resource_types.not_found.types) == 0
}
`

func TestAccDataSourceResourceTypes_filterByRegion(t *testing.T) {
	var (
		byRegion         = "data.huaweicloud_tms_resource_types.filter_by_region"
		regionNotFound   = "data.huaweicloud_tms_resource_types.not_found"
		dcByRegion       = acceptance.InitDataSourceCheck(byRegion)
		dcRegionNotFound = acceptance.InitDataSourceCheck(regionNotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceResourceTypes_filterByRegion(),
				Check: resource.ComposeTestCheckFunc(
					dcByRegion.CheckResourceExists(),
					resource.TestCheckOutput("is_region_filter_useful", "true"),
					dcRegionNotFound.CheckResourceExists(),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccDataSourceResourceTypes_filterByRegion() string {
	return fmt.Sprintf(`
data "huaweicloud_tms_resource_types" "filter_by_region" {
  region = "%[1]s"
}

data "huaweicloud_tms_resource_types" "not_found" {
  region = "not_found"
}

output "is_region_filter_useful" {
  value = length(data.huaweicloud_tms_resource_types.filter_by_region.types) > 0
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_tms_resource_types.not_found.types) == 0
}
`, acceptance.HW_REGION_NAME)
}
