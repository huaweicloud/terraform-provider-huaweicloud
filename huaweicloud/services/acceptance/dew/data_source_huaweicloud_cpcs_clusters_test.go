package dew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCpcsClusters_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cpcs_clusters.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckCpcsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCpcsClusters_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.cluster_id"),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.cluster_name"),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.service_type"),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.instance_num"),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "clusters.0.az"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_service_type_filter_useful", "true"),
				),
			},
		},
	})
}

const testDataSourceCpcsClusters_basic = `
data "huaweicloud_cpcs_clusters" "test" {}

locals {
  cluster_name = data.huaweicloud_cpcs_clusters.test.clusters.0.cluster_name
  service_type = data.huaweicloud_cpcs_clusters.test.clusters.0.service_type
}

# Filter by cluster name
data "huaweicloud_cpcs_clusters" "name_filter" {
  name = local.cluster_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_cpcs_clusters.name_filter.clusters[*].cluster_name : v == local.cluster_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by service type
data "huaweicloud_cpcs_clusters" "service_type_filter" {
  service_type = local.service_type
}

locals {
  service_type_filter_result = [
    for v in data.huaweicloud_cpcs_clusters.service_type_filter.clusters[*].service_type : v == local.service_type
  ]
}

output "is_service_type_filter_useful" {
  value = length(local.service_type_filter_result) > 0 && alltrue(local.service_type_filter_result)
}
`
