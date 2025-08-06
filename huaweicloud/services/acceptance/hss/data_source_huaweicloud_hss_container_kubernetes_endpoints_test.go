package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceContainerKubernetesEndpoints_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_container_kubernetes_endpoints.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires preparing a CCE cluster and synchronizing the cluster information to HSS.
			acceptance.TestAccPreCheckHSSCCEProtection(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceContainerKubernetesEndpoints_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "total_num"),
					resource.TestCheckResourceAttrSet(dataSource, "last_update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoint_info_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoint_info_list.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoint_info_list.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoint_info_list.0.service_name"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoint_info_list.0.namespace"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoint_info_list.0.creation_timestamp"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoint_info_list.0.cluster_name"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoint_info_list.0.cluster_type"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoint_info_list.0.association_service"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_cluster_name_filter_useful", "true"),
					resource.TestCheckOutput("is_namespace_filter_useful", "true"),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

const testDataSourceContainerKubernetesEndpoints_basic string = `
data "huaweicloud_hss_container_kubernetes_endpoints" "test" {}

# Filter using name.
locals {
  name = data.huaweicloud_hss_container_kubernetes_endpoints.test.endpoint_info_list[0].name
}

data "huaweicloud_hss_container_kubernetes_endpoints" "name_filter" {
  name = local.name
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_hss_container_kubernetes_endpoints.name_filter.endpoint_info_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_container_kubernetes_endpoints.name_filter.endpoint_info_list[*].name : v == local.name]
  )
}

# Filter using cluster_name.
locals {
  cluster_name = data.huaweicloud_hss_container_kubernetes_endpoints.test.endpoint_info_list[0].cluster_name
}

data "huaweicloud_hss_container_kubernetes_endpoints" "cluster_name_filter" {
  cluster_name = local.cluster_name
}

output "is_cluster_name_filter_useful" {
  value = length(data.huaweicloud_hss_container_kubernetes_endpoints.cluster_name_filter.endpoint_info_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_container_kubernetes_endpoints.cluster_name_filter.endpoint_info_list[*].cluster_name : v == local.cluster_name]
  )
}

# Filter using namespace.
locals {
  namespace = data.huaweicloud_hss_container_kubernetes_endpoints.test.endpoint_info_list[0].namespace
}

data "huaweicloud_hss_container_kubernetes_endpoints" "namespace_filter" {
  namespace = local.namespace
}

output "is_namespace_filter_useful" {
  value = length(data.huaweicloud_hss_container_kubernetes_endpoints.namespace_filter.endpoint_info_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_container_kubernetes_endpoints.namespace_filter.endpoint_info_list[*].namespace : v == local.namespace]
  )
}

# Filter using enterprise_project_id.
data "huaweicloud_hss_container_kubernetes_endpoints" "enterprise_project_id_filter" {
  enterprise_project_id = "all_granted_eps"
}

output "is_enterprise_project_id_filter_useful" {
  value = length(data.huaweicloud_hss_container_kubernetes_endpoints.enterprise_project_id_filter.endpoint_info_list) > 0
}

# Filter using non existent name.
data "huaweicloud_hss_container_kubernetes_endpoints" "not_found" {
  name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_container_kubernetes_endpoints.not_found.endpoint_info_list) == 0
}
`
