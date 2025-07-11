package codeartsbuild

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBuildLogDownload_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testBuildLogDownload_base(name),
				// wait seconds for delay
				Check: func(*terraform.State) error {
					// lintignore:R018
					time.Sleep(10 * time.Second)
					return nil
				},
			},
			{
				Config: testBuildLogDownload_basic(name),
			},
		},
	})
}

func testBuildLogDownload_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_build_tasks" "test" {
  depends_on = [huaweicloud_codearts_build_task_action.stop]

  project_id = huaweicloud_codearts_project.test.id
}

data "huaweicloud_codearts_build_task_records" "test" {
  build_project_id = data.huaweicloud_codearts_build_tasks.test.tasks[0].build_project_id
}
`, testBuildTaskAction_basic(name))
}

func testBuildLogDownload_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_build_log_download" "test" {
  record_id = data.huaweicloud_codearts_build_task_records.test.records[0].id
  log_file  = "./buildLog.txt"
}
`, testBuildLogDownload_base(name))
}
