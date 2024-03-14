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

func getKmsKeyMaterialResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.KmsKeyV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating KMS client: %s", err)
	}

	key, err := keys.Get(client, state.Primary.ID).ExtractKeyInfo()

	if key.KeyState == dew.PendingDeletionState || key.KeyState == dew.PendingImportState {
		return nil, golangsdk.ErrDefault404{}
	}

	return key, err
}

func TestAccKmsKeyMaterial_Symmetric(t *testing.T) {
	var resourceName = "huaweicloud_kms_key_material.test"
	var key keys.Key

	rc := acceptance.InitResourceCheck(
		resourceName,
		&key,
		getKmsKeyMaterialResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			// The key status must be pending import.
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckKms(t)
			acceptance.TestAccPreCheckKmsKeyID(t)
			acceptance.TestAccPreCheckKmsImportToken(t)
			acceptance.TestAccPreCheckKmsKeyMaterial(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKeyMaterial_Symmetric(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "key_state", "2"),
					resource.TestCheckResourceAttr(resourceName, "expiration_time", "2999886177"),
					resource.TestCheckResourceAttr(resourceName, "region", acceptance.HW_REGION_NAME),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"import_token", "encrypted_key_material", "encrypted_privatekey",
				},
			},
		},
	})
}

func testAccKmsKeyMaterial_Symmetric() string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key_material" "test" {
  key_id                 = "%[1]s"
  import_token           = "%[2]s"
  encrypted_key_material = "%[3]s"
  expiration_time        = "2999886177"
}
`, acceptance.HW_KMS_KEY_ID, acceptance.HW_KMS_IMPORT_TOKEN, acceptance.HW_KMS_KEY_MATERIAL)
}

// The key material of the asymmetric key does not support deletion.
func TestAccKmsKeyMaterial_Asymmetric(t *testing.T) {
	var resourceName = "huaweicloud_kms_key_material.test"
	var key keys.Key

	rc := acceptance.InitResourceCheck(
		resourceName,
		&key,
		getKmsKeyMaterialResourceFunc,
	)
	// Avoid CheckDestroy because the asymmetric key material can not be destroyed.
	// lintignore:AT001
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			// The key status must be pending import.
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckKms(t)
			acceptance.TestAccPreCheckKmsKeyID(t)
			acceptance.TestAccPreCheckKmsImportToken(t)
			acceptance.TestAccPreCheckKmsKeyMaterial(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKeyMaterial_Asymmetric(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "key_state", "2"),
					resource.TestCheckResourceAttr(resourceName, "region", acceptance.HW_REGION_NAME),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"import_token", "encrypted_key_material", "encrypted_privatekey",
				},
			},
		},
	})
}

func testAccKmsKeyMaterial_Asymmetric() string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key_material" "test" {
  key_id                 = "%[1]s"
  import_token           = "%[2]s"
  encrypted_key_material = "%[3]s"
  encrypted_privatekey   = "%[4]s"
}
`, acceptance.HW_KMS_KEY_ID, acceptance.HW_KMS_IMPORT_TOKEN, acceptance.HW_KMS_KEY_MATERIAL, acceptance.HW_KMS_KEY_PRIVATE_KEY)
}
