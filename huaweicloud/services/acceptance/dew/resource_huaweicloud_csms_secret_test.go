package dew

import (
	"fmt"
	"testing"

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
		return nil, fmt.Errorf("error creating CSMS(KMS) client: %s", err)
	}
	name := state.Primary.Attributes["name"]
	return secrets.Get(client, name)
}

func TestAccDewCsmsSecret_basic(t *testing.T) {
	var secret secrets.Secret
	name := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_csms_secret.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secret,
		geCsmsSecretFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDewCsmsSecret_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "secret_type", "COMMON"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", "csms secret test"),
					resource.TestCheckResourceAttr(
						resourceName,
						"secret_text",
						utils.HashAndHexEncode("this is a password"),
					),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccDewCsmsSecret_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(
						resourceName,
						"secret_text",
						utils.HashAndHexEncode(`{"password":"123456","username":"admin"}`),
					),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "new_bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.hello", "world"),
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

func testAccDewCsmsSecret_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_csms_secret" "test" {
  name                  = "%s"
  description           = "csms secret test"
  secret_text           = "this is a password"
  secret_type           = "COMMON"
  enterprise_project_id = "0"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name)
}

func testAccDewCsmsSecret_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_csms_secret" "test" {
  name                  = "%s"
  description           = ""
  secret_type           = "COMMON"
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
`, name)
}

func TestAccDewCsmsSecret_SecretBinaryAndExpireTime(t *testing.T) {
	var secret secrets.Secret
	name := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_csms_secret.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secret,
		geCsmsSecretFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDewCsmsSecret_SecretBinary(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(
						resourceName,
						"secret_binary",
						utils.HashAndHexEncode(`1010`),
					),
				),
			},
			{
				Config: testAccDewCsmsSecret_updateSecretBinary(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(
						resourceName,
						"secret_binary",
						utils.HashAndHexEncode(`0101`),
					),
					resource.TestCheckResourceAttr(resourceName, "expire_time", "2999886177000"),
				),
			},
			{
				Config: testAccDewCsmsSecret_updateExpireTime(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(
						resourceName,
						"secret_binary",
						utils.HashAndHexEncode(`0101`),
					),
					resource.TestCheckResourceAttr(resourceName, "expire_time", "3999886177000"),
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

func testAccDewCsmsSecret_SecretBinary(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_csms_secret" "test" {
  name          = "%s"
  secret_binary = "1010"
}
`, name)
}

func testAccDewCsmsSecret_updateSecretBinary(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_csms_secret" "test" {
  name          = "%s"
  secret_binary = "0101"
  expire_time   = 2999886177000
}
`, name)
}

func testAccDewCsmsSecret_updateExpireTime(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_csms_secret" "test" {
  name          = "%s"
  secret_binary = "0101"
  expire_time   = 3999886177000
}
`, name)
}

func TestAccDewCsmsSecret_associatedEvent(t *testing.T) {
	var secret secrets.Secret
	name := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_csms_secret.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secret,
		geCsmsSecretFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDewCsmsSecret_associatedEvent_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "secret_type", "COMMON"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttrPair(resourceName, "event_subscriptions.0", "huaweicloud_csms_event.test", "name"),
					resource.TestCheckResourceAttr(
						resourceName,
						"secret_text",
						utils.HashAndHexEncode("terraform"),
					),
				),
			},
			{
				Config: testAccDewCsmsSecret_associatedEvent_update_1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "event_subscriptions.0", "huaweicloud_csms_event.retest", "name"),
				),
			},
			{
				Config: testAccDewCsmsSecret_associatedEvent_update_2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "event_subscriptions.#", "0"),
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

func testAccDewCsmsSecret_associatedEvent_base(name string) string {
	return fmt.Sprintf(`
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

func testAccDewCsmsSecret_associatedEvent_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_csms_secret" "test" {
  name                  = "%[2]s"
  description           = "test acc"
  secret_text           = "terraform"
  secret_type           = "COMMON"
  enterprise_project_id = "0"
  event_subscriptions   = [huaweicloud_csms_event.test.name]
}
`, testAccDewCsmsSecret_associatedEvent_base(name), name)
}

func testAccDewCsmsSecret_associatedEvent_update_1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_csms_secret" "test" {
  name                  = "%[2]s"
  description           = "test acc"
  secret_text           = "terraform"
  secret_type           = "COMMON"
  enterprise_project_id = "0"
  event_subscriptions   = [huaweicloud_csms_event.retest.name]
}
`, testAccDewCsmsSecret_associatedEvent_base(name), name)
}

func testAccDewCsmsSecret_associatedEvent_update_2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_csms_secret" "test" {
  name                  = "%[2]s"
  description           = "test acc"
  secret_text           = "terraform"
  secret_type           = "COMMON"
  enterprise_project_id = "0"
}
`, testAccDewCsmsSecret_associatedEvent_base(name), name)
}
