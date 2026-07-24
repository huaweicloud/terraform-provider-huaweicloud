package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscDbMaskTasks_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_db_mask_tasks.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDSCScanTemplateID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDscDbMaskTasks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.db_type"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.execute_line"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.progress"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.run_status"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.task_template_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.type"),
				),
			},
		},
	})
}

func testDataSourceDscDbMaskTasks_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dsc_db_mask_tasks" "test" {
  template_id = "%s"
}
`, acceptance.HW_DSC_SCAN_TEMPLATE_ID)
}
