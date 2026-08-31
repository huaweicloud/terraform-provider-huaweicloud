package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBackupBatchDelete_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccCheckGeminidbBackupID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testBackupBatchDelete_basic(),
			},
		},
	})
}

func testBackupBatchDelete_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_geminidb_backup_batch_delete" "test" {
  backup_ids = split(",", "%s")
}
`, acceptance.HW_GEMINIDB_BACKUP_ID)
}
