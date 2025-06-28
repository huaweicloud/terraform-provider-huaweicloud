package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceReplicateBackup_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare backup ID using to replicate destination backup.
			acceptance.TestAccPreCheckCBRBackupID(t)
			// Please prepare destination project id for replicating.
			acceptance.TestAccPreCheckCBRDestinationProjectID(t)
			// Please prepare destination region name for replicating.
			acceptance.TestAccPreCheckCBRRegionName(t)
			// Please prepare destination vault id for replicating.
			acceptance.TestAccPreCheckImsVaultId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceReplicateBackup_basic(),
			},
		},
	})
}

func testResourceReplicateBackup_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_replicate_backup" "test" {
  backup_id = "%s"
  
  replicate {
    destination_project_id = "%s"
    destination_region     = "%s"
    destination_vault_id   = "%s"
    name                   = "test-replicate-backup"
    description            = "test replicate backup description"
    enable_acceleration    = false
  }
}
`, acceptance.HW_CBR_BACKUP_ID, acceptance.HW_CBR_DESTINATION_PROJECT_ID, acceptance.HW_REGION_NAME_1, acceptance.HW_IMS_VAULT_ID)
}
