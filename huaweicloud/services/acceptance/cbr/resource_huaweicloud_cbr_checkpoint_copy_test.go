package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Please prepare environment variables using to copy CBR checkpoint.
//
// HW_CBR_VAULT_ID: The ID of the source vault where the backup to be copied is located.
// HW_CBR_DESTINATION_PROJECT_ID: The ID of the destination project to which the backup is to be copied.
// HW_CBR_DESTINATION_REGION: The ID of the destination region to which the backup is to be copied.
// HW_CBR_DESTINATION_VAULT_ID: The ID of the destination vault to which the backup is to be copied.
func TestAccResourceCheckpointCopy_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCbrCheckpointCopy(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceCheckpointCopy_basic(),
			},
		},
	})
}

func testResourceCheckpointCopy_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_checkpoint_copy" "test" {
  vault_id               = "%[1]s"
  destination_project_id = "%[2]s"
  destination_region     = "%[3]s"
  destination_vault_id   = "%[4]s"
  auto_trigger           = true
  enable_acceleration    = false
}
`, acceptance.HW_CBR_VAULT_ID, acceptance.HW_CBR_DESTINATION_PROJECT_ID, acceptance.HW_CBR_DESTINATION_REGION,
		acceptance.HW_CBR_DESTINATION_VAULT_ID)
}
