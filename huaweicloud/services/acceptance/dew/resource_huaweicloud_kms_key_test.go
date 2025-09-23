package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/kms/v1/keys"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dew"
)

func getKmsKeyResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.KmsKeyV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating KMS client: %s", err)
	}
	key, err := keys.Get(client, state.Primary.ID).ExtractKeyInfo()

	if err == nil && key.KeyState == dew.PendingDeletionState {
		return nil, golangsdk.ErrDefault404{}
	}
	return key, err
}

// keystore_id scenario testing is currently not supported.
func TestAccKmsKey_basic(t *testing.T) {
	var (
		name         = acceptance.RandomAccResourceName()
		updateName   = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_kms_key.test"
		key          keys.Key
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&key,
		getKmsKeyResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckKms(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKey_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "key_alias", name),
					resource.TestCheckResourceAttr(resourceName, "key_algorithm", "AES_256"),
					resource.TestCheckResourceAttr(resourceName, "key_usage", "ENCRYPT_DECRYPT"),
					resource.TestCheckResourceAttr(resourceName, "key_description", "test acc"),
					resource.TestCheckResourceAttr(resourceName, "origin", "kms"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(resourceName, "key_id"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_id"),
					resource.TestCheckResourceAttrSet(resourceName, "keystore_id"),
					resource.TestCheckResourceAttrSet(resourceName, "is_enabled"),
					resource.TestCheckResourceAttrSet(resourceName, "rotation_enabled"),
					resource.TestCheckResourceAttrSet(resourceName, "key_state"),
					resource.TestCheckResourceAttrSet(resourceName, "default_key_flag"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_date"),
				),
			},
			{
				Config: testAccKmsKey_update1(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "key_alias", updateName),
					resource.TestCheckResourceAttr(resourceName, "key_description", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "newbar"),
					resource.TestCheckResourceAttr(resourceName, "tags.acc", "test"),
					resource.TestCheckResourceAttrSet(resourceName, "is_enabled"),
					resource.TestCheckResourceAttrSet(resourceName, "rotation_enabled"),
				),
			},
			{
				Config: testAccKmsKey_update2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "key_alias", updateName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "rotation_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "rotation_interval"),
				),
			},
			{
				Config: testAccKmsKey_update3(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "key_alias", updateName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "rotation_interval", "30"),
				),
			},
			{
				Config: testAccKmsKey_update4(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "key_alias", updateName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				Config: testAccKmsKey_update5(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "key_alias", updateName),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "key_state"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"pending_days",
				},
			},
		},
	})
}

func TestAccKmsKey_ExternalKey(t *testing.T) {
	var (
		keyAlias     = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_kms_key.test"
		key          keys.Key
	)
	rc := acceptance.InitResourceCheck(
		resourceName,
		&key,
		getKmsKeyResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckKms(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKey_ExternalKey(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(resourceName, "origin", "external"),
					resource.TestCheckResourceAttr(resourceName, "key_usage", "ENCRYPT_DECRYPT"),
					resource.TestCheckResourceAttr(resourceName, "key_state", "5"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"pending_days", "is_enabled",
				},
			},
		},
	})
}

func testAccKmsKey_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "test" {
  key_alias             = "%s"
  key_algorithm         = "AES_256"
  key_usage             = "ENCRYPT_DECRYPT"
  origin                = "kms"
  key_description       = "test acc"
  enterprise_project_id = "0"
  pending_days          = "7"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name)
}

func testAccKmsKey_update1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "test" {
  key_alias             = "%s"
  key_algorithm         = "AES_256"
  key_usage             = "ENCRYPT_DECRYPT"
  origin                = "kms"
  key_description       = "terraform"
  enterprise_project_id = "0"
  pending_days          = "7"

  tags = {
    foo = "newbar"
    acc = "test"
  }
}
`, name)
}

func testAccKmsKey_update2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "test" {
  key_alias             = "%s"
  key_algorithm         = "AES_256"
  key_usage             = "ENCRYPT_DECRYPT"
  origin                = "kms"
  key_description       = "terraform"
  enterprise_project_id = "0"
  rotation_enabled      = true
  pending_days          = "7"
}
`, name)
}

func testAccKmsKey_update3(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "test" {
  key_alias             = "%s"
  key_algorithm         = "AES_256"
  key_usage             = "ENCRYPT_DECRYPT"
  origin                = "kms"
  key_description       = "terraform"
  enterprise_project_id = "0"
  rotation_enabled      = true
  rotation_interval     = 30
  pending_days          = "7"
}
`, name)
}

func testAccKmsKey_update4(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "test" {
  key_alias             = "%[1]s"
  key_algorithm         = "AES_256"
  key_usage             = "ENCRYPT_DECRYPT"
  origin                = "kms"
  key_description       = "terraform"
  enterprise_project_id = "%[2]s"
  rotation_enabled      = true
  rotation_interval     = 30
  pending_days          = "7"
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccKmsKey_update5(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "test" {
  key_alias             = "%[1]s"
  key_algorithm         = "AES_256"
  key_usage             = "ENCRYPT_DECRYPT"
  origin                = "kms"
  key_description       = "terraform"
  enterprise_project_id = "%[2]s"
  rotation_enabled      = true
  rotation_interval     = 30
  is_enabled            = false
  pending_days          = "7"
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccKmsKey_ExternalKey(keyAlias string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "test" {
  key_alias    = "%s"
  pending_days = "7"
  origin       = "external"
  key_usage    = "ENCRYPT_DECRYPT"
}
`, keyAlias)
}
