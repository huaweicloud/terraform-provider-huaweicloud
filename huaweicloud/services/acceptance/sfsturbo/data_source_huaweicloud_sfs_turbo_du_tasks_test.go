package sfsturbo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDuTasks_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_sfs_turbo_du_tasks.test"
		rName      = acceptance.RandomAccResourceName()
		path       = "/temp" + acctest.RandString(5)
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDuTasks_basic(rName, path),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.begin_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.end_time"),

					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_exist_task", "true"),
				),
			},
		},
	})
}

func testDataSourceDuTasks_basic(name, path string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_sfs_turbo_du_tasks" "test" {
  depends_on = [
    huaweicloud_sfs_turbo_du_task.test
  ]

  share_id = huaweicloud_sfs_turbo.test.id
}

output "is_exist_task" {
  value = length(data.huaweicloud_sfs_turbo_du_tasks.test.tasks) > 0
}
`, testAccDuTask_basic(name, path))
}
