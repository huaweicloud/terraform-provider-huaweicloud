package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCCEReleases_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cce_releases.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCceChartPath(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCceReleases_basic(rName),
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

func testDataSourceDataSourceCceReleases_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cce_releases" "test" {
  depends_on = [huaweicloud_cce_release.test]

  cluster_id = huaweicloud_cce_cluster.test.id
}

data "huaweicloud_cce_releases" "chart_id_filter" {
  depends_on = [huaweicloud_cce_release.test]

  cluster_id = huaweicloud_cce_cluster.test.id
  chart_id   = huaweicloud_cce_chart.test.id
}
locals {
  chart_name = huaweicloud_cce_chart.test.name
}
output "chart_id_filter_is_useful" {
  value = length(data.huaweicloud_cce_releases.chart_id_filter.releases) > 0 && alltrue(
    [for v in data.huaweicloud_cce_releases.chart_id_filter.releases[*].chart_name : v == local.chart_name]
  )
}

data "huaweicloud_cce_releases" "namespace_filter" {
  depends_on = [huaweicloud_cce_release.test]

  cluster_id = huaweicloud_cce_cluster.test.id
  namespace  = huaweicloud_cce_release.test.namespace
}
locals {
  namespace = huaweicloud_cce_release.test.namespace
}
output "namespace_filter_is_useful" {
  value = length(data.huaweicloud_cce_releases.namespace_filter.releases) > 0 && alltrue(
    [for v in data.huaweicloud_cce_releases.namespace_filter.releases[*].namespace : v == local.namespace]
  )
}
`, testAccRelease_basic(name))
}
