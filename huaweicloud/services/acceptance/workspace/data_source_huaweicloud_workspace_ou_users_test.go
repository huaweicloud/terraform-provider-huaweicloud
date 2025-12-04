package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running the test, ensure that there is at least one user under this OU.
func TestAccDataOuUsers_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_workspace_ou_users.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByUserName   = "data.huaweicloud_workspace_ou_users.filter_by_user_name"
		dcFilterByUserName = acceptance.InitDataSourceCheck(filterByUserName)

		filterByHasExisted   = "data.huaweicloud_workspace_ou_users.filter_by_has_existed"
		dcFilterByHasExisted = acceptance.InitDataSourceCheck(filterByHasExisted)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceOUName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataOuUsers_basic(),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "users.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "enable_create_count"),
					resource.TestCheckResourceAttrSet(all, "users.0.name"),
					resource.TestCheckResourceAttrSet(all, "users.0.expired_time"),
					// Filter by 'user_name' parameter.
					dcFilterByUserName.CheckResourceExists(),
					resource.TestCheckOutput("is_user_name_filter_useful", "true"),
					// Filter by 'has_existed' parameter.
					dcFilterByHasExisted.CheckResourceExists(),
					resource.TestCheckOutput("is_has_existed_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataOuUsers_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_ous" "test" {
  ou_name = "%[1]s"
}

locals {
  ou_dn = try(data.huaweicloud_workspace_ous.test.ous[0].ou_dn, "NOT_FOUND")
}

# Without any filter parameter.
data "huaweicloud_workspace_ou_users" "all" {
  ou_dn = local.ou_dn
}

# Filter by 'user_name' parameter.
locals {
  user_name = try(data.huaweicloud_workspace_ou_users.all.users[0].name, "NOT_FOUND")
}

data "huaweicloud_workspace_ou_users" "filter_by_user_name" {
  ou_dn     = local.ou_dn
  user_name = local.user_name
}

locals {
  user_name_filter_result = [for v in data.huaweicloud_workspace_ou_users.filter_by_user_name.users[*].name :
    strcontains(v, local.user_name)
  ]
}

output "is_user_name_filter_useful" {
  value = length(local.user_name_filter_result) > 0 && alltrue(local.user_name_filter_result)
}

# Filter by 'has_existed' parameter.
locals {
  has_existed = try(data.huaweicloud_workspace_ou_users.all.users[0].has_existed, null)
}

data "huaweicloud_workspace_ou_users" "filter_by_has_existed" {
  ou_dn       = local.ou_dn
  has_existed = local.has_existed
}

locals {
  has_existed_filter_result = [for v in data.huaweicloud_workspace_ou_users.filter_by_has_existed.users[*].has_existed :
    v == local.has_existed
  ]
}

output "is_has_existed_filter_useful" {
  value = length(local.has_existed_filter_result) > 0 && alltrue(local.has_existed_filter_result)
}
`, acceptance.HW_WORKSPACE_OU_NAME)
}
