package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// This resource can only restore secrets that no longer exist.
func TestAccCsmsRestoreSecret_basic(t *testing.T) {
	resourceName := "huaweicloud_csms_restore_secret.test"

	// lintignore:AT001
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCsmsSecretBackupFilePath(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCsmsRestoreSecret_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "kms_key_id"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttrSet(resourceName, "update_time"),
					resource.TestCheckResourceAttrSet(resourceName, "secret_type"),
					resource.TestCheckResourceAttrSet(resourceName, "auto_rotation"),
					resource.TestCheckResourceAttrSet(resourceName, "enterprise_project_id"),
				),
			},
		},
	})
}

func testAccCsmsRestoreSecret_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_csms_restore_secret" "test" {
  secret_blob = file("%s")
}
`, acceptance.HW_CSMS_SECRET_BACKUP_FILE_PATH)
}
