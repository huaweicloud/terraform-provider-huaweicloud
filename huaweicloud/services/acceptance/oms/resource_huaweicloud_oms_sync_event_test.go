package oms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSyncEvent_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOmsSyncTaskId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testSyncEvent_basic(),
			},
		},
	})
}

func testSyncEvent_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_oms_sync_event" "test" {
  sync_task_id = "%s"
  object_keys  = ["test.txt"]
}
`, acceptance.HW_OMS_SYNC_TASK_ID)
}
