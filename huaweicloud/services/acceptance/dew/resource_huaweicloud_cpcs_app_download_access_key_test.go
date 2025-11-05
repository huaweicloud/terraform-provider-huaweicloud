package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Currently, this resource is valid only in cn-north-9 region.
// Please note that each key can only be downloaded once.
func TestAccCpcsAppDownloadAccessKey_basic(t *testing.T) {
	rName := "huaweicloud_cpcs_app_download_access_key.test"
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare a valid app ID and access key ID in the environment variables.
			acceptance.TestAccPrecheckCpcsAppId(t)
			acceptance.TestAccPrecheckCpcsAccessKeyId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testCpcsAppDownloadAccessKey_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "access_key"),
					resource.TestCheckResourceAttrSet(rName, "secret_key"),
					resource.TestCheckResourceAttrSet(rName, "key_name"),
					resource.TestCheckResourceAttrSet(rName, "is_imported"),
				),
			},
		},
	})
}

func testCpcsAppDownloadAccessKey_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cpcs_app_download_access_key" "test" {
  app_id        = "%s"
  access_key_id = "%s"
}
`, acceptance.HW_CPCS_APP_ID, acceptance.HW_CPCS_ACCESS_KEY_ID)
}
