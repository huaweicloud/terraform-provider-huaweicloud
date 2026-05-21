package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBatchRetryTask_basic(t *testing.T) {
	resourceName := "huaweicloud_drs_batch_retry_task.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testBatchRetryTask_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "results.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.status"),
				),
			},
		},
	})
}

func testBatchRetryTask_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_batch_retry_task" "test" {
  jobs {
    job_id          = "%s"
    is_sync_re_edit = "false"
  }
}
`, acceptance.HW_DRS_JOB_ID)
}
