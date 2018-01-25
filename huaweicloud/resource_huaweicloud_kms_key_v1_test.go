package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/kms/v1/keys"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccKmsKeyV1_basic(t *testing.T) {
	var key keys.Key
	var keyAlias = fmt.Sprintf("kms_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKmsV1KeyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccKmsV1Key_basic(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsV1KeyExists("huaweicloud_kms_key_V1.key_2", &key),
					resource.TestCheckResourceAttr(
						"huaweicloud_kms_key_V1.key_2", "key_alias", keyAlias),
				),
			},
			resource.TestStep{
				Config: testAccKmsV1Key_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsV1KeyExists("huaweicloud_kms_key_V1.key_2", &key),
					resource.TestCheckResourceAttr(
						"huaweicloud_kms_key_V1.key_2", "key_alias", "key_2-updated"),
					resource.TestCheckResourceAttr(
						"huaweicloud_kms_key_V1.key_2", "key_description", "key_2-description-updated"),
				),
			},
		},
	})
}

func testAccCheckKmsV1KeyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	kmsClient, err := config.kmsKeyV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack kms client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_kms_key_V1" {
			continue
		}
		v, err := keys.Get(kmsClient, rs.Primary.ID).ExtractKeyInfo()
		if err != nil {
			return err
		}
		if v.KeyState != "4" {
			return fmt.Errorf("key still exists")
		}
	}
	return nil
}

func testAccCheckKmsV1KeyExists(n string, key *keys.Key) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		kmsClient, err := config.kmsKeyV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack kms client: %s", err)
		}
		found, err := keys.Get(kmsClient, rs.Primary.ID).ExtractKeyInfo()
		if err != nil {
			return err
		}
		if found.KeyID != rs.Primary.ID {
			return fmt.Errorf("key not found")
		}

		*key = *found
		return nil
	}
}

func testAccKmsV1Key_basic(keyAlias string) string {
	return fmt.Sprintf(`
		resource "huaweicloud_kms_key_v1" "key_2" {
			key_alias = "%s"

			pending_days = "7"
		}
	`, keyAlias)
}

const testAccKmsV1Key_update = `
resource "huaweicloud_kms_key_v1" "key_2" {
  key_description = "key_2-description-updated"
  key_alias = "key_2-updated"
}
`
