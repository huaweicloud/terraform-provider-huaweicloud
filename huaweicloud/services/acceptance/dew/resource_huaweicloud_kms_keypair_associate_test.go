package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKeypairsAssociate_basic(t *testing.T) {
	// lintignore:AT001
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare two valid KMS keypair and one ECS, then config it to the environment variable.
			acceptance.TestAccPreCheckKmsKeyPair(t)
			acceptance.TestAccPreCheckECSAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKeypairsAssociate_basic(),
			},
			{
				Config: testAccKeypairsAssociate_replace(),
			},
		},
	})
}

func testAccKeypairsAssociate_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_keypair_associate" "test" {
  keypair_name = "%[1]s"

  server {
    id   = "%[2]s"
    port = %[3]s

    auth {
      type = "password"
      key  = "%[4]s"
    }
  }
}
`, acceptance.HW_KMS_KEYPAIR_NAME_1, acceptance.HW_ECS_ID, acceptance.HW_KMS_KEYPAIR_SSH_PORT, acceptance.HW_ECS_ROOT_PWD)
}

func testAccKeypairsAssociate_replace() string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_keypair_associate" "test1" {
  keypair_name = "%[1]s"
  
  server {
    id   = "%[2]s"
    port = %[3]s

    auth {
      type = "keypair"
      key  = "%[4]s"
    }
  }
}
`, acceptance.HW_KMS_KEYPAIR_NAME_2, acceptance.HW_ECS_ID, acceptance.HW_KMS_KEYPAIR_SSH_PORT, acceptance.HW_KMS_KEYPAIR_KEY_1)
}
