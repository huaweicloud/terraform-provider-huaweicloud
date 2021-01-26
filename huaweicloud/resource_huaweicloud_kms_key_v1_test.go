package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/kms/v1/keys"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccKmsKeyV1_Basic(t *testing.T) {
	var key keys.Key
	var keyAlias = fmt.Sprintf("kms_%s", acctest.RandString(5))
	var keyAliasUpdate = fmt.Sprintf("kms_updated_%s", acctest.RandString(5))
	var resourceName = "huaweicloud_kms_key.key_2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckKms(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKmsV1KeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsV1Key_Basic(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsV1KeyExists(resourceName, &key),
					resource.TestCheckResourceAttr(resourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(resourceName, "region", HW_REGION_NAME),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"key_usage",
					"pending_days",
				},
			},
			{
				Config: testAccKmsV1KeyUpdate(keyAliasUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsV1KeyExists(resourceName, &key),
					resource.TestCheckResourceAttr(resourceName, "key_alias", keyAliasUpdate),
					resource.TestCheckResourceAttr(resourceName, "key_description", "key update description"),
					resource.TestCheckResourceAttr(resourceName, "region", HW_REGION_NAME),
				),
			},
		},
	})
}

func TestAccKmsKey_WithTags(t *testing.T) {
	var key keys.Key
	var keyAlias = fmt.Sprintf("kms_%s", acctest.RandString(5))
	var resourceName = "huaweicloud_kms_key.key_2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckKms(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKmsV1KeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKey_WithTags(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsV1KeyExists(resourceName, &key),
					resource.TestCheckResourceAttr(resourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
		},
	})
}

func TestAccKmsKeyV1_WithEpsId(t *testing.T) {
	var key keys.Key
	var keyAlias = fmt.Sprintf("kms_%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckKms(t); testAccPreCheckEpsID(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKmsV1KeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsV1Key_epsId(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsV1KeyExists("huaweicloud_kms_key_v1.key_2", &key),
					resource.TestCheckResourceAttr(
						"huaweicloud_kms_key_v1.key_2", "key_alias", keyAlias),
					resource.TestCheckResourceAttr(
						"huaweicloud_kms_key_v1.key_2", "enterprise_project_id", HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccCheckKmsV1KeyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	kmsClient, err := config.kmsKeyV1Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud kms client: %s", err)
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
		kmsClient, err := config.kmsKeyV1Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud kms client: %s", err)
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
	// lintignore:AT009
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckKms(t) },
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
			return fmt.Errorf("Expected key %s to have is_enabled=%t, given %s",
				key.KeyID, isEnabled, key.KeyState)
		}

		return nil
	}
}

func testAccKmsV1Key_Basic(keyAlias string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "key_2" {
  key_alias    = "%s"
  pending_days = "7"
  region       = "%s"
}
`, keyAlias, HW_REGION_NAME)
}

func testAccKmsKey_WithTags(keyAlias string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "key_2" {
  key_alias    = "%s"
  pending_days = "7"
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, keyAlias)
}

func testAccKmsV1Key_epsId(keyAlias string) string {
	return fmt.Sprintf(`
		resource "huaweicloud_kms_key_v1" "key_2" {
			key_alias = "%s"
			pending_days = "7"
			enterprise_project_id = "%s"
		}
	`, keyAlias, HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccKmsV1KeyUpdate(keyAliasUpdate string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "key_2" {
  key_alias       = "%s"
  key_description = "key update description"
  pending_days    = "7"
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
