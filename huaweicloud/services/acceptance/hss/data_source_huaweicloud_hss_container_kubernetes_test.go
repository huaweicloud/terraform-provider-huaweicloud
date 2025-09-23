package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceContainerKubernetes_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_container_kubernetes.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID with container protection enabled.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceContainerKubernetes_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "total_num"),
					resource.TestCheckResourceAttrSet(dataSource, "last_update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.container_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.container_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.image_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.cpu_limit"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.memory_limit"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.restart_count"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.pod_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.cluster_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.cluster_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.cluster_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.risky"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.low_risk"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.medium_risk"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.high_risk"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.fatal_risk"),

					resource.TestCheckOutput("is_container_name_filter_useful", "true"),
					resource.TestCheckOutput("is_pod_name_filter_useful", "true"),
					resource.TestCheckOutput("is_image_name_filter_useful", "true"),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

const testDataSourceContainerKubernetes_basic string = `
data "huaweicloud_hss_container_kubernetes" "test" {}

# Filter using container_name.
locals {
  container_name = data.huaweicloud_hss_container_kubernetes.test.data_list[0].container_name
}

data "huaweicloud_hss_container_kubernetes" "container_name_filter" {
  container_name = local.container_name
}

output "is_container_name_filter_useful" {
  value = length(data.huaweicloud_hss_container_kubernetes.container_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_container_kubernetes.container_name_filter.data_list[*].container_name : v == local.container_name]
  )
}

# Filter using pod_name.
locals {
  pod_name = data.huaweicloud_hss_container_kubernetes.test.data_list[0].pod_name
}

data "huaweicloud_hss_container_kubernetes" "pod_name_filter" {
  pod_name = local.pod_name
}

output "is_pod_name_filter_useful" {
  value = length(data.huaweicloud_hss_container_kubernetes.pod_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_container_kubernetes.pod_name_filter.data_list[*].pod_name : v == local.pod_name]
  )
}

# Filter using image_name.
locals {
  image_name = data.huaweicloud_hss_container_kubernetes.test.data_list[0].image_name
}

data "huaweicloud_hss_container_kubernetes" "image_name_filter" {
  image_name = local.image_name
}

output "is_image_name_filter_useful" {
  value = length(data.huaweicloud_hss_container_kubernetes.image_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_container_kubernetes.image_name_filter.data_list[*].image_name : v == local.image_name]
  )
}

# Filter using enterprise_project_id.
data "huaweicloud_hss_container_kubernetes" "enterprise_project_id_filter" {
  enterprise_project_id = "all_granted_eps"
}

output "is_enterprise_project_id_filter_useful" {
  value = length(data.huaweicloud_hss_container_kubernetes.enterprise_project_id_filter.data_list) > 0
}

# Filter using non existent container_name.
data "huaweicloud_hss_container_kubernetes" "not_found" {
  container_name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_container_kubernetes.not_found.data_list) == 0
}
`
