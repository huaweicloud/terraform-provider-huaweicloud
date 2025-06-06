package dew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccFailedTasksDelete_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare some failed KPS tasks.
			acceptance.TestAccPreCheckKpsEnable(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFailedTasksDelete_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("is_delete_success", "true"),
				),
			},
		},
	})
}

const testAccFailedTasksDelete_basic string = `
data "huaweicloud_kps_failed_tasks" "before_delete" {}

resource "huaweicloud_kps_failed_tasks_delete" "test" {
    depends_on = [data.huaweicloud_kps_failed_tasks.before_delete]
}

data "huaweicloud_kps_failed_tasks" "after_delete" {
    depends_on = [huaweicloud_kps_failed_tasks_delete.test]
}

output "is_delete_success" {
    value = alltrue([length(data.huaweicloud_kps_failed_tasks.before_delete.tasks) > 0, 
    length(data.huaweicloud_kps_failed_tasks.after_delete.tasks) == 0])
}
`
