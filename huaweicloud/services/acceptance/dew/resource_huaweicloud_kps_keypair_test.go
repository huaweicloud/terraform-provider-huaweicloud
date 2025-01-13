package dew

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getKpsKeypairResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v3/{project_id}/keypairs/{keypair_name}"
		product = "kms"
	)

	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating KMS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{keypair_name}", state.Primary.ID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func TestAccKeypair_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_kps_keypair.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getKpsKeypairResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckUserId(t)
			// Please provide a file path to write the private key file, for example: /temp/XXX/temp.crt.
			acceptance.TestAccPreCheckKPSKeyFilePath(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKeypair_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "encryption_type", "kms"),
					resource.TestCheckResourceAttr(rName, "key_file", acceptance.HW_KPS_KEY_FILE_PATH),
					resource.TestCheckResourceAttr(rName, "user_id", acceptance.HW_USER_ID),
					resource.TestCheckResourceAttr(rName, "scope", "user"),

					resource.TestCheckResourceAttrPair(rName, "kms_key_id", "huaweicloud_kms_key.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "kms_key_name", "huaweicloud_kms_key.test", "key_alias"),

					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "fingerprint"),
					resource.TestCheckResourceAttrSet(rName, "is_managed"),
					resource.TestCheckResourceAttrSet(rName, "public_key"),
				),
			},
			{
				Config: testAccKeypair_basic_update1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
					resource.TestCheckResourceAttr(rName, "encryption_type", "default"),
					resource.TestCheckResourceAttr(rName, "key_file", acceptance.HW_KPS_KEY_FILE_PATH),
					resource.TestCheckResourceAttr(rName, "user_id", acceptance.HW_USER_ID),
					resource.TestCheckResourceAttr(rName, "scope", "user"),
					resource.TestCheckResourceAttr(rName, "kms_key_id", ""),
					resource.TestCheckResourceAttr(rName, "kms_key_name", ""),

					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "fingerprint"),
					resource.TestCheckResourceAttrSet(rName, "is_managed"),
					resource.TestCheckResourceAttrSet(rName, "public_key"),
					resource.TestCheckResourceAttrSet(rName, "private_key"),
				),
			},
			{
				Config: testAccKeypair_basic_update2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
					resource.TestCheckResourceAttr(rName, "encryption_type", "default"),
					resource.TestCheckResourceAttr(rName, "key_file", acceptance.HW_KPS_KEY_FILE_PATH),
					resource.TestCheckResourceAttr(rName, "user_id", acceptance.HW_USER_ID),
					resource.TestCheckResourceAttr(rName, "scope", "user"),
					resource.TestCheckResourceAttr(rName, "kms_key_id", ""),
					resource.TestCheckResourceAttr(rName, "kms_key_name", ""),
					resource.TestCheckResourceAttr(rName, "private_key", ""),

					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "fingerprint"),
					resource.TestCheckResourceAttrSet(rName, "is_managed"),
					resource.TestCheckResourceAttrSet(rName, "public_key"),
				),
			},
			{
				Config: testAccKeypair_basic_update3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
					resource.TestCheckResourceAttr(rName, "encryption_type", "kms"),
					resource.TestCheckResourceAttr(rName, "key_file", acceptance.HW_KPS_KEY_FILE_PATH),
					resource.TestCheckResourceAttr(rName, "user_id", acceptance.HW_USER_ID),
					resource.TestCheckResourceAttr(rName, "scope", "user"),

					resource.TestCheckResourceAttrPair(rName, "kms_key_id", "huaweicloud_kms_key.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "kms_key_name", "huaweicloud_kms_key.test", "key_alias"),

					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "fingerprint"),
					resource.TestCheckResourceAttrSet(rName, "is_managed"),
					resource.TestCheckResourceAttrSet(rName, "public_key"),
					resource.TestCheckResourceAttrSet(rName, "private_key"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"private_key",
					"encryption_type",
					"kms_key_id",
					"kms_key_name",
					"key_file",
				},
			},
		},
	})
}

func testAccKeypair_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "test" {
  key_alias    = "%[1]s"
  pending_days = "7"
}
`, name)
}

func testAccKeypair_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_kps_keypair" "test" {
  name            = "%[2]s"
  scope           = "user"
  user_id         = "%[3]s"
  encryption_type = "kms"
  kms_key_id      = huaweicloud_kms_key.test.id
  kms_key_name    = huaweicloud_kms_key.test.key_alias
  description     = "test description"
  key_file        = "%[4]s"
}
`, testAccKeypair_base(name), name, acceptance.HW_USER_ID, acceptance.HW_KPS_KEY_FILE_PATH)
}

func testAccKeypair_basic_update1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_kps_keypair" "test" {
  name            = "%[2]s"
  scope           = "user"
  user_id         = "%[3]s"
  encryption_type = "default"
  private_key     = file("%[4]s")
  description     = "test description update"
  key_file        = "%[4]s"
}
`, testAccKeypair_base(name), name, acceptance.HW_USER_ID, acceptance.HW_KPS_KEY_FILE_PATH)
}

func testAccKeypair_basic_update2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_kps_keypair" "test" {
  name            = "%[2]s"
  scope           = "user"
  user_id         = "%[3]s"
  encryption_type = "default"
  description     = "test description update"
  key_file        = "%[4]s"
}
`, testAccKeypair_base(name), name, acceptance.HW_USER_ID, acceptance.HW_KPS_KEY_FILE_PATH)
}

func testAccKeypair_basic_update3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_kps_keypair" "test" {
  name            = "%[2]s"
  scope           = "user"
  user_id         = "%[3]s"
  encryption_type = "kms"
  private_key     = file("%[4]s")
  kms_key_id      = huaweicloud_kms_key.test.id
  kms_key_name    = huaweicloud_kms_key.test.key_alias
  description     = "test description update"
  key_file        = "%[4]s"
}
`, testAccKeypair_base(name), name, acceptance.HW_USER_ID, acceptance.HW_KPS_KEY_FILE_PATH)
}

func TestAccKeypair_existKeypair(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_kps_keypair.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getKpsKeypairResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please provide a file path to write the private key file, for example: /temp/XXX/temp.crt.
			acceptance.TestAccPreCheckKPSKeyFilePath(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKeypair_import(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "scope", "account"),

					resource.TestCheckResourceAttrPair(rName, "public_key", "huaweicloud_kps_keypair.test-base", "public_key"),

					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "fingerprint"),
					resource.TestCheckResourceAttrSet(rName, "is_managed"),
					resource.TestCheckResourceAttrSet(rName, "user_id"),
				),
			},
			{
				Config: testAccKeypair_import_update1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
					resource.TestCheckResourceAttr(rName, "encryption_type", "kms"),
					resource.TestCheckResourceAttr(rName, "scope", "account"),

					resource.TestCheckResourceAttrPair(rName, "kms_key_id", "huaweicloud_kms_key.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "kms_key_name", "huaweicloud_kms_key.test", "key_alias"),
					resource.TestCheckResourceAttrPair(rName, "public_key", "huaweicloud_kps_keypair.test-base", "public_key"),

					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "fingerprint"),
					resource.TestCheckResourceAttrSet(rName, "is_managed"),
					resource.TestCheckResourceAttrSet(rName, "private_key"),
					resource.TestCheckResourceAttrSet(rName, "user_id"),
				),
			},
			{
				Config: testAccKeypair_import_update2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
					resource.TestCheckResourceAttr(rName, "encryption_type", "default"),
					resource.TestCheckResourceAttr(rName, "scope", "account"),
					resource.TestCheckResourceAttr(rName, "kms_key_id", ""),
					resource.TestCheckResourceAttr(rName, "kms_key_name", ""),

					resource.TestCheckResourceAttrPair(rName, "public_key", "huaweicloud_kps_keypair.test-base", "public_key"),

					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "fingerprint"),
					resource.TestCheckResourceAttrSet(rName, "is_managed"),
					resource.TestCheckResourceAttrSet(rName, "private_key"),
					resource.TestCheckResourceAttrSet(rName, "user_id"),
				),
			},
			{
				Config: testAccKeypair_import_update3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
					resource.TestCheckResourceAttr(rName, "encryption_type", "default"),
					resource.TestCheckResourceAttr(rName, "scope", "account"),
					resource.TestCheckResourceAttr(rName, "kms_key_id", ""),
					resource.TestCheckResourceAttr(rName, "kms_key_name", ""),
					resource.TestCheckResourceAttr(rName, "private_key", ""),

					resource.TestCheckResourceAttrPair(rName, "public_key", "huaweicloud_kps_keypair.test-base", "public_key"),

					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "fingerprint"),
					resource.TestCheckResourceAttrSet(rName, "is_managed"),
					resource.TestCheckResourceAttrSet(rName, "user_id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"encryption_type",
					"kms_key_id",
					"kms_key_name",
					"private_key",
				},
			},
		},
	})
}

func testAccKeypair_import_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "test" {
  key_alias    = "%[1]s"
  pending_days = "7"
}

resource "huaweicloud_kps_keypair" "test-base" {
  name            = "%[1]s-base"
  scope           = "user"
  encryption_type = "default"
  description     = "test description"
  key_file        = "%[2]s"
}
`, name, acceptance.HW_KPS_KEY_FILE_PATH)
}

func testAccKeypair_import(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_kps_keypair" "test" {
  name        = "%[2]s"
  scope       = "account"
  description = "test description"
  public_key  = huaweicloud_kps_keypair.test-base.public_key
}
`, testAccKeypair_import_base(name), name)
}

func testAccKeypair_import_update1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_kps_keypair" "test" {
  name            = "%[2]s"
  scope           = "account"
  encryption_type = "kms"
  kms_key_id      = huaweicloud_kms_key.test.id
  kms_key_name    = huaweicloud_kms_key.test.key_alias
  description     = "test description update"
  private_key     = file("%[3]s")
  public_key      = huaweicloud_kps_keypair.test-base.public_key
}
`, testAccKeypair_import_base(name), name, acceptance.HW_KPS_KEY_FILE_PATH)
}

func testAccKeypair_import_update2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_kps_keypair" "test" {
  name            = "%[2]s"
  scope           = "account"
  encryption_type = "default"
  description     = "test description update"
  private_key     = file("%[3]s")
  public_key      = huaweicloud_kps_keypair.test-base.public_key
}
`, testAccKeypair_import_base(name), name, acceptance.HW_KPS_KEY_FILE_PATH)
}

func testAccKeypair_import_update3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_kps_keypair" "test" {
  name            = "%[2]s"
  scope           = "account"
  encryption_type = "default"
  description     = "test description update"
  public_key      = huaweicloud_kps_keypair.test-base.public_key
}
`, testAccKeypair_import_base(name), name)
}
