package dew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceInstances_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cpcs_instances.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running test, prepare a password service cluster instance.
			acceptance.TestAccPrecheckDewFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceInstances_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "result.#"),
					resource.TestCheckResourceAttrSet(dataSource, "result.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "result.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "result.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "result.0.service_type"),
					resource.TestCheckResourceAttrSet(dataSource, "result.0.cluster_id"),
					resource.TestCheckResourceAttrSet(dataSource, "result.0.is_normal"),
					resource.TestCheckResourceAttrSet(dataSource, "result.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "result.0.image_name"),
					resource.TestCheckResourceAttrSet(dataSource, "result.0.specification"),
					resource.TestCheckResourceAttrSet(dataSource, "result.0.az"),
					resource.TestCheckResourceAttrSet(dataSource, "result.0.expired_time"),
					resource.TestCheckResourceAttrSet(dataSource, "result.0.create_time"),

					resource.TestCheckOutput("cluster_id_filter_useful", "true"),
					resource.TestCheckOutput("is_normal_filter_useful", "true"),
				),
			},
		},
	})
}

const testDataSourceInstances_basic = `
data "huaweicloud_cpcs_instances" "test" {}

locals {
  cluster_id = data.huaweicloud_cpcs_instances.test.result.0.cluster_id
}

data "huaweicloud_cpcs_instances" "cluster_id_filter" {
  cluster_id = local.cluster_id
}

output "cluster_id_filter_useful" {
  value = length(data.huaweicloud_cpcs_instances.cluster_id_filter.result) > 0 && alltrue(
    [for v in data.huaweicloud_cpcs_instances.cluster_id_filter.result[*].cluster_id : v == local.cluster_id]
  )
}

data "huaweicloud_cpcs_instances" "is_normal_filter" {
  is_normal = format(data.huaweicloud_cpcs_instances.test.result.0.is_normal)
}

output "is_normal_filter_useful" {
  value = length(data.huaweicloud_cpcs_instances.is_normal_filter.result) > 0
}
`
