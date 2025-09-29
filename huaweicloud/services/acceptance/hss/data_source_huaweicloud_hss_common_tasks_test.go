package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHssCommonTasks_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_common_tasks.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHssCommonTasks_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
				),
			},
		},
	})
}

const testAccDataSourceHssCommonTasks_basic = `
data "huaweicloud_hss_common_tasks" "test" {
  task_type             = "cluster_scan"
  enterprise_project_id = "0"
  task_name             = "non-exist-name"
  trigger_type          = "manual"
  task_status           = "finished"
  sort_key              = "start_time"
  sort_dir              = "desc"
}
`
