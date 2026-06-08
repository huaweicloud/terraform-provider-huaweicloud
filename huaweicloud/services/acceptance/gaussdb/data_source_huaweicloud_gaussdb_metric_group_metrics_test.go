package gaussdb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussDbMetricGroupMetrics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_metric_group_metrics.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussDbMetricGroupMetrics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "metric_names.#"),
					resource.TestCheckResourceAttrSet(dataSource, "metric_names.0.metric"),
					resource.TestCheckResourceAttrSet(dataSource, "metric_names.0.name"),
				),
			},
		},
	})
}

func testDataSourceGaussDbMetricGroupMetrics_basic() string {
	return `
data "huaweicloud_gaussdb_metric_group_metrics" "test" {
  group_name = "CPUMEMORY"
}
`
}
