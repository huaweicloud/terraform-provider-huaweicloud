package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dataarts"
)

func getStudioWorkspaceUserFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("dataarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio client: %s", err)
	}

	return dataarts.GetStudioWorkspaceUserById(client, state.Primary.Attributes["workspace_id"], state.Primary.ID)
}

func TestAccStudioWorkspaceUser_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_dataarts_studio_workspace_user.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getStudioWorkspaceUserFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckUserName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccStudioWorkspaceUser_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttrPair(resourceName, "user_id", "data.huaweicloud_identity_users.test", "users.0.id"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "roles.0.id",
						"data.huaweicloud_dataarts_studio_workspace_user_roles.test", "roles.0.id"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccStudioWorkspaceUser_basic_step2(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "roles.0.id",
						"data.huaweicloud_dataarts_studio_workspace_user_roles.test", "roles.1.id"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(resourceName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStudioWorkspaceUserImportStateIDFunc(resourceName),
			},
		},
	})
}

func testAccStudioWorkspaceUserImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		workspaceId := rs.Primary.Attributes["workspace_id"]
		userId := rs.Primary.Attributes["user_id"]
		if workspaceId == "" || userId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<user_id>', but got '%s/%s'",
				workspaceId, userId)
		}
		return fmt.Sprintf("%s/%s", workspaceId, userId), nil
	}
}

func testAccStudioWorkspaceUser_basic_step1() string {
	return fmt.Sprintf(`
data "huaweicloud_dataarts_studio_workspace_user_roles" "test" {
  workspace_id = "%[1]s"
}

data "huaweicloud_identity_users" "test" {
  name = "%[2]s"
}

resource "huaweicloud_dataarts_studio_workspace_user" "test" {
  workspace_id = "%[1]s"
  user_id      = try(data.huaweicloud_identity_users.test.users[0].id, "NOT_FOUND")

  dynamic "roles" {
    for_each = slice(data.huaweicloud_dataarts_studio_workspace_user_roles.test.roles, 0, 1)

    content {
      id = roles.value.id
    }
  }
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_USER_NAME)
}

func testAccStudioWorkspaceUser_basic_step2() string {
	return fmt.Sprintf(`
data "huaweicloud_dataarts_studio_workspace_user_roles" "test" {
  workspace_id = "%[1]s"
}

data "huaweicloud_identity_users" "test" {
  name = "%[2]s"
}

resource "huaweicloud_dataarts_studio_workspace_user" "test" {
  workspace_id = "%[1]s"
  user_id      = try(data.huaweicloud_identity_users.test.users[0].id, "NOT_FOUND")

  dynamic "roles" {
    for_each = slice(data.huaweicloud_dataarts_studio_workspace_user_roles.test.roles, 1, 2)

    content {
      id = roles.value.id
    }
  }
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_USER_NAME)
}
