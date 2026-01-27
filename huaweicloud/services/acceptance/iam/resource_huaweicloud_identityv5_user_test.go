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

func getV5UserResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam", acceptance.HW_REGION_NAME)
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

// Please ensure that the user executing the acceptance test has 'admin' permission.
func TestAccV5User_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_identityv5_user.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getV5UserResourceFunc)
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
				Config: testAccV5User_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "tested by terraform"),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
					// Check attributes.
					resource.TestCheckResourceAttrSet(rName, "is_root_user"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "urn"),
				),
			},
			{
				Config: testAccV5User_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", "update by terraform"),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
					resource.TestCheckResourceAttr(rName, "tags.#", "0"),
				),
			},
			// Only used to check the 'tags' attribute.
			{
				Config: testAccV5User_basic_step3(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "tags.#", "1"),
					resource.TestCheckResourceAttr(rName, "tags.0.tag_key", "foo"),
					resource.TestCheckResourceAttr(rName, "tags.0.tag_value", "bar"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccV5User_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_user" "test" {
  name        = "%[1]s"
  description = "tested by terraform"
}
`, name)
}

func testAccV5User_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_user" "test" {
  name        = "%[1]s"
  enabled     = false
  description = "update by terraform"
}

resource "huaweicloud_identityv5_resource_tag" "test" {
  resource_type = "user"
  resource_id   = huaweicloud_identityv5_user.test.id

  tags = {
    foo = "bar"
  }
}
`, name)
}

// Only used to check the 'tags' attribute.
// In step two, a tag is added to the user, but 'tags' is only an attribute,
// so the value of 'tags' can only be obtained after refreshing the resource in step three.
func testAccV5User_basic_step3(name string) string {
	return testAccV5User_basic_step2(name)
}
