package aom

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDashboards_basic(t *testing.T) {
	dataSource := "data.huaweicloud_aom_dashboards.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDashboards_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "dashboards.#"),
					resource.TestCheckResourceAttrSet(dataSource, "dashboards.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "dashboards.0.dashboard_title"),
					resource.TestCheckResourceAttrSet(dataSource, "dashboards.0.dashboard_type"),
					resource.TestCheckResourceAttrSet(dataSource, "dashboards.0.enterprise_project_id"),
				),
			},
		},
	})
}

func testDataSourceDashboards_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_aom_dashboards" "test" {
  depends_on = [huaweicloud_aom_dashboard.test]

  enterprise_project_id = huaweicloud_aom_dashboard.test.enterprise_project_id
  dashboard_type        = huaweicloud_aom_dashboard.test.dashboard_type
}
`, testDashboard_basic(name))
}
