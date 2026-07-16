package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGaussDBDownloadBackupsAuthorize_basic(t *testing.T) {
	resourceName := "huaweicloud_gaussdb_download_backups_authorize.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBBackupId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBDownloadBackupsAuthorize_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "bucket"),
					resource.TestCheckResourceAttrSet(resourceName, "file_paths.#"),
				),
			},
		},
	})
}

func testAccGaussDBDownloadBackupsAuthorize_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_download_backups_authorize" "test" {
  backup_id = "%s"
}
`, acceptance.HW_GAUSSDB_BACKUP_ID)
}
