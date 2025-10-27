package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKpsBatchImportKeypair_basic(t *testing.T) {
	rName := "huaweicloud_kps_batch_import_keypair.test"
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare public key and private key file path
			acceptance.TestAccPreCheckKpsKeyFilePath(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKpsBatchImportKeypair_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "succeeded_keypairs.#"),
					resource.TestCheckResourceAttrSet(rName, "succeeded_keypairs.0.name"),
					resource.TestCheckResourceAttrSet(rName, "succeeded_keypairs.0.type"),
					resource.TestCheckResourceAttrSet(rName, "succeeded_keypairs.0.public_key"),
					resource.TestCheckResourceAttrSet(rName, "succeeded_keypairs.0.private_key"),
					resource.TestCheckResourceAttrSet(rName, "succeeded_keypairs.0.fingerprint"),
					resource.TestCheckResourceAttrSet(rName, "succeeded_keypairs.0.user_id"),
				),
			},
		},
	})
}

func testAccKpsBatchImportKeypair_basic() string {
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
resource "huaweicloud_kps_batch_import_keypair" "test" {
  keypairs {
    name       = "%[1]s"
    type       = "ssh"
    public_key = file("%[2]s")
    scope      = "domain"

    key_protection {
      private_key = file("%[3]s")
      encryption {
        type = "default"
      }
    }
  }
}
`, name, acceptance.HW_KPS_PUBLIC_KEY_FILE_PATH, acceptance.HW_KPS_PRIVATE_KEY_FILE_PATH)
}
