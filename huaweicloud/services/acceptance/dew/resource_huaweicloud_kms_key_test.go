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
		return nil, fmt.Errorf("error creating kms client: %s", err)
	}
	key, err := keys.Get(client, state.Primary.ID).ExtractKeyInfo()

	if err == nil && key.KeyState == dew.PendingDeletionState {
		return nil, golangsdk.ErrDefault404{}
	}
	return key, err
}

// Keystore_id scenario testing is currently not supported.
func TestAccKmsKey_Basic(t *testing.T) {
	var keyAlias = acceptance.RandomAccResourceName()
	var keyAliasUpdate = acceptance.RandomAccResourceName()
	var resourceName = "huaweicloud_kms_key.key_1"

	var key keys.Key

	rc := acceptance.InitResourceCheck(
		resourceName,
		&key,
		getKmsKeyResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckKms(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKey_Basic(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(resourceName, "rotation_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "region", acceptance.HW_REGION_NAME),
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
			{
				Config: testAccKmsKeyUpdate(keyAliasUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "key_alias", keyAliasUpdate),
					resource.TestCheckResourceAttr(resourceName, "key_description", "key update description"),
					resource.TestCheckResourceAttr(resourceName, "region", acceptance.HW_REGION_NAME),
				),
			},
		},
	})
}

func TestAccKmsKey_ExternalKey(t *testing.T) {
	var keyAlias = acceptance.RandomAccResourceName()
	var resourceName = "huaweicloud_kms_key.key_1"

	var key keys.Key

	rc := acceptance.InitResourceCheck(
		resourceName,
		&key,
		getKmsKeyResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckKms(t) },
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
					resource.TestCheckResourceAttr(resourceName, "region", acceptance.HW_REGION_NAME),
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

func TestAccKmsKey_Enable(t *testing.T) {
	var rName = acceptance.RandomAccResourceName()
	var resourceName = "huaweicloud_kms_key.key_1"

	var key keys.Key
	rc := acceptance.InitResourceCheck(
		resourceName,
		&key,
		getKmsKeyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckKms(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKey_enabled(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "true"),
				),
			},
			{
				Config: testAccKmsKey_disabled(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "false"),
				),
			},
			{
				Config: testAccKmsKey_enabled(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "true"),
				),
			},
		},
	})
}

func TestAccKmsKey_WithTags(t *testing.T) {
	var keyAlias = acceptance.RandomAccResourceName()
	var resourceName = "huaweicloud_kms_key.key_1"

	var key keys.Key
	rc := acceptance.InitResourceCheck(
		resourceName,
		&key,
		getKmsKeyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckKms(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKey_WithTags(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
		},
	})
}

func TestAccKmsKey_WithEpsId(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_kms_key.test"

	var key keys.Key
	rc := acceptance.InitResourceCheck(
		resourceName,
		&key,
		getKmsKeyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckKms(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKey_epsId_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "key_alias", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
				),
			},
			{
				Config: testAccKmsKey_epsId_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "key_alias", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccKmsKey_rotation(t *testing.T) {
	var keyAlias = acceptance.RandomAccResourceName()
	var resourceName = "huaweicloud_kms_key.key_1"

	var key keys.Key
	rc := acceptance.InitResourceCheck(
		resourceName,
		&key,
		getKmsKeyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckKms(t); acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKey_Basic(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(resourceName, "rotation_enabled", "false"),
				),
			},
			{
				Config: testAccKmsKey_rotation(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(resourceName, "rotation_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "rotation_interval", "365"),
				),
			},
			{
				Config: testAccKmsKey_rotation_interval(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(resourceName, "rotation_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "rotation_interval", "200"),
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

func testAccKmsKey_Basic(keyAlias string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "key_1" {
  key_alias    = "%s"
  pending_days = "7"
  region       = "%s"
}
`, keyAlias, acceptance.HW_REGION_NAME)
}

func testAccKmsKey_ExternalKey(keyAlias string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "key_1" {
  key_alias    = "%s"
  pending_days = "7"
  region       = "%s"
  origin       = "external"
  key_usage    = "ENCRYPT_DECRYPT"
}
`, keyAlias, acceptance.HW_REGION_NAME)
}

func testAccKmsKey_WithTags(keyAlias string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "key_1" {
  key_alias    = "%s"
  pending_days = "7"
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, keyAlias)
}

func testAccKmsKey_epsId_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "test" {
  key_alias             = "%[1]s"
  pending_days          = "7"
  enterprise_project_id = "0"
}
`, name)
}

func testAccKmsKey_epsId_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "test" {
  key_alias             = "%[1]s"
  pending_days          = "7"
  enterprise_project_id = "%[2]s"
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccKmsKeyUpdate(keyAliasUpdate string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "key_1" {
  key_alias       = "%s"
  key_description = "key update description"
  pending_days    = "7"
}
`, keyAliasUpdate)
}

func testAccKmsKey_enabled(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "key_1" {
  key_description = "Terraform acc test is_enabled %s"
  pending_days    = "7"
  key_alias       = "%s"
}`, rName, rName)
}

func testAccKmsKey_disabled(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "key_1" {
  key_description = "Terraform acc test is_enabled %s"
  pending_days    = "7"
  key_alias       = "%s"
  is_enabled      = false
}`, rName, rName)
}

func testAccKmsKey_rotation(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "key_1" {
  key_alias        = "%s"
  pending_days     = "7"
  rotation_enabled = true
}`, rName)
}

func testAccKmsKey_rotation_interval(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "key_1" {
  key_alias         = "%s"
  pending_days      = "7"
  rotation_enabled  = true
  rotation_interval = 200
}`, rName)
}
