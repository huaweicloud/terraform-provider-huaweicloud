package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTaurusDBBackupsBatchDelete_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBBackupId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBBackupsBatchDelete_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("huaweicloud_taurusdb_backups_batch_delete.test", "id"),
					resource.TestCheckResourceAttr("huaweicloud_taurusdb_backups_batch_delete.test", "success_count", "1"),
				),
			},
		},
	})
}

func testAccTaurusDBBackupsBatchDelete_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_taurusdb_backups_batch_delete" "test" {
  backup_ids = ["%s"]
}
`, acceptance.HW_TAURUSDB_BACKUP_ID)
}
