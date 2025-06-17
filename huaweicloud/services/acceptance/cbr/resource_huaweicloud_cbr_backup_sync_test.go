package cbr

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Because there is a lack of scenarios for testing the API, the test case only tests one failure error.
func TestAccResourceBackupSync_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testResourceBackupSync_basic,
				ExpectError: regexp.MustCompile(`error creating CBR backup sync`),
			},
		},
	})
}

const testResourceBackupSync_basic = `
resource "huaweicloud_cbr_backup_sync" "test" {
  backup_id     = "not-exist-backup-id"
  resource_id   = "not-exist-resource-id"
  resource_name = "not-exist-resource-name"
  resource_type = "OS::Native::Server"
  bucket_name   = "not-exist-bucket-name"
  image_path    = "not-exist-image-path"
  created_at    = 1553587260
  backup_name   = "not-exist-backup-name"
}
`
