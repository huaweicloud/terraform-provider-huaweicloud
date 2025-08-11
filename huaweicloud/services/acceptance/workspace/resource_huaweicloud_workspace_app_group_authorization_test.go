package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this test, please create a workspace APP server group with SESSION_DESKTOP_APP type.
func TestAccResourceAppGroupAuthorization_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroupId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAppGroupAuthorization_basic(name),
			},
		},
	})
}

func testAccAppGroupAuthorization_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_user" "test" {
  name  = "%[2]s"
  email = "tf@example.com"
}

resource "huaweicloud_workspace_user_group" "test" {
  count = 2

  name  = "%[2]s${count.index}"
  type  = "LOCAL"
}
`, testDataSourceAppGroups_base(name), name)
}

func testAccAppGroupAuthorization_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_group_authorization" "test" {
  app_group_id = huaweicloud_workspace_app_group.test.id

  accounts {
    account = huaweicloud_workspace_user.test.name
    type    = "USER"
  }

  dynamic "accounts" {
    for_each = huaweicloud_workspace_user_group.test[*]

    content {
      id      = accounts.value.id
      account = accounts.value.name
      type    = "USER_GROUP"
    }
  }
}
`, testAccAppGroupAuthorization_base(name))
}

func TestAccResourceAppGroupAuthorization_expectErr(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroupId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppGroupAuthorization_expectErr(name),
				ExpectError: regexp.MustCompile(`unable to authorize for some accounts: not_exist_user_group_tf | USER`),
			},
		},
	})
}

func testAccAppGroupAuthorization_expectErr(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_user" "test" {
  count = 2
  name  = "%[2]s${count.index}"
  email = "tf@example.com"
}

resource "huaweicloud_workspace_app_group_authorization" "test" {
  app_group_id = huaweicloud_workspace_app_group.test.id

  dynamic "accounts" {
    for_each = huaweicloud_workspace_user.test[*]
  
    content {
      account = accounts.value.name
      type    = "USER"
    }
  }

  accounts {
    account = "not_exist_user_group_tf"
    type    = "USER"
  }
}
`, testDataSourceAppGroups_base(name), name)
}
