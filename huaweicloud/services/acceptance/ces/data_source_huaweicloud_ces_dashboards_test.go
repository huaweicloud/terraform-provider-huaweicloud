package ces

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCesDashboards_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ces_dashboards.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCesDashboards_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "dashboards.0.dashboard_id"),
					resource.TestCheckResourceAttrSet(dataSource, "dashboards.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "dashboards.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "dashboards.0.row_widget_num"),
					resource.TestCheckResourceAttrSet(dataSource, "dashboards.0.creator_name"),
					resource.TestMatchResourceAttr(dataSource,
						"dashboards.0.created_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckOutput("is_default_filter_useful", "true"),
					resource.TestCheckOutput("is_filter_by_name_useful", "true"),
					resource.TestCheckOutput("is_filter_by_favorite_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCesDashboards_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  name                  = "%[2]s1"
  enterprise_project_id = "0"
  is_favorite           = true
}

data "huaweicloud_ces_dashboards" "test" {
  depends_on = [
    huaweicloud_ces_dashboard.d1,
    huaweicloud_ces_dashboard.d2,
  ]
}

output "is_default_filter_useful" {
  value = length(data.huaweicloud_ces_dashboards.test.dashboards) >= 2
}

data "huaweicloud_ces_dashboards" "filter_by_name" {
  name = local.name

  depends_on = [
    huaweicloud_ces_dashboard.d1,
    huaweicloud_ces_dashboard.d2,
  ]
}

output "is_filter_by_name_useful" {
  value = length(data.huaweicloud_ces_dashboards.filter_by_name.dashboards) >= 1 && alltrue(
    [for r in data.huaweicloud_ces_dashboards.filter_by_name.dashboards[*]: r.name == local.name]
  )
}

data "huaweicloud_ces_dashboards" "filter_by_favorite" {
  enterprise_project_id = local.enterprise_project_id
  is_favorite           = local.is_favorite

  depends_on = [
    huaweicloud_ces_dashboard.d1,
    huaweicloud_ces_dashboard.d2,
  ]
}

output "is_filter_by_favorite_useful" {
  value = length(data.huaweicloud_ces_dashboards.filter_by_favorite.dashboards) >= 1 && alltrue(
    [for r in data.huaweicloud_ces_dashboards.filter_by_favorite.dashboards[*]: r.is_favorite == local.is_favorite]
  )
}
`, testDataSourceCesDashboards_base(name), name)
}

func testDataSourceCesDashboards_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ces_dashboard" "d1" {
  name           = "%[1]s1"
  row_widget_num = 1
  is_favorite    = true
}

resource "huaweicloud_ces_dashboard" "d2" {
  name           = "%[1]s2"
  row_widget_num = 2
}
`, name)
}
