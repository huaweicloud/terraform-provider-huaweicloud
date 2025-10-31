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

func getIdentityUserResourceFuncV5(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam_no_version", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}
	getUserHttpUrl := "v5/users/{user_id}"
	getUserPath := client.Endpoint + getUserHttpUrl
	getUserPath = strings.ReplaceAll(getUserPath, "{user_id}", state.Primary.ID)
	getUserOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getUserResp, err := client.Request("GET", getUserPath, &getUserOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IAM User: %s", err)
	}
	return utils.FlattenResponse(getUserResp)
}

func TestAccIdentityV5User_basic(t *testing.T) {
	var user users.User
	userName := acceptance.RandomAccResourceName()
	userNameUpdate := acceptance.RandomAccResourceName()

	resourceName := "huaweicloud_identityv5_user.user_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&user,
		getIdentityUserResourceFuncV5,
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
				Config: testAccIdentityV5User_basic(userName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", userName),
					resource.TestCheckResourceAttr(resourceName, "description", "tested by terraform"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "is_root_user"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIdentityV5User_update(userNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", userNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "description", "update by terraform"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "is_root_user"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),
				),
			},
		},
	})
}

func testAccIdentityV5User_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_user" "user_1" {
 name        = "%[1]s"
 description = "tested by terraform"
}
`, name)
}

func testAccIdentityV5User_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_user" "user_1" {
  name        = "%[1]s"
  enabled     = false
  description = "update by terraform"
}
`, name)
}
