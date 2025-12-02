package workspace

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataFlavors_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_workspace_flavors.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByVcpus   = "data.huaweicloud_workspace_flavors.filter_by_vcpus"
		dcFilterByVcpus = acceptance.InitDataSourceCheck(filterByVcpus)

		filterByNotFoundVcpus   = "data.huaweicloud_workspace_flavors.filter_by_not_found_vcpus"
		dcFilterByNotFoundVcpus = acceptance.InitDataSourceCheck(filterByNotFoundVcpus)

		filterByMemory   = "data.huaweicloud_workspace_flavors.filter_by_memory"
		dcFilterByMemory = acceptance.InitDataSourceCheck(filterByMemory)

		filterByOsType   = "data.huaweicloud_workspace_flavors.filter_by_os_type"
		dcFilterByOsType = acceptance.InitDataSourceCheck(filterByOsType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataFlavors_basic,
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "flavors.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by 'vcpus' parameter.
					dcFilterByVcpus.CheckResourceExists(),
					resource.TestCheckOutput("is_vcpus_filter_useful", "true"),
					// Filter by not found vpcus.
					dcFilterByNotFoundVcpus.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found_validation_pass", "true"),
					// Filter by 'memory' parameter.
					dcFilterByMemory.CheckResourceExists(),
					resource.TestCheckOutput("is_memory_filter_useful", "true"),
					// Filter by 'os_type' parameter.
					dcFilterByOsType.CheckResourceExists(),
					resource.TestCheckOutput("is_os_type_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataFlavors_basic = `
// Without any filter parameter.
data "huaweicloud_workspace_flavors" "all" {}

# Filter by 'vcpus' parameter.
locals {
  vcpus = data.huaweicloud_workspace_flavors.all.flavors[0].vcpus
}

data "huaweicloud_workspace_flavors" "filter_by_vcpus" {
  vcpus = local.vcpus
}

locals {
  vcpus_filter_result = [for v in data.huaweicloud_workspace_flavors.filter_by_vcpus.flavors : v.vcpus == local.vcpus]
}

output "is_vcpus_filter_useful" {
  value = alltrue(local.vcpus_filter_result) && length(local.vcpus_filter_result) > 0
}

# Filter by 'vcpus' parameter and with invalid value.
data "huaweicloud_workspace_flavors" "filter_by_not_found_vcpus" {
  vcpus = -1
}

output "is_not_found_validation_pass" {
  value = length(data.huaweicloud_workspace_flavors.filter_by_not_found_vcpus.flavors) == 0
}

# Filter by 'memory' parameter.
locals {
  memory = data.huaweicloud_workspace_flavors.all.flavors[0].memory
}

data "huaweicloud_workspace_flavors" "filter_by_memory" {
  memory = local.memory
}

locals {
  memory_filter_result = [for v in data.huaweicloud_workspace_flavors.filter_by_memory.flavors : v.memory == local.memory]
}

output "is_memory_filter_useful" {
  value = alltrue(local.memory_filter_result) && length(local.memory_filter_result) > 0
}

# Filter by 'os_type' parameter.
locals {
  os_type = "Windows"
}

data "huaweicloud_workspace_flavors" "filter_by_os_type" {
  os_type = local.os_type
}

output "is_os_type_filter_useful" {
  value = length(data.huaweicloud_workspace_flavors.filter_by_os_type.flavors) > 0
}
`
