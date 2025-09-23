package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAppGroupAuthorizations_basic(t *testing.T) {
	var (
		rName      = acceptance.RandomAccResourceName()
		dataSource = "data.huaweicloud_workspace_app_group_authorizations.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byAppGroupId   = "data.huaweicloud_workspace_app_group_authorizations.filter_by_app_group_id"
		dcByAppGroupId = acceptance.InitDataSourceCheck(byAppGroupId)

		byAccount   = "data.huaweicloud_workspace_app_group_authorizations.filter_by_account"
		dcByAccount = acceptance.InitDataSourceCheck(byAccount)

		byAccountType   = "data.huaweicloud_workspace_app_group_authorizations.filter_by_account_type"
		dcByAccountType = acceptance.InitDataSourceCheck(byAccountType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroup(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAppGroupAuthorizations_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "authorizations.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByAppGroupId.CheckResourceExists(),
					resource.TestCheckOutput("is_app_group_id_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(byAppGroupId, "authorizations.0.id"),
					resource.TestCheckResourceAttrSet(byAppGroupId, "authorizations.0.account_id"),
					resource.TestCheckResourceAttrSet(byAppGroupId, "authorizations.0.app_group_id"),
					resource.TestCheckResourceAttrSet(byAppGroupId, "authorizations.0.app_group_name"),
					resource.TestMatchResourceAttr(byAppGroupId, "authorizations.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					dcByAccount.CheckResourceExists(),
					resource.TestCheckOutput("is_account_filter_useful", "true"),
					dcByAccountType.CheckResourceExists(),
					resource.TestCheckOutput("is_account_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceAppGroupAuthorizations_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_workspace_app_group_authorizations" "test" {
  depends_on = [huaweicloud_workspace_app_group_authorization.test]
}

locals {
  app_group_id            = huaweicloud_workspace_app_group.test.id
  authorized_account      = huaweicloud_workspace_app_group_authorization.test.accounts[0].account
  authorized_account_type = huaweicloud_workspace_app_group_authorization.test.accounts[0].type
}

data "huaweicloud_workspace_app_group_authorizations" "filter_by_app_group_id" {
  depends_on = [huaweicloud_workspace_app_group_authorization.test]

  app_group_id = local.app_group_id
}

locals {
  app_group_id_filter_result = [
    for v in data.huaweicloud_workspace_app_group_authorizations.filter_by_app_group_id.authorizations[*].app_group_id : v == local.app_group_id
  ]
}

output "is_app_group_id_filter_useful" {
  value = length(local.app_group_id_filter_result) > 0 && alltrue(local.app_group_id_filter_result)
}

# Fuzzy search is supported.
data "huaweicloud_workspace_app_group_authorizations" "filter_by_account" {
  depends_on = [huaweicloud_workspace_app_group_authorization.test]

  account = local.authorized_account
}

locals {
  account_filter_result = [
    for v in data.huaweicloud_workspace_app_group_authorizations.filter_by_account.authorizations[*].account :
    strcontains(v, local.authorized_account)
  ]
}

output "is_account_filter_useful" {
  value = length(local.account_filter_result) > 0 && alltrue(local.account_filter_result)
}

data "huaweicloud_workspace_app_group_authorizations" "filter_by_account_type" {
  depends_on = [huaweicloud_workspace_app_group_authorization.test]

  account_type = local.authorized_account_type
}

locals {
  account_type_filter_result = [
    for v in data.huaweicloud_workspace_app_group_authorizations.filter_by_account_type.authorizations[*].account_type :
	v == local.authorized_account_type
  ]
}
  
output "is_account_type_filter_useful" {
  value = length(local.account_type_filter_result) > 0 && alltrue(local.account_type_filter_result)
}
`, testAccAppGroupAuthorization_basic(name))
}
