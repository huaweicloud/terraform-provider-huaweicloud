package dws

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcStatistics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dws_statistics.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceStatistics_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_cluster_total", "true"),
					resource.TestCheckOutput("is_cluster_normal", "true"),
					resource.TestCheckOutput("is_instance_total", "true"),
					resource.TestCheckOutput("is_instance_normal", "true"),
					resource.TestCheckOutput("is_storage", "true"),
				),
			},
		},
	})
}

const testDataSourceStatistics_basic = `
data "huaweicloud_dws_statistics" "test" {}

locals {
  statistics = data.huaweicloud_dws_statistics.test.statistics
  storage    = try([for v in local.statistics : v if v.name == "storage.total"][0], {})
}

output "is_cluster_total" {
  value = try([for v in local.statistics : v.value if v.name == "cluster.total"][0] >= 1, false)
}

output "is_cluster_normal" {
  value = try([for v in local.statistics : v.value if v.name == "cluster.normal"][0] >= 1, false)
}

output "is_instance_total" {
  value = try([for v in local.statistics : v.value if v.name == "instance.total"][0] >= 1, false)
}

output "is_instance_normal" {
  value = try([for v in local.statistics : v.value if v.name == "instance.normal"][0] >= 1, false)
}

output "is_storage" {
  value = try(local.storage.value > 0 && local.storage.unit != "", false)
}
`
