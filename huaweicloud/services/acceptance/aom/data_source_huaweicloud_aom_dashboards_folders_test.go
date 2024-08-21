package aom

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDashboardsFolders_basic(t *testing.T) {
	dataSource := "data.huaweicloud_aom_dashboards_folders.test"
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
				Config: testDataSourceDashboardsFolders_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "folders.#"),
					resource.TestCheckResourceAttrSet(dataSource, "folders.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "folders.0.folder_title"),
					resource.TestCheckResourceAttrSet(dataSource, "folders.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "folders.0.is_template"),
					resource.TestCheckResourceAttrSet(dataSource, "folders.0.created_by"),
				),
			},
		},
	})
}

func testDataSourceDashboardsFolders_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_aom_dashboards_folders" "test" {
  depends_on = [huaweicloud_aom_dashboards_folder.test]

  enterprise_project_id = huaweicloud_aom_dashboards_folder.test.enterprise_project_id
}
`, testDashboardsFolder_basic(name, false))
}
