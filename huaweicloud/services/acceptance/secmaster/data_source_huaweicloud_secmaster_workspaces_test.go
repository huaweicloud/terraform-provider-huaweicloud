package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecmasterWorkspaces_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_workspaces.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSecmasterWorkspaces_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "workspaces.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "workspaces.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "workspaces.0.description"),

					resource.TestCheckOutput("ids_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("desc_filter_is_useful", "true"),
					resource.TestCheckOutput("view_bind_id_filter_is_useful", "true"),
					resource.TestCheckOutput("view_bind_name_filter_is_useful", "true"),
					resource.TestCheckOutput("eps_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSecmasterWorkspaces_basic() string {
	randUUID, _ := uuid.GenerateUUID()
	randName := acceptance.RandomAccResourceNameWithDash()
	return fmt.Sprintf(`
data "huaweicloud_secmaster_workspaces" "test" {}

locals {
  workspace_id   = data.huaweicloud_secmaster_workspaces.test.workspaces[0].id
  name           = data.huaweicloud_secmaster_workspaces.test.workspaces[0].name
  description    = data.huaweicloud_secmaster_workspaces.test.workspaces[0].description
  view_bind_id   = data.huaweicloud_secmaster_workspaces.test.workspaces[0].view_bind_id
  view_bind_name = data.huaweicloud_secmaster_workspaces.test.workspaces[0].view_bind_name
  eps_id         = data.huaweicloud_secmaster_workspaces.test.workspaces[0].enterprise_project_id
}

data "huaweicloud_secmaster_workspaces" "filter_by_ids" {
  ids = local.workspace_id
}

data "huaweicloud_secmaster_workspaces" "filter_not_found_by_ids" {
  ids = "%[1]s"
}

data "huaweicloud_secmaster_workspaces" "filter_by_name" {
  name = local.name
}

data "huaweicloud_secmaster_workspaces" "filter_not_found_by_name" {
  name = "%[2]s"
}

data "huaweicloud_secmaster_workspaces" "filter_by_desc" {
  description = local.description
}

data "huaweicloud_secmaster_workspaces" "filter_not_found_by_desc" {
  description = "%[2]s"
}

data "huaweicloud_secmaster_workspaces" "filter_by_view_bind_id" {
  view_bind_id = local.view_bind_id
}

data "huaweicloud_secmaster_workspaces" "filter_null_by_view_bind_id" {
  view_bind_id = "%[1]s"
}

data "huaweicloud_secmaster_workspaces" "filter_by_view_bind_name" {
  view_bind_name = local.view_bind_name
}

data "huaweicloud_secmaster_workspaces" "filter_null_by_view_bind_name" {
  view_bind_name = "%[2]s"
}

data "huaweicloud_secmaster_workspaces" "filter_by_eps_id" {
  enterprise_project_id = local.eps_id
}

data "huaweicloud_secmaster_workspaces" "filter_null_by_eps_id" {
  enterprise_project_id = "%[1]s"
}

locals {
  list_by_ids                 = data.huaweicloud_secmaster_workspaces.filter_by_ids.workspaces
  list_null_by_ids            = data.huaweicloud_secmaster_workspaces.filter_not_found_by_ids.workspaces
  list_by_name                = data.huaweicloud_secmaster_workspaces.filter_by_name.workspaces
  list_null_by_name           = data.huaweicloud_secmaster_workspaces.filter_not_found_by_name.workspaces
  list_by_desc                = data.huaweicloud_secmaster_workspaces.filter_by_desc.workspaces
  list_null_by_desc           = data.huaweicloud_secmaster_workspaces.filter_not_found_by_desc.workspaces
  list_by_view_bind_id        = data.huaweicloud_secmaster_workspaces.filter_by_view_bind_id.workspaces
  list_null_by_view_bind_id   = data.huaweicloud_secmaster_workspaces.filter_null_by_view_bind_id.workspaces
  list_by_view_bind_name      = data.huaweicloud_secmaster_workspaces.filter_by_view_bind_name.workspaces
  list_null_by_view_bind_name = data.huaweicloud_secmaster_workspaces.filter_null_by_view_bind_name.workspaces
  list_by_eps_id              = data.huaweicloud_secmaster_workspaces.filter_by_eps_id.workspaces
  list_null_by_eps_id         = data.huaweicloud_secmaster_workspaces.filter_null_by_eps_id.workspaces
}

output "ids_filter_is_useful" {
  value = length(local.list_by_ids) == 1 && length(local.list_null_by_ids) == 0
}

output "name_filter_is_useful" {
  value = length(local.list_by_name) == 1 && length(local.list_null_by_name) == 0
}

output "desc_filter_is_useful" {
  value = length(local.list_by_desc) >= 1 && length(local.list_null_by_desc) == 0
}

output "view_bind_id_filter_is_useful" {
  value = length(local.list_by_view_bind_id) >= 1 && length(local.list_null_by_view_bind_id) == 0
}

output "view_bind_name_filter_is_useful" {
  value = length(local.list_by_view_bind_name) >= 1 && length(local.list_null_by_view_bind_name) == 0
}

output "eps_id_filter_is_useful" {
  value = length(local.list_by_eps_id) >= 1 && length(local.list_null_by_eps_id) == 0
}
`, randUUID, randName)
}
