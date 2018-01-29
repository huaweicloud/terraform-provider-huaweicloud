package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/huawei-clouds/golangsdk/openstack/kms/v1/keys"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccKmsKeyV1_basic(t *testing.T) {
	var key keys.Key
	var keyAlias = fmt.Sprintf("kms_%s", acctest.RandString(5))
	var keyAliasUpdate = fmt.Sprintf("kms_updated_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKmsV1KeyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccKmsV1Key_basic(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsV1KeyExists("huaweicloud_kms_key_v1.key_2", &key),
					resource.TestCheckResourceAttr(
						"huaweicloud_kms_key_v1.key_2", "key_alias", keyAlias),
				),
			},
			resource.TestStep{
				Config: testAccKmsV1Key_update(keyAliasUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsV1KeyExists("huaweicloud_kms_key_v1.key_2", &key),
					resource.TestCheckResourceAttr(
						"huaweicloud_kms_key_v1.key_2", "key_alias", keyAliasUpdate),
					resource.TestCheckResourceAttr(
						"huaweicloud_kms_key_v1.key_2", "key_description", "key update description"),
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
		if rs.Type != "huaweicloud_kms_key_v1" {
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

func TestAccKmsKey_isEnabled(t *testing.T) {
	var key1, key2, key3 keys.Key
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKmsV1KeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKey_enabled(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsV1KeyExists("huaweicloud_kms_key_v1.bar", &key1),
					resource.TestCheckResourceAttr("huaweicloud_kms_key_v1.bar", "is_enabled", "true"),
					testAccCheckKmsKeyIsEnabled(&key1, true),
				),
			},
			{
				Config: testAccKmsKey_disabled(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsV1KeyExists("huaweicloud_kms_key_v1.bar", &key2),
					resource.TestCheckResourceAttr("huaweicloud_kms_key_v1.bar", "is_enabled", "false"),
					testAccCheckKmsKeyIsEnabled(&key2, false),
				),
			},
			{
				Config: testAccKmsKey_enabled(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsV1KeyExists("huaweicloud_kms_key_v1.bar", &key3),
					resource.TestCheckResourceAttr("huaweicloud_kms_key_v1.bar", "is_enabled", "true"),
					testAccCheckKmsKeyIsEnabled(&key3, true),
				),
			},
		},
	})
}

func testAccCheckKmsKeyIsEnabled(key *keys.Key, isEnabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if (key.KeyState == EnabledState) != isEnabled {
			return fmt.Errorf("Expected key %q to have is_enabled=%t, given %t",
				key.KeyID, isEnabled, key.KeyState)
		}

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

func testAccKmsV1Key_update(keyAliasUpdate string) string {
	return fmt.Sprintf(`
		resource "huaweicloud_kms_key_v1" "key_2" {
           key_alias       = "%s"
           key_description = "key update description"
           pending_days = "7"
		}
	`, keyAliasUpdate)
}

func testAccKmsKey_enabled(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key_v1" "bar" {
    key_description = "Terraform acc test is_enabled %s"
    pending_days    = "7"
    key_alias       = "tf-acc-test-kms-key-%s"
}`, rName, rName)
}

func testAccKmsKey_disabled(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key_v1" "bar" {
    key_description = "Terraform acc test is_enabled %s"
    pending_days    = "7"
    key_alias       = "tf-acc-test-kms-key-%s"
    is_enabled      = false
}`, rName, rName)
}
