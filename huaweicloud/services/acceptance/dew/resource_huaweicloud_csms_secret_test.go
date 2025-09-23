package dew

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/csms/v1/secrets"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func geCsmsSecretFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.KmsV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating KMS client: %s", err)
	}
	name := state.Primary.Attributes["name"]
	return secrets.Get(client, name)
}

func TestAccDewCsmsSecret_basic(t *testing.T) {
	var (
		secret           secrets.Secret
		name             = acceptance.RandomAccResourceName()
		resourceName     = "huaweicloud_csms_secret.test"
		secretText       = utils.HashAndHexEncode("this is a password")
		secretTextUpdate = utils.HashAndHexEncode(`{"password":"123456","username":"admin"}`)
		expireTime       = time.Now().Add(48*time.Hour).Unix() * 1000
		expireUpdateTime = time.Now().Add(72*time.Hour).Unix() * 1000
		binary           = utils.HashAndHexEncode(`1010`)
		binaryUpdate     = utils.HashAndHexEncode(`0101`)
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secret,
		geCsmsSecretFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDewCsmsSecret_basic(name, expireTime),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "secret_text", secretText),
					resource.TestCheckResourceAttr(resourceName, "expire_time", fmt.Sprintf("%d", expireTime)),
					resource.TestCheckResourceAttr(resourceName, "description", "csms secret test"),
					resource.TestCheckResourceAttr(resourceName, "secret_type", "COMMON"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "secret_binary", ""),
					resource.TestCheckResourceAttrPair(resourceName, "kms_key_id",
						"huaweicloud_kms_key.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "event_subscriptions.0",
						"huaweicloud_csms_event.test", "name"),
					resource.TestCheckResourceAttrSet(resourceName, "secret_id"),
					resource.TestCheckResourceAttrSet(resourceName, "latest_version"),
					resource.TestCheckResourceAttrSet(resourceName, "version_stages.#"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
			{
				Config: testAccDewCsmsSecret_update1(name, expireUpdateTime),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "secret_text", secretTextUpdate),
					resource.TestCheckResourceAttr(resourceName, "expire_time", fmt.Sprintf("%d", expireUpdateTime)),
					resource.TestCheckResourceAttr(resourceName, "description", "csms secret test update"),
					resource.TestCheckResourceAttr(resourceName, "secret_type", "COMMON"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "new_bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.hello", "world"),
					resource.TestCheckResourceAttr(resourceName, "secret_binary", ""),
					resource.TestCheckResourceAttrPair(resourceName, "kms_key_id",
						"huaweicloud_kms_key.test_second", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "event_subscriptions.0",
						"huaweicloud_csms_event.retest", "name"),
					resource.TestCheckResourceAttrSet(resourceName, "secret_id"),
					resource.TestCheckResourceAttrSet(resourceName, "latest_version"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
			{
				Config: testAccDewCsmsSecret_update2(name, expireUpdateTime),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "event_subscriptions.#", "0"),
				),
			},
			{
				Config: testAccDewCsmsSecret_update3(name, expireUpdateTime),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "secret_binary", binary),
				),
			},
			{
				Config: testAccDewCsmsSecret_update4(name, expireUpdateTime),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "secret_binary", binaryUpdate),
				),
			},
			{
				Config: testAccDewCsmsSecret_update5(name, expireUpdateTime),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDewCsmsSecret_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "test" {
  key_alias    = "%[1]s"
  pending_days = "7"
}

resource "huaweicloud_kms_key" "test_second" {
  key_alias    = "%[1]s_second"
  pending_days = "7"
}

resource "huaweicloud_smn_topic" "test" {
  name = "%[1]s"
}

resource "huaweicloud_csms_event" "test" {
  name                     = "%[1]s-01"
  event_types              = ["SECRET_VERSION_CREATED", "SECRET_ROTATED"]
  status                   = "ENABLED"
  notification_target_type = "SMN"
  notification_target_id   = huaweicloud_smn_topic.test.id
  notification_target_name = huaweicloud_smn_topic.test.name
}

resource "huaweicloud_csms_event" "retest" {
  name                     = "%[1]s-02"
  event_types              = ["SECRET_VERSION_CREATED"]
  status                   = "ENABLED"
  notification_target_type = "SMN"
  notification_target_id   = huaweicloud_smn_topic.test.id
  notification_target_name = huaweicloud_smn_topic.test.name
}
`, name)
}

func testAccDewCsmsSecret_basic(name string, expireTime int64) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_csms_secret" "test" {
  name                  = "%[2]s"
  secret_text           = "this is a password"
  description           = "csms secret test"
  expire_time           = %[3]d
  kms_key_id            = huaweicloud_kms_key.test.id
  secret_type           = "COMMON"
  event_subscriptions   = [huaweicloud_csms_event.test.name]
  enterprise_project_id = "0"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccDewCsmsSecret_base(name), name, expireTime)
}

func testAccDewCsmsSecret_update1(name string, expireTime int64) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_csms_secret" "test" {
  name                  = "%[2]s"
  description           = "csms secret test update"
  expire_time           = %[3]d
  kms_key_id            = huaweicloud_kms_key.test_second.id
  secret_type           = "COMMON"
  event_subscriptions   = [huaweicloud_csms_event.retest.name]
  enterprise_project_id = "0"

  secret_text = jsonencode({
    username = "admin"
    password = "123456"
  })

  tags = {
    foo   = "new_bar"
    hello = "world"
  }
}
`, testAccDewCsmsSecret_base(name), name, expireTime)
}

func testAccDewCsmsSecret_update2(name string, expireTime int64) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_csms_secret" "test" {
  name                  = "%[2]s"
  description           = ""
  expire_time           = %[3]d
  kms_key_id            = huaweicloud_kms_key.test_second.id
  secret_type           = "COMMON"
  enterprise_project_id = "0"

  secret_text = jsonencode({
    username = "admin"
    password = "123456"
  })
}
`, testAccDewCsmsSecret_base(name), name, expireTime)
}

func testAccDewCsmsSecret_update3(name string, expireTime int64) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_csms_secret" "test" {
  name                  = "%[2]s"
  secret_binary         = "1010"
  description           = ""
  expire_time           = %[3]d
  kms_key_id            = huaweicloud_kms_key.test_second.id
  secret_type           = "COMMON"
  enterprise_project_id = "0"
}
`, testAccDewCsmsSecret_base(name), name, expireTime)
}

func testAccDewCsmsSecret_update4(name string, expireTime int64) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_csms_secret" "test" {
  name                  = "%[2]s"
  secret_binary         = "0101"
  description           = ""
  expire_time           = %[3]d
  kms_key_id            = huaweicloud_kms_key.test_second.id
  secret_type           = "COMMON"
  enterprise_project_id = "0"
}
`, testAccDewCsmsSecret_base(name), name, expireTime)
}

func testAccDewCsmsSecret_update5(name string, expireTime int64) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_csms_secret" "test" {
  name                  = "%[2]s"
  secret_binary         = "0101"
  description           = ""
  expire_time           = %[3]d
  kms_key_id            = huaweicloud_kms_key.test_second.id
  secret_type           = "COMMON"
  enterprise_project_id = "%[4]s"
}
`, testAccDewCsmsSecret_base(name), name, expireTime, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
