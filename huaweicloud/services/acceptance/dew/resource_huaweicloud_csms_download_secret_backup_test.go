package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCsmsDownloadSecretBackup_basic(t *testing.T) {
	resourceName := "huaweicloud_csms_download_secret_backup.test"

	// lintignore:AT001
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCsmsSecretName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCsmsDownloadSecretBackup_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "secret_blob"),
				),
			},
		},
	})
}

func testAccCsmsDownloadSecretBackup_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_csms_download_secret_backup" "test" {
  secret_name = "%s"
}
`, acceptance.HW_CSMS_SECRET_NAME)
}
