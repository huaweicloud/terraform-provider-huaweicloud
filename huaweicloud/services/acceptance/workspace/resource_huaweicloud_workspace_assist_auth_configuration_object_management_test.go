package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/workspace"
)

func GetAssistAuthConfigurationObjectManagementResourceFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("workspace", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace client: %s", err)
	}
	return workspace.ListAndReorderAssistAuthConfigurationAppliedObjects(client, nil)
}

// Before running this test, make sure the assist auth configuration is enabled.
func TestAccAssistAuthConfigurationObjectManagement_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_workspace_assist_auth_configuration_object_management.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, GetAssistAuthConfigurationObjectManagementResourceFunc)

		name       = acceptance.RandomAccResourceNameWithDash()
		baseConfig = testAccAssistAuthConfigurationObjectManagement_base(name)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// Acceptance test will not record the origin objects, so we cannot filter the objects which managed by provider.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAssistAuthConfigurationObjectManagement_basic_step1(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				Config: testAccAssistAuthConfigurationObjectManagement_basic_step2(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"objects_origin",
					"objects",
				},
			},
		},
	})
}

func testAccAssistAuthConfigurationObjectManagement_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_user" "test" {
  count = 9

  name  = format("%[1]s-user-%%d", count.index)
  email = format("test%%d@example.com", count.index)
}

resource "huaweicloud_workspace_user_group" "test" {
  count = 3

  name = format("%[1]s-group-%%d", count.index)
  type = "LOCAL"

  dynamic "users" {
    for_each = slice(huaweicloud_workspace_user.test[*].id, count.index*2+3, count.index*2+5)

    content {
      id = users.value
    }
  }
}
`, name)
}

func testAccAssistAuthConfigurationObjectManagement_basic_step1(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_assist_auth_configuration_object_management" "test" {
  dynamic "objects" {
    for_each = slice(huaweicloud_workspace_user.test, 0, 2)

    content {
      type = "USER"
      id   = objects.value.id
      name = objects.value.name
    }
  }

  dynamic "objects" {
    for_each = slice(huaweicloud_workspace_user_group.test, 0, 2)

    content {
      type = "USER_GROUP"
      id   = objects.value.id
      name = objects.value.name
    }
  }
}
`, baseConfig)
}

func testAccAssistAuthConfigurationObjectManagement_basic_step2(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_assist_auth_configuration_object_management" "test" {
  dynamic "objects" {
    for_each = slice(huaweicloud_workspace_user.test, 1, 3)

    content {
      type = "USER"
      id   = objects.value.id
      name = objects.value.name
    }
  }

  dynamic "objects" {
    for_each = slice(huaweicloud_workspace_user_group.test, 1, 3)

    content {
      type = "USER_GROUP"
      id   = objects.value.id
      name = objects.value.name
    }
  }
}
`, baseConfig)
}
