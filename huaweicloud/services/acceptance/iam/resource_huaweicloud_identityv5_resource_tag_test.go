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

func getIdentityResourceTagResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.IAMNoVersionClient(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}
	getResourceTagHttpUrl := "v5/{resource_type}/{resource_id}/tags"
	getResourceTagPath := client.Endpoint + getResourceTagHttpUrl
	getResourceTagPath = strings.ReplaceAll(getResourceTagPath, "{resource_type}", state.Primary.Attributes["resource_type"])
	getResourceTagPath = strings.ReplaceAll(getResourceTagPath, "{resource_id}", state.Primary.Attributes["resource_id"])
	getResourceTagOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	getResourceTagResp, err := client.Request("GET", getResourceTagPath, &getResourceTagOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IAM resource tag: %s", err)
	}
	return utils.FlattenResponse(getResourceTagResp)
}

func TestAccV5ResourceTag_basic(t *testing.T) {
	username := acceptance.RandomAccResourceName()
	var object interface{}
	resourceName := "huaweicloud_identityv5_resource_tag.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&object,
		getIdentityResourceTagResourceFunc,
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
				Config: testAccV5ResourceTag_basic(username),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "resource_type", "user"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_id"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccV5ResourceTag_update(username),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "resource_type", "user"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_id"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
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

func testAccV5ResourceTag_basic(username string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identityv5_resource_tag" "test" {
  resource_type = "user"
  resource_id   = huaweicloud_identityv5_user.user_1.id
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccIdentityV5User_basic(username))
}

func testAccV5ResourceTag_update(username string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identityv5_resource_tag" "test" {
  resource_type = "user"
  resource_id   = huaweicloud_identityv5_user.user_1.id
  tags = {
    foo  = "bar1"
    key1 = "value"
  }
}
`, testAccIdentityV5User_basic(username))
}
