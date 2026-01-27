package cceautopilot

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceReleases_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cce_autopilot_releases.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCceClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceReleases_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "releases.#"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.chart_name"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.chart_public"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.chart_version"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.cluster_id"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.cluster_name"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.create_at"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.namespace"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.status_description"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.update_at"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.values"),
					resource.TestCheckResourceAttrSet(dataSource, "releases.0.version"),
					resource.TestCheckOutput("chart_id_filter_is_useful", "true"),
					resource.TestCheckOutput("namespace_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceReleases_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cce_autopilot_charts" "test" {}

data "huaweicloud_cce_autopilot_releases" "test" {
  cluster_id = "%[1]s"
}

data "huaweicloud_cce_autopilot_releases" "chart_id_filter" {
  cluster_id = "%[1]s"
  chart_id   = data.huaweicloud_cce_autopilot_charts.test.charts[0].id
}
locals {
  chart_name = data.huaweicloud_cce_autopilot_charts.test.charts[0].name
}
output "chart_id_filter_is_useful" {
  value = length(data.huaweicloud_cce_autopilot_releases.chart_id_filter.releases) > 0 && alltrue(
    [for v in data.huaweicloud_cce_autopilot_releases.chart_id_filter.releases[*].chart_name : v == local.chart_name]
  )
}

data "huaweicloud_cce_autopilot_releases" "namespace_filter" {
  cluster_id = "%[1]s"
  namespace  = "monitoring"
}
locals {
  namespace = "monitoring"
}
output "namespace_filter_is_useful" {
  value = length(data.huaweicloud_cce_autopilot_releases.namespace_filter.releases) > 0 && alltrue(
    [for v in data.huaweicloud_cce_autopilot_releases.namespace_filter.releases[*].namespace : v == "monitoring"]
  )
}
`, acceptance.HW_CCE_CLUSTER_ID)
}
