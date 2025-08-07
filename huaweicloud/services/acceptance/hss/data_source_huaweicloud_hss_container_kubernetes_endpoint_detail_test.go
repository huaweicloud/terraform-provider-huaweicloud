package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceContainerKubernetesEndpointDetail_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_container_kubernetes_endpoint_detail.test"
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
				Config: testDataSourceContainerKubernetesEndpointDetail_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "name"),
					resource.TestCheckResourceAttrSet(dataSource, "service_name"),
					resource.TestCheckResourceAttrSet(dataSource, "namespace"),
					resource.TestCheckResourceAttrSet(dataSource, "creation_timestamp"),
					resource.TestCheckResourceAttrSet(dataSource, "cluster_name"),
					resource.TestCheckResourceAttrSet(dataSource, "labels"),
					resource.TestCheckResourceAttrSet(dataSource, "association_service"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoint_pod_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoint_port_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoint_port_list.0.endpoint_id"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoint_port_list.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoint_port_list.0.protocol"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoint_port_list.0.port"),

					resource.TestCheckOutput("is_endpoint_id_filter_useful", "true"),
				),
			},
		},
	})
}

const testDataSourceContainerKubernetesEndpointDetail_basic string = `
data "huaweicloud_hss_container_kubernetes_endpoints" "test" {}

# Filter using endpoint_id.
locals {
  endpoint_id = data.huaweicloud_hss_container_kubernetes_endpoints.test.endpoint_info_list[0].id
}

data "huaweicloud_hss_container_kubernetes_endpoint_detail" "test" {
  endpoint_id = local.endpoint_id
}

output "is_endpoint_id_filter_useful" {
  value = data.huaweicloud_hss_container_kubernetes_endpoint_detail.test.id == local.endpoint_id
}
`
