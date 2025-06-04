package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceRestore_ECS_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare environment ID using to restore ECS backup.
			acceptance.TestAccPreCheckECSBackupRestore(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceECSRestore_basic(),
			},
		},
	})
}

func testResourceECSRestore_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_restore" "test" {
  backup_id = "%[1]s"
  server_id = "%[2]s"
  power_on  = true

  mappings {
    backup_id = "%[3]s"
    volume_id = "%[4]s"
  }
}
`, acceptance.HW_CBR_ECS_BACKUP_ID,
		acceptance.HW_CBR_ECS_SERVER_ID,
		acceptance.HW_CBR_EVS_BACKUP_ID_FOR_ECS,
		acceptance.HW_CBR_EVS_VOLUME_ID_FOR_ECS,
	)
}

func TestAccResourceRestore_EVS_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare environment ID using to restore EVS backup.
			acceptance.TestAccPreCheckEVSBackupRestore(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceEVSRestore_basic(),
			},
		},
	})
}

func testResourceEVSRestore_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_restore" "test" {
  backup_id = "%[1]s"
  volume_id = "%[2]s"
}
`, acceptance.HW_CBR_EVS_BACKUP_ID, acceptance.HW_CBR_EVS_VOLUME_ID)
}

func TestAccResourceRestore_Workspace_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare environment ID using to restore Workspace backup.
			acceptance.TestAccPreCheckWorkspaceBackupRestore(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceWorkspaceRestore_basic(),
			},
		},
	})
}

func testResourceWorkspaceRestore_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_restore" "test" {
  backup_id   = "%[1]s"
  resource_id = "%[2]s"
  power_on    = true
}
`, acceptance.HW_CBR_WORKSPACE_BACKUP_ID, acceptance.HW_CBR_WORKSPACE_RESOURCE_ID)
}
