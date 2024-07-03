package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/signs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getSignatureFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	opts := signs.ListOpts{
		InstanceId: state.Primary.Attributes["instance_id"],
		ID:         state.Primary.ID,
	}
	resp, err := signs.List(client, opts)
	if err != nil {
		return nil, err
	}
	if len(resp) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	return resp[0], nil
}

func TestAccSignature_basic(t *testing.T) {
	var (
		signature signs.Signature

		rName1 = "huaweicloud_apig_signature.with_key"
		rName2 = "huaweicloud_apig_signature.without_key"
		name   = acceptance.RandomAccResourceName()

		// lintignore:AT009
		signKey    = acctest.RandStringFromCharSet(16, acctest.CharSetAlphaNum)
		revSignKey = utils.Reverse(signKey)
		// lintignore:AT009
		signSecret    = acctest.RandStringFromCharSet(16, acctest.CharSetAlphaNum)
		revSignSecret = utils.Reverse(signSecret)

		rc1 = acceptance.InitResourceCheck(rName1, &signature, getSignatureFunc)
		rc2 = acceptance.InitResourceCheck(rName2, &signature, getSignatureFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc1.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSignature_basic_step1(name, signKey, signSecret),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName1, "name", name+"_with_key"),
					resource.TestCheckResourceAttr(rName1, "type", "basic"),
					resource.TestCheckResourceAttr(rName1, "key", signKey),
					resource.TestCheckResourceAttr(rName1, "secret", signSecret),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "name", name+"_without_key"),
					resource.TestCheckResourceAttr(rName2, "type", "basic"),
					resource.TestCheckResourceAttrSet(rName2, "key"),
					resource.TestCheckResourceAttrSet(rName2, "secret"),
				),
			},
			{
				Config: testAccSignature_basic_step2(name, revSignKey, revSignSecret),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName1, "name", name+"_with_key_update"),
					resource.TestCheckResourceAttr(rName1, "key", revSignKey),
					resource.TestCheckResourceAttr(rName1, "secret", revSignSecret),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "name", name+"_without_key_update"),
					resource.TestCheckResourceAttr(rName2, "key", revSignKey),
					resource.TestCheckResourceAttr(rName2, "secret", revSignSecret),
				),
			},
			{
				ResourceName:      rName1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSignatureImportStateFunc(rName1),
			},
			{
				ResourceName:      rName2,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSignatureImportStateFunc(rName2),
			},
		},
	})
}

func testAccSignatureImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", rName)
		}
		if rs.Primary.Attributes["instance_id"] == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>', but got '%s/%s'",
				rs.Primary.Attributes["instance_id"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID), nil
	}
}

func testAccSignature_basic_step1(name, signKey, signSecret string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_signature" "with_key" {
  instance_id = "%[1]s"
  name        = "%[2]s_with_key"
  type        = "basic"
  key         = "%[3]s"
  secret      = "%[4]s"
}

resource "huaweicloud_apig_signature" "without_key" {
  instance_id = "%[1]s"
  name        = "%[2]s_without_key"
  type        = "basic"
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name, signKey, signSecret)
}

func testAccSignature_basic_step2(name, signKey, signSecret string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_signature" "with_key" {
  instance_id = "%[1]s"
  name        = "%[2]s_with_key_update"
  type        = "basic"
  key         = "%[3]s"
  secret      = "%[4]s"
}

resource "huaweicloud_apig_signature" "without_key" {
  instance_id = "%[1]s"
  name        = "%[2]s_without_key_update"
  type        = "basic"
  key         = "%[3]s"
  secret      = "%[4]s"
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name, signKey, signSecret)
}

func TestAccSignature_hmac(t *testing.T) {
	var (
		signature signs.Signature

		rName1 = "huaweicloud_apig_signature.with_key"
		rName2 = "huaweicloud_apig_signature.without_key"
		name   = acceptance.RandomAccResourceName()

		// lintignore:AT009
		signKey    = acctest.RandStringFromCharSet(16, acctest.CharSetAlphaNum)
		revSignKey = utils.Reverse(signKey)
		// lintignore:AT009
		signSecret    = acctest.RandStringFromCharSet(16, acctest.CharSetAlphaNum)
		revSignSecret = utils.Reverse(signSecret)

		rc1 = acceptance.InitResourceCheck(rName1, &signature, getSignatureFunc)
		rc2 = acceptance.InitResourceCheck(rName2, &signature, getSignatureFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc1.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSignature_hmac_step1(name, signKey, signSecret),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName1, "name", name+"_with_key"),
					resource.TestCheckResourceAttr(rName1, "type", "hmac"),
					resource.TestCheckResourceAttr(rName1, "key", signKey),
					resource.TestCheckResourceAttr(rName1, "secret", signSecret),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "name", name+"_without_key"),
					resource.TestCheckResourceAttr(rName2, "type", "hmac"),
					resource.TestCheckResourceAttrSet(rName2, "key"),
					resource.TestCheckResourceAttrSet(rName2, "secret"),
				),
			},
			{
				Config: testAccSignature_hmac_step2(name, revSignKey, revSignSecret),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName1, "name", name+"_with_key_update"),
					resource.TestCheckResourceAttr(rName1, "key", revSignKey),
					resource.TestCheckResourceAttr(rName1, "secret", revSignSecret),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "name", name+"_without_key_update"),
					resource.TestCheckResourceAttr(rName2, "key", revSignKey),
					resource.TestCheckResourceAttr(rName2, "secret", revSignSecret),
				),
			},
			{
				ResourceName:      rName1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSignatureImportStateFunc(rName1),
			},
			{
				ResourceName:      rName2,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSignatureImportStateFunc(rName2),
			},
		},
	})
}

func testAccSignature_hmac_step1(name, signKey, signSecret string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_signature" "with_key" {
  instance_id = "%[1]s"
  name        = "%[2]s_with_key"
  type        = "hmac"
  key         = "%[3]s"
  secret      = "%[4]s"
}

resource "huaweicloud_apig_signature" "without_key" {
  instance_id = "%[1]s"
  name        = "%[2]s_without_key"
  type        = "hmac"
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name, signKey, signSecret)
}

func testAccSignature_hmac_step2(name, signKey, signSecret string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_signature" "with_key" {
  instance_id = "%[1]s"
  name        = "%[2]s_with_key_update"
  type        = "hmac"
  key         = "%[3]s"
  secret      = "%[4]s"
}

resource "huaweicloud_apig_signature" "without_key" {
  instance_id = "%[1]s"
  name        = "%[2]s_without_key_update"
  type        = "hmac"
  key         = "%[3]s"
  secret      = "%[4]s"
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name, signKey, signSecret)
}

func TestAccSignature_aes(t *testing.T) {
	var (
		signature signs.Signature

		rName1 = "huaweicloud_apig_signature.with_key"
		rName2 = "huaweicloud_apig_signature.without_key"
		name   = acceptance.RandomAccResourceName()

		// lintignore:AT009
		signKey    = acctest.RandStringFromCharSet(16, acctest.CharSetAlphaNum)
		revSignKey = utils.Reverse(signKey)
		// lintignore:AT009
		signSecret    = acctest.RandStringFromCharSet(16, acctest.CharSetAlphaNum)
		revSignSecret = utils.Reverse(signSecret)

		rc1 = acceptance.InitResourceCheck(rName1, &signature, getSignatureFunc)
		rc2 = acceptance.InitResourceCheck(rName2, &signature, getSignatureFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc1.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSignature_aes_step1(name, signKey, signSecret),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName1, "name", name+"_with_key"),
					resource.TestCheckResourceAttr(rName1, "type", "aes"),
					resource.TestCheckResourceAttr(rName1, "algorithm", "aes-128-cfb"),
					resource.TestCheckResourceAttr(rName1, "key", signKey),
					resource.TestCheckResourceAttr(rName1, "secret", signSecret),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "name", name+"_without_key"),
					resource.TestCheckResourceAttr(rName2, "type", "aes"),
					resource.TestCheckResourceAttrSet(rName2, "key"),
					resource.TestCheckResourceAttrSet(rName2, "secret"),
				),
			},
			{
				Config: testAccSignature_aes_step2(name, revSignKey, revSignSecret),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName1, "name", name+"_with_key_update"),
					resource.TestCheckResourceAttr(rName1, "key", revSignKey),
					resource.TestCheckResourceAttr(rName1, "secret", revSignSecret),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "name", name+"_without_key_update"),
					resource.TestCheckResourceAttr(rName2, "key", revSignKey+signKey),
					resource.TestCheckResourceAttr(rName2, "secret", revSignSecret),
				),
			},
			{
				ResourceName:      rName1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSignatureImportStateFunc(rName1),
			},
			{
				ResourceName:      rName2,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSignatureImportStateFunc(rName2),
			},
		},
	})
}

func testAccSignature_aes_step1(name, signKey, signSecret string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_signature" "with_key" {
  instance_id = "%[1]s"
  name        = "%[2]s_with_key"
  type        = "aes"
  algorithm   = "aes-128-cfb"
  key         = "%[3]s"
  secret      = "%[4]s"
}

resource "huaweicloud_apig_signature" "without_key" {
  instance_id = "%[1]s"
  name        = "%[2]s_without_key"
  type        = "aes"
  algorithm   = "aes-256-cfb"
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name, signKey, signSecret)
}

// The length of the signature key and signature secret are both 16.
func testAccSignature_aes_step2(name, signKey, signSecret string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_signature" "with_key" {
  instance_id = "%[1]s"
  name        = "%[2]s_with_key_update"
  type        = "aes"
  algorithm   = "aes-128-cfb"
  key         = "%[3]s"
  secret      = "%[4]s"
}

resource "huaweicloud_apig_signature" "without_key" {
  instance_id = "%[1]s"
  name        = "%[2]s_without_key_update"
  type        = "aes"
  algorithm   = "aes-256-cfb"
  key         = format("%%s%%s", "%[3]s", strrev("%[3]s")) # the length of the 256 signature key is 32.
  secret      = "%[4]s"
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name, signKey, signSecret)
}
