package dew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssociations_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cpcs_associations.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running test, prepare a password service cluster bind an application.
			acceptance.TestAccPrecheckDewFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAssociations_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "result.0.cluster_id"),
					resource.TestCheckResourceAttrSet(dataSource, "result.0.cluster_name"),
					resource.TestCheckResourceAttrSet(dataSource, "result.0.app_id"),
					resource.TestCheckResourceAttrSet(dataSource, "result.0.app_name"),
					resource.TestCheckResourceAttrSet(dataSource, "result.0.vpc_name"),
					resource.TestCheckResourceAttrSet(dataSource, "result.0.subnet_name"),
					resource.TestCheckResourceAttrSet(dataSource, "result.0.cluster_server_type"),
					resource.TestCheckResourceAttrSet(dataSource, "result.0.vpcep_address"),
					resource.TestCheckResourceAttrSet(dataSource, "result.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "result.0.create_time"),

					resource.TestCheckOutput("cluster_id_filter_useful", "true"),
					resource.TestCheckOutput("app_id_filter_useful", "true"),
				),
			},
		},
	})
}

const testDataSourceAssociations_basic = `
data "huaweicloud_cpcs_associations" "test" {}

locals {
  cluster_id = data.huaweicloud_cpcs_associations.test.result.0.cluster_id
  app_id     = data.huaweicloud_cpcs_associations.test.result.0.app_id
}

data "huaweicloud_cpcs_associations" "cluster_id_filter" {
  cluster_id = local.cluster_id
}

output "cluster_id_filter_useful" {
  value = length(data.huaweicloud_cpcs_associations.cluster_id_filter.result) > 0 && alltrue(
    [for v in data.huaweicloud_cpcs_associations.cluster_id_filter.result[*].cluster_id : v == local.cluster_id]
  )
}

data "huaweicloud_cpcs_associations" "app_id_filter" {
  app_id = local.app_id
}

output "app_id_filter_useful" {
  value = length(data.huaweicloud_cpcs_associations.app_id_filter.result) > 0 && alltrue(
    [for v in data.huaweicloud_cpcs_associations.app_id_filter.result[*].app_id : v == local.app_id]
  )
}
`
