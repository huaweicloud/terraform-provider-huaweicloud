package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKeypairsDisassociate_basic(t *testing.T) {
	// lintignore:AT001
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare one ECS with KMS keypair, then config it to the environment variable.
			acceptance.TestAccPreCheckKmsKeyPrivateKey(t)
			acceptance.TestAccPreCheckKmsSSHPort(t)
			acceptance.TestAccPreCheckECSID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKeypairsDisassociate_basic(),
			},
		},
	})
}

func testAccKeypairsDisassociate_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_keypair_disassociate" "test" {
  server {
    id   = "%[1]s"
    port = %[2]s

    auth {
      type = "keypair"
      key  = "%[3]s"
    }
  }
}
`, acceptance.HW_ECS_ID, acceptance.HW_KMS_KEYPAIR_SSH_PORT, acceptance.HW_KMS_KEY_PRIVATE_KEY)
}
