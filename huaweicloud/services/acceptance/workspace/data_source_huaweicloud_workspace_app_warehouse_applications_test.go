package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAppWarehouseApplications_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_workspace_app_warehouse_applications.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByAppId   = "data.huaweicloud_workspace_app_warehouse_applications.filter_by_app_id"
		dcFilterByAppId = acceptance.InitDataSourceCheck(filterByAppId)

		filterByName   = "data.huaweicloud_workspace_app_warehouse_applications.filter_by_name"
		dcFilterByName = acceptance.InitDataSourceCheck(filterByName)

		filterByCategory   = "data.huaweicloud_workspace_app_warehouse_applications.filter_by_category"
		dcFilterByCategory = acceptance.InitDataSourceCheck(filterByCategory)

		filterByVerifyStatus   = "data.huaweicloud_workspace_app_warehouse_applications.filter_by_verify_status"
		dcFilterByVerifyStatus = acceptance.InitDataSourceCheck(filterByVerifyStatus)
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
				Config: testAccDataAppWarehouseApplications_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "applications.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by 'app_id' parameter.
					dcFilterByAppId.CheckResourceExists(),
					resource.TestCheckOutput("is_app_id_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(filterByAppId, "applications.0.id"),
					resource.TestCheckResourceAttrSet(filterByAppId, "applications.0.app_id"),
					resource.TestCheckResourceAttrSet(filterByAppId, "applications.0.name"),
					resource.TestCheckResourceAttrSet(filterByAppId, "applications.0.category"),
					resource.TestCheckResourceAttrSet(filterByAppId, "applications.0.os_type"),
					resource.TestCheckResourceAttrSet(filterByAppId, "applications.0.version"),
					resource.TestCheckResourceAttrSet(filterByAppId, "applications.0.version_name"),
					resource.TestCheckResourceAttr(filterByAppId, "applications.0.file_store_path", acceptance.HW_WORKSPACE_APP_FILE_NAME),
					resource.TestCheckResourceAttrSet(filterByAppId, "applications.0.app_file_size"),
					resource.TestCheckResourceAttrSet(filterByAppId, "applications.0.description"),
					resource.TestCheckResourceAttrSet(filterByAppId, "applications.0.verify_status"),
					resource.TestMatchResourceAttr(filterByAppId, "applications.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(filterByAppId, "applications.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Filter by 'name' parameter (Fuzzy search).
					dcFilterByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					// Filter by 'category' parameter.
					dcFilterByCategory.CheckResourceExists(),
					resource.TestCheckOutput("is_category_filter_useful", "true"),
					// Filter by 'verify_status' parameter.
					dcFilterByVerifyStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_verify_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataAppWarehouseApplications_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameter.
data "huaweicloud_workspace_app_warehouse_applications" "all" {
  depends_on = [huaweicloud_workspace_app_warehouse_application.test]
}

locals {
  application_id = huaweicloud_workspace_app_warehouse_application.test.id
  name           = huaweicloud_workspace_app_warehouse_application.test.name
  category       = huaweicloud_workspace_app_warehouse_application.test.category
  verify_status  = try(data.huaweicloud_workspace_app_warehouse_applications.all.applications[0].verify_status, null)
}

# Filter by 'app_id' parameter.
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

# Filter by 'name' parameter (Fuzzy search).
data "huaweicloud_workspace_app_warehouse_applications" "filter_by_name" {
  name       = "tf_test"
  depends_on = [huaweicloud_workspace_app_warehouse_application.test]
}

locals {
  name_filter_result = [for v in data.huaweicloud_workspace_app_warehouse_applications.filter_by_name.applications :
  strcontains(v.name, "tf_test")]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by 'category' parameter.
data "huaweicloud_workspace_app_warehouse_applications" "filter_by_category" {
  category   = local.category
  depends_on = [huaweicloud_workspace_app_warehouse_application.test]
}

locals {
  category_filter_result = [for v in data.huaweicloud_workspace_app_warehouse_applications.filter_by_category.applications :
  v.category == local.category]
}

output "is_category_filter_useful" {
  value = length(local.category_filter_result) > 0 && alltrue(local.category_filter_result)
}

# Filter by 'verify_status' parameter.
data "huaweicloud_workspace_app_warehouse_applications" "filter_by_verify_status" {
  verify_status = local.verify_status
  depends_on    = [huaweicloud_workspace_app_warehouse_application.test]
}

locals {
  verify_status_filter_result = [for v in data.huaweicloud_workspace_app_warehouse_applications.filter_by_verify_status.applications :
  v.verify_status == local.verify_status]
}

output "is_verify_status_filter_useful" {
  value = length(local.verify_status_filter_result) > 0 && alltrue(local.verify_status_filter_result)
}
`, testAccAppWarehouseApplication_basic_step1(name))
}
