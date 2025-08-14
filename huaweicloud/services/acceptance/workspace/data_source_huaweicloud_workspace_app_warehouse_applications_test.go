package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAppWarehouseApplications_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		dataSource = "data.huaweicloud_workspace_app_warehouse_applications.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byAppId   = "data.huaweicloud_workspace_app_warehouse_applications.filter_by_app_id"
		dcByAppId = acceptance.InitDataSourceCheck(byAppId)

		exactByName   = "data.huaweicloud_workspace_app_warehouse_applications.exact_filter_by_name"
		dcExactByName = acceptance.InitDataSourceCheck(exactByName)

		byName   = "data.huaweicloud_workspace_app_warehouse_applications.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byCategory   = "data.huaweicloud_workspace_app_warehouse_applications.filter_by_category"
		dcByCategory = acceptance.InitDataSourceCheck(byCategory)

		byVerifyStatus   = "data.huaweicloud_workspace_app_warehouse_applications.filter_by_verify_status"
		dcByVerifyStatus = acceptance.InitDataSourceCheck(byVerifyStatus)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppFileName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"null": {
				Source:            "hashicorp/null",
				VersionConstraint: "3.2.1",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppWarehouseApplications_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "applications.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByAppId.CheckResourceExists(),
					resource.TestCheckOutput("is_app_id_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(byAppId, "applications.0.id"),
					resource.TestCheckResourceAttrSet(byAppId, "applications.0.app_id"),
					resource.TestCheckResourceAttrSet(byAppId, "applications.0.name"),
					resource.TestCheckResourceAttrSet(byAppId, "applications.0.category"),
					resource.TestCheckResourceAttrSet(byAppId, "applications.0.os_type"),
					resource.TestCheckResourceAttrSet(byAppId, "applications.0.version"),
					resource.TestCheckResourceAttrSet(byAppId, "applications.0.version_name"),
					resource.TestCheckResourceAttr(byAppId, "applications.0.file_store_path", acceptance.HW_WORKSPACE_APP_FILE_NAME),
					resource.TestCheckResourceAttrSet(byAppId, "applications.0.app_file_size"),
					resource.TestCheckResourceAttrSet(byAppId, "applications.0.description"),
					resource.TestCheckResourceAttrSet(byAppId, "applications.0.verify_status"),
					resource.TestMatchResourceAttr(byAppId, "applications.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(byAppId, "applications.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					dcExactByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_exact_filter_useful", "true"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByCategory.CheckResourceExists(),
					resource.TestCheckOutput("is_category_filter_useful", "true"),
					dcByVerifyStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_verify_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAppWarehouseApplications_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_workspace_app_warehouse_applications" "test" {
  depends_on = [huaweicloud_workspace_app_warehouse_app.test]
}

locals {
  application_id = huaweicloud_workspace_app_warehouse_app.test.id
  name           = huaweicloud_workspace_app_warehouse_app.test.name
  category       = huaweicloud_workspace_app_warehouse_app.test.category
  verify_status  = try(data.huaweicloud_workspace_app_warehouse_applications.test.applications[0].verify_status, null)
}

# Filter by application ID.
data "huaweicloud_workspace_app_warehouse_applications" "filter_by_app_id" {
  app_id = local.application_id
}

locals {
  app_id_filter_result = [for v in data.huaweicloud_workspace_app_warehouse_applications.filter_by_app_id.applications :
  v.app_id == local.application_id]
}

output "is_app_id_filter_useful" {
  value = length(local.app_id_filter_result) > 0 && alltrue(local.app_id_filter_result)
}

# Filter by application name (Exact match).
data "huaweicloud_workspace_app_warehouse_applications" "exact_filter_by_name" {
  name       = local.name
  depends_on = [huaweicloud_workspace_app_warehouse_app.test]
}

locals {
  name_exact_filter_result = [for v in data.huaweicloud_workspace_app_warehouse_applications.exact_filter_by_name.applications : v.name == local.name]
}

output "is_name_exact_filter_useful" {
  value = length(local.name_exact_filter_result) > 0 && alltrue(local.name_exact_filter_result)
}

# Filter by application name (Fuzzy search).
data "huaweicloud_workspace_app_warehouse_applications" "filter_by_name" {
  name       = "tf_test"
  depends_on = [huaweicloud_workspace_app_warehouse_app.test]
}

locals {
  name_filter_result = [for v in data.huaweicloud_workspace_app_warehouse_applications.filter_by_name.applications :
  strcontains(v.name, "tf_test")]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by application category.
data "huaweicloud_workspace_app_warehouse_applications" "filter_by_category" {
  category   = local.category
  depends_on = [huaweicloud_workspace_app_warehouse_app.test]
}

locals {
  category_filter_result = [for v in data.huaweicloud_workspace_app_warehouse_applications.filter_by_category.applications :
  v.category == local.category]
}

output "is_category_filter_useful" {
  value = length(local.category_filter_result) > 0 && alltrue(local.category_filter_result)
}

# Filter by application verify status.
data "huaweicloud_workspace_app_warehouse_applications" "filter_by_verify_status" {
  verify_status = local.verify_status
  depends_on    = [huaweicloud_workspace_app_warehouse_app.test]
}

locals {
  verify_status_filter_result = [for v in data.huaweicloud_workspace_app_warehouse_applications.filter_by_verify_status.applications :
  v.verify_status == local.verify_status]
}

output "is_verify_status_filter_useful" {
  value = length(local.verify_status_filter_result) > 0 && alltrue(local.verify_status_filter_result)
}
`, testAccWarehouseApp_basic_step1(name))
}
