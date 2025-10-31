package iam

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/v3.0/users"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getIdentityLoginProfileResourceFuncV5(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam_no_version", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}
	getLoginProfileHttpUrl := "v5/users/{user_id}/login-profile"
	getLoginProfilePath := client.Endpoint + getLoginProfileHttpUrl
	getLoginProfilePath = strings.ReplaceAll(getLoginProfilePath, "{user_id}", state.Primary.ID)
	getLoginProfileOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getLoginProfileResp, err := client.Request("GET", getLoginProfilePath, &getLoginProfileOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IAM login profile: %s", err)
	}
	return utils.FlattenResponse(getLoginProfileResp)
}

func TestAccIdentityV5LoginProfile_basic(t *testing.T) {
	var user users.User
	userName := acceptance.RandomAccResourceName()

	resourceName := "huaweicloud_identityv5_login_profile.login_profile_11"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&user,
		getIdentityLoginProfileResourceFuncV5,
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
				Config: testAccIdentityV5LoginProfile_basic(userName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "password_reset_required"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				Config: testAccIdentityV5LoginProfile_update(userName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "password_reset_required"),
				),
			},
		},
	})
}

func testAccIdentityV5LoginProfile_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_user" "user" {
 name = "%[1]s"
}

resource "huaweicloud_identityv5_login_profile" "login_profile_11" {
  user_id                 = huaweicloud_identityv5_user.user.id
  password                = "default8881"
  password_reset_required = false
}
`, name)
}

func testAccIdentityV5LoginProfile_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_user" "user" {
 name = "%[1]s"
}

resource "huaweicloud_identityv5_login_profile" "login_profile_11" {
  user_id                 = huaweicloud_identityv5_user.user.id
  password                = "default8881"
  password_reset_required = false
}
`, name)
}
