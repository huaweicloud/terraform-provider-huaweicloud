package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCssScanTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_css_scan_tasks.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCssScanTasks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "scan_tasks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "scan_tasks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "scan_tasks.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "scan_tasks.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "scan_tasks.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "scan_tasks.0.smn_status"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCssScanTasks_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_css_scan_tasks" "test" {
  depends_on = [huaweicloud_css_scan_task.test]

  cluster_id = huaweicloud_css_cluster.test.id
}
`, testAccScanTask_basic(name))
}
