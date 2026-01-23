package iam

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

func getV5LoginProfileFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam", acceptance.HW_REGION_NAME)
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

// Please ensure that the user executing the acceptance test has 'admin' permission.
func TestAccV5LoginProfile_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_identityv5_login_profile.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getV5LoginProfileFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source: "hashicorp/random",
				// The version of the random provider must be greater than 3.3.0 to support the 'numeric' parameter.
				VersionConstraint: "3.3.0",
			},
		},
		CheckDestroy: rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV5LoginProfile_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "user_id", "huaweicloud_identityv5_user.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "password", "random_string.test.0", "result"),
					resource.TestCheckResourceAttr(rName, "password_reset_required", "true"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				Config: testAccV5LoginProfile_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrPair(rName, "password", "random_string.test.1", "result"),
					resource.TestCheckResourceAttr(rName, "password_reset_required", "false"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccV5LoginProfile_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_user" "test" {
 name = "%[1]s"
}

resource "random_string" "test" {
  count = 2

  length           = 10
  min_numeric      = 1
  min_special      = 1
  min_lower        = 1
  override_special = "@!"
}
`, name)
}

func testAccV5LoginProfile_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identityv5_login_profile" "test" {
  user_id                 = huaweicloud_identityv5_user.test.id
  password                = random_string.test[0].result
  password_reset_required = true
}
`, testAccV5LoginProfile_base(name))
}

func testAccV5LoginProfile_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identityv5_login_profile" "test" {
  user_id                 = huaweicloud_identityv5_user.test.id
  password                = random_string.test[1].result
  password_reset_required = false
}
`, testAccV5LoginProfile_base(name))
}
