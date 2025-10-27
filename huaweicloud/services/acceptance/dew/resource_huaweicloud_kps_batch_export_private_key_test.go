package dew

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKpsBatchExportPrivateKey_basic(t *testing.T) {
	exportFileName := fmt.Sprintf("./%s.zip", acceptance.RandomAccResourceName())
	defer func() {
		if err := os.Remove(exportFileName); err != nil {
			log.Printf("error deleting testing file %s: %s", exportFileName, err)
		}
	}()

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please ensure that there is at least one available keypair in the test environment.
			acceptance.TestAccPreCheckKpsEnable(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKpsBatchExportPrivateKey_basic(exportFileName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKpsBatchExportPrivateKeyFileExists(exportFileName),
				),
			},
		},
	})
}

func testAccCheckKpsBatchExportPrivateKeyFileExists(exportFileName string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if _, err := os.Stat(exportFileName); os.IsNotExist(err) {
			return fmt.Errorf("exported file %s does not exist", exportFileName)
		}
		return nil
	}
}

func testAccKpsBatchExportPrivateKey_basic(exportFileName string) string {
	return fmt.Sprintf(`
data "huaweicloud_kps_keypairs" "test" {}

locals {
  name              = data.huaweicloud_kps_keypairs.test.keypairs.0.name
  type              = data.huaweicloud_kps_keypairs.test.keypairs.0.type
  scope             = data.huaweicloud_kps_keypairs.test.keypairs.0.scope
  public_key        = data.huaweicloud_kps_keypairs.test.keypairs.0.public_key
  fingerprint       = data.huaweicloud_kps_keypairs.test.keypairs.0.fingerprint
  is_key_protection = data.huaweicloud_kps_keypairs.test.keypairs.0.is_managed
  frozen_state      = data.huaweicloud_kps_keypairs.test.keypairs.0.frozen_state
}

resource "huaweicloud_kps_batch_export_private_key" "test" {
  export_file_name = "%[1]s"

  keypairs {
    name              = local.name
    type              = local.type
    scope             = local.scope == "account" ? "domain" : local.scope
    public_key        = local.public_key
    fingerprint       = local.fingerprint
    is_key_protection = local.is_key_protection
    frozen_state      = local.frozen_state
  }
}
`, exportFileName)
}
