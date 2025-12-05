package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this test, make sure the assist auth configuration is enabled.
func TestAccAssistAuthConfigurationObjectBatchApply_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_workspace_assist_auth_configuration_object_batch_apply.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAssistAuthConfigurationObjectBatchApply_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "add.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "add.0.object_type", "USER"),
					resource.TestCheckResourceAttrPair(resourceName, "add.0.object_id", "huaweicloud_workspace_user.test.0", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "add.0.object_name", "huaweicloud_workspace_user.test.0", "name"),
					resource.TestCheckResourceAttr(resourceName, "add.1.object_type", "USER"),
					resource.TestCheckResourceAttrPair(resourceName, "add.1.object_id", "huaweicloud_workspace_user.test.1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "add.1.object_name", "huaweicloud_workspace_user.test.1", "name"),
					resource.TestCheckResourceAttr(resourceName, "add.2.object_type", "USER_GROUP"),
					resource.TestCheckResourceAttrPair(resourceName, "add.2.object_id", "huaweicloud_workspace_user_group.test.0", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "add.2.object_name", "huaweicloud_workspace_user_group.test.0", "name"),
					resource.TestCheckResourceAttr(resourceName, "add.3.object_type", "USER_GROUP"),
					resource.TestCheckResourceAttrPair(resourceName, "add.3.object_id", "huaweicloud_workspace_user_group.test.1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "add.3.object_name", "huaweicloud_workspace_user_group.test.1", "name"),
				),
			},
		},
	})
}

func testAccAssistAuthConfigurationObjectBatchApply_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_user" "test" {
  count = 4

  name  = format("%[1]s-user-%%d", count.index)
  email = format("test%%d@example.com", count.index)
}

resource "huaweicloud_workspace_user_group" "test" {
  count = 2

  name = format("%[1]s-group-%%d", count.index)
  type = "LOCAL"

  users {
    id = huaweicloud_workspace_user.test[count.index+2].id
  }
}

resource "huaweicloud_workspace_assist_auth_configuration_object_batch_apply" "test" {
  dynamic "add" {
    for_each = slice(huaweicloud_workspace_user.test, 0, 2)

    content {
      object_type = "USER"
      object_id   = add.value.id
      object_name = add.value.name
    }
  }

  dynamic "add" {
    for_each = huaweicloud_workspace_user_group.test

    content {
      object_type = "USER_GROUP"
      object_id   = add.value.id
      object_name = add.value.name
    }
  }
}
`, name)
}
