package workspace

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceFlavors_basic(t *testing.T) {
	var (
		byVcpus           = "data.huaweicloud_workspace_flavors.filter_by_vcpus"
		vcpusNotFound     = "data.huaweicloud_workspace_flavors.not_found"
		dcByVcpus         = acceptance.InitDataSourceCheck(byVcpus)
		dcServiceNotFound = acceptance.InitDataSourceCheck(vcpusNotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFlavors_basic,
				Check: resource.ComposeTestCheckFunc(
					dcByVcpus.CheckResourceExists(),
					resource.TestCheckOutput("internal_charging_mode_useful", "true"),
					resource.TestCheckOutput("internal_status_useful", "true"),
					resource.TestCheckOutput("is_vcpus_filter_useful", "true"),
					dcServiceNotFound.CheckResourceExists(),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

const testAccDataSourceFlavors_basic = `
data "huaweicloud_workspace_flavors" "filter_by_vcpus" {
  vcpus = 2
}

locals {
  vcpus_filter_result = [for v in data.huaweicloud_workspace_flavors.filter_by_vcpus.flavors : v.vcpus == 2]
}

output "is_vcpus_filter_useful" {
  value = alltrue(local.vcpus_filter_result) && length(local.vcpus_filter_result) > 0
}

locals {
  charging_mode_filter_result = [for v in data.huaweicloud_workspace_flavors.filter_by_vcpus.flavors : v.charging_mode == "postPaid"]
}

output "internal_charging_mode_useful" {
  value = alltrue(local.charging_mode_filter_result) && length(local.charging_mode_filter_result) > 0
}

locals {
  status_filter_result = [for v in data.huaweicloud_workspace_flavors.filter_by_vcpus.flavors : v.status == "normal"]
}

output "internal_status_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}

data "huaweicloud_workspace_flavors" "not_found" {
  vcpus = -1
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_workspace_flavors.not_found.flavors) == 0
}`

func TestAccDataSourceFlavors_filterByMemory(t *testing.T) {
	var (
		byMemory         = "data.huaweicloud_workspace_flavors.filter_by_memory"
		memoryNotFound   = "data.huaweicloud_workspace_flavors.not_found"
		dcByMemory       = acceptance.InitDataSourceCheck(byMemory)
		dcMemoryNotFound = acceptance.InitDataSourceCheck(memoryNotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFlavors_filterByMemory,
				Check: resource.ComposeTestCheckFunc(
					dcByMemory.CheckResourceExists(),
					resource.TestCheckOutput("is_memory_filter_useful", "true"),
					dcMemoryNotFound.CheckResourceExists(),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

const testAccDataSourceFlavors_filterByMemory = `
data "huaweicloud_workspace_flavors" "filter_by_memory" {
  memory = 4
}

locals {
  memory_filter_result = [for v in data.huaweicloud_workspace_flavors.filter_by_memory.flavors : v.memory == 4]
}

output "is_memory_filter_useful" {
  value = alltrue(local.memory_filter_result) && length(local.memory_filter_result) > 0
}

data "huaweicloud_workspace_flavors" "not_found" {
  memory = -5
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_workspace_flavors.not_found.flavors) == 0
}`

func TestAccDataSourceFlavors_filterByOsType(t *testing.T) {
	var (
		byWindows   = "data.huaweicloud_workspace_flavors.filter_by_os_type"
		dcByWindows = acceptance.InitDataSourceCheck(byWindows)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFlavors_filterByOsType,
				Check: resource.ComposeTestCheckFunc(
					dcByWindows.CheckResourceExists(),
					resource.TestCheckOutput("is_os_type_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceFlavors_filterByOsType = `
data "huaweicloud_workspace_flavors" "filter_by_os_type" {
  os_type = "Windows"
}

output "is_os_type_filter_useful" {
  value = length(data.huaweicloud_workspace_flavors.filter_by_os_type.flavors) > 0
}`

func TestAccDataSourceFlavors_filterByAvailabilityZone(t *testing.T) {
	var (
		byAZone   = "data.huaweicloud_workspace_flavors.filter_by_availability_zone"
		dcByAZone = acceptance.InitDataSourceCheck(byAZone)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFlavors_filterByAvailabilityZone,
				Check: resource.ComposeTestCheckFunc(
					dcByAZone.CheckResourceExists(),
					resource.TestCheckOutput("is_availability_zone_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceFlavors_filterByAvailabilityZone = `
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_workspace_flavors" "filter_by_availability_zone" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

output "is_availability_zone_filter_useful" {
  value = length(data.huaweicloud_workspace_flavors.filter_by_availability_zone.flavors) > 0
}`
