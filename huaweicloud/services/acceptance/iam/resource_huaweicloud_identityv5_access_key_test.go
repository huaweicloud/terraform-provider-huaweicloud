package iam

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/v3.0/credentials"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getIdentityAccessKeyResourceFuncV5(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam_no_version", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}
	getAccessKeyHttpUrl := "v5/users/{user_id}/access-keys"
	getAccessKeyPath := client.Endpoint + getAccessKeyHttpUrl
	getAccessKeyPath = strings.ReplaceAll(getAccessKeyPath, "{user_id}", state.Primary.Attributes["user_id"])
	getAccessKeyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getAccessKeyResp, err := client.Request("GET", getAccessKeyPath, &getAccessKeyOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IAM access key: %s", err)
	}
	getAccessKeyRespBody, err := utils.FlattenResponse(getAccessKeyResp)
	if err != nil {
		return nil, err
	}

	accessKey := utils.PathSearch(fmt.Sprintf("access_keys[?access_key_id=='%s']|[0]", state.Primary.ID), getAccessKeyRespBody, nil)
	if accessKey == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return accessKey, nil
}

func TestAccIdentityV5AccessKey_basic(t *testing.T) {
	var cred credentials.Credential
	var userName = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_identityv5_access_key.key_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&cred,
		getIdentityAccessKeyResourceFuncV5,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityAccessKeyV5_basic(userName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "inactive"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "access_key_id"),
					resource.TestCheckResourceAttrSet(resourceName, "secret_access_key"),
				),
			},
			{
				Config: testAccIdentityAccessKeyV5_update(userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testIdentityV5AccessKeyImportState(resourceName),
				ImportStateVerifyIgnore: []string{"secret_access_key"},
			},
		},
	})
}

func testAccIdentityAccessKeyV5_basic(userName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_user" "user_1" {
 name        = "%[1]s"
 description = "tested by terraform"
}

resource "huaweicloud_identityv5_access_key" "key_1" {
  user_id = huaweicloud_identityv5_user.user_1.id
  status  = "inactive"
}
`, userName)
}

func testAccIdentityAccessKeyV5_update(userName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_user" "user_1" {
 name        = "%[1]s"
 description = "tested by terraform"
}

resource "huaweicloud_identityv5_access_key" "key_1" {
  user_id = huaweicloud_identityv5_user.user_1.id
  status  = "active"
}
`, userName)
}

func testIdentityV5AccessKeyImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		userId := rs.Primary.Attributes["user_id"]
		if userId == "" {
			return "", fmt.Errorf("attribute (user_id) of Resource (%s) not found: %s", name, rs)
		}
		return userId + "/" + rs.Primary.ID, nil
	}
}
