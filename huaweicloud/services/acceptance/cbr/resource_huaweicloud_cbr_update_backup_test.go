package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceUpdateBackup_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCBRBackupID(t)
			acceptance.TestAccPreCheckCBRBackupName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceUpdateBackup_basic(),
			},
		},
	})
}

func testResourceUpdateBackup_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_update_backup" "test" {
  backup_id = "%s"
  name      = "%s"
}
`, acceptance.HW_CBR_BACKUP_ID, acceptance.HW_CBR_BACKUP_NAME)
}
