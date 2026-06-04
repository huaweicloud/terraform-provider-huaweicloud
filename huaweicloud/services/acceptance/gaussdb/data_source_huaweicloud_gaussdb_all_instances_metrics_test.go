package gaussdb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussDbAllInstancesMetrics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_all_instances_metrics.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussDbAllInstancesMetrics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.mode"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.engine_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.solution"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.disk_used_size"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.disk_total_size"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.disk_usage"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.p80"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.p95"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.deadlocks"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.buffer_hit_ratio"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.nodes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.nodes.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.nodes.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.nodes.0.role"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.nodes.0.component_ids.#"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceGaussDbAllInstancesMetrics_basic() string {
	return `
data "huaweicloud_gaussdb_all_instances_metrics" "test" {
}

data "huaweicloud_gaussdb_all_instances_metrics" "name_filter" {
  name = "gauss-ee0a"
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_all_instances_metrics.name_filter.instances) > 0
}
`
}
