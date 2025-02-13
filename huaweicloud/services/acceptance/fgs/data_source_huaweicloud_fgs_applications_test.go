package fgs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataApplications_basic(t *testing.T) {
	var (
		rcName       = "huaweicloud_fgs_application.test"
		all          = "data.huaweicloud_fgs_applications.all"
		dcForAllApps = acceptance.InitDataSourceCheck(all)

		byAppId           = "data.huaweicloud_fgs_applications.filter_by_app_id"
		dcByAppId         = acceptance.InitDataSourceCheck(byAppId)
		byNotFoundAppId   = "data.huaweicloud_fgs_applications.filter_by_not_found_app_id"
		dcByNotFoundAppId = acceptance.InitDataSourceCheck(byNotFoundAppId)

		byAppName           = "data.huaweicloud_fgs_applications.filter_by_app_name"
		dcByAppName         = acceptance.InitDataSourceCheck(byAppName)
		byNotFoundAppName   = "data.huaweicloud_fgs_applications.filter_by_not_found_app_name"
		dcByNotFoundAppName = acceptance.InitDataSourceCheck(byNotFoundAppName)

		byAppStatus           = "data.huaweicloud_fgs_applications.filter_by_app_status"
		dcByAppStatus         = acceptance.InitDataSourceCheck(byAppStatus)
		byNotFoundAppStatus   = "data.huaweicloud_fgs_applications.filter_by_not_found_app_status"
		dcByNotFoundAppStatus = acceptance.InitDataSourceCheck(byNotFoundAppStatus)

		byAppDesc           = "data.huaweicloud_fgs_applications.filter_by_app_desc"
		dcByAppDesc         = acceptance.InitDataSourceCheck(byAppDesc)
		byNotFoundAppDesc   = "data.huaweicloud_fgs_applications.filter_by_not_found_app_desc"
		dcByNotFoundAppDesc = acceptance.InitDataSourceCheck(byNotFoundAppDesc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataApplications_basic(),
				Check: resource.ComposeTestCheckFunc(
					dcForAllApps.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "applications.#", regexp.MustCompile(`[1-9][0-9]*`)),
					// Filter by application ID.
					dcByAppId.CheckResourceExists(),
					resource.TestCheckOutput("is_app_id_filter_useful", "true"),
					dcByNotFoundAppId.CheckResourceExists(),
					resource.TestCheckOutput("app_id_not_found_validation_pass", "true"),
					// Filter by application name.
					dcByAppName.CheckResourceExists(),
					resource.TestCheckOutput("is_app_name_filter_useful", "true"),
					dcByNotFoundAppName.CheckResourceExists(),
					resource.TestCheckOutput("app_name_not_found_validation_pass", "true"),
					// Filter by application status.
					dcByAppStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_app_status_filter_useful", "true"),
					dcByNotFoundAppStatus.CheckResourceExists(),
					resource.TestCheckOutput("app_status_not_found_validation_pass", "true"),
					// Filter by application description.
					dcByAppDesc.CheckResourceExists(),
					resource.TestCheckOutput("is_app_desc_filter_useful", "true"),
					dcByNotFoundAppDesc.CheckResourceExists(),
					resource.TestCheckOutput("app_desc_not_found_validation_pass", "true"),
					// Check attributes.
					resource.TestCheckResourceAttrPair(byAppId, "applications.0.id", rcName, "id"),
					resource.TestCheckResourceAttrPair(byAppId, "applications.0.name", rcName, "name"),
					resource.TestCheckResourceAttrPair(byAppId, "applications.0.status", rcName, "status"),
					resource.TestCheckResourceAttrPair(byAppId, "applications.0.description", rcName, "description"),
					resource.TestMatchResourceAttr(byAppId, "applications.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccDataApplications_basic() string {
	name := acceptance.RandomAccResourceName()
	randAppId, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
%[1]s

# Without any filter parameter.
data "huaweicloud_fgs_applications" "all" {
  // Query applications after application resource create.
  depends_on = [
    huaweicloud_fgs_application.test,
  ]
}

// Filter by application ID.
locals {
  app_id = huaweicloud_fgs_application.test.id
}

data "huaweicloud_fgs_applications" "filter_by_app_id" {
  application_id = local.app_id
}

data "huaweicloud_fgs_applications" "filter_by_not_found_app_id" {
  // Query applications using a not exist ID after application resource create.
  depends_on = [
    huaweicloud_fgs_application.test,
  ]

  application_id = "%[2]s"
}

locals {
  app_id_filter_result = [for v in data.huaweicloud_fgs_applications.filter_by_app_id.applications[*].id :
    v == local.app_id]
}

output "is_app_id_filter_useful" {
  value = length(local.app_id_filter_result) > 0 && alltrue(local.app_id_filter_result)
}

output "app_id_not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_applications.filter_by_not_found_app_id.applications) == 0
}

// Filter by application name.
locals {
  app_name = huaweicloud_fgs_application.test.name
}

data "huaweicloud_fgs_applications" "filter_by_app_name" {
  // The behavior of parameter 'name' of the application resource is 'Required', means this parameter does not
  // have 'Know After Apply' behavior.
  depends_on = [
    huaweicloud_fgs_application.test,
  ]

  name = local.app_name
}

data "huaweicloud_fgs_applications" "filter_by_not_found_app_name" {
  // Query applications using a not exist name after application resource create.
  depends_on = [
    huaweicloud_fgs_application.test,
  ]

  name = "app_name_not_found"
}

locals {
  app_name_filter_result = [for v in data.huaweicloud_fgs_applications.filter_by_app_name.applications[*].name :
    v == local.app_name]
}

output "is_app_name_filter_useful" {
  value = length(local.app_name_filter_result) > 0 && alltrue(local.app_name_filter_result)
}

output "app_name_not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_applications.filter_by_not_found_app_name.applications) == 0
}

// Filter by application status.
locals {
  app_status = huaweicloud_fgs_application.test.status
}

data "huaweicloud_fgs_applications" "filter_by_app_status" {
  status = local.app_status
}

data "huaweicloud_fgs_applications" "filter_by_not_found_app_status" {
  // Query application using a not exist status after application resource create.
  depends_on = [
    huaweicloud_fgs_application.test,
  ]

  status = "app_status_not_found"
}

locals {
  app_status_filter_result = [for v in data.huaweicloud_fgs_applications.filter_by_app_status.applications[*].status :
    v == local.app_status]
}

output "is_app_status_filter_useful" {
  value = length(local.app_status_filter_result) > 0 && alltrue(local.app_status_filter_result)
}

output "app_status_not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_applications.filter_by_not_found_app_status.applications) == 0
}

// Filter by application description.
locals {
  app_desc = huaweicloud_fgs_application.test.description
}

data "huaweicloud_fgs_applications" "filter_by_app_desc" {
  // The behavior of parameter 'description' of the resource is not have the 'Computed', means this parameter does not
  // have 'Know After Apply' behavior.
  depends_on = [
    huaweicloud_fgs_application.test,
  ]

  description = local.app_desc
}

data "huaweicloud_fgs_applications" "filter_by_not_found_app_desc" {
  // Query applications using a not exist description after application resource create.
  depends_on = [
    huaweicloud_fgs_application.test,
  ]

  description = "app_desc_not_found"
}

locals {
  app_desc_filter_result = [for v in data.huaweicloud_fgs_applications.filter_by_app_desc.applications[*].description :
    v == local.app_desc]
}

output "is_app_desc_filter_useful" {
  value = length(local.app_desc_filter_result) > 0 && alltrue(local.app_desc_filter_result)
}

output "app_desc_not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_applications.filter_by_not_found_app_desc.applications) == 0
}
`, testAccApplication_basic(name), randAppId)
}
