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
		base = "huaweicloud_fgs_application.test"

		all          = "data.huaweicloud_fgs_applications.all"
		dcForAllApps = acceptance.InitDataSourceCheck(all)

		byAppId           = "data.huaweicloud_fgs_applications.filter_by_application_id"
		dcByAppId         = acceptance.InitDataSourceCheck(byAppId)
		byNotFoundAppId   = "data.huaweicloud_fgs_applications.filter_by_not_found_application_id"
		dcByNotFoundAppId = acceptance.InitDataSourceCheck(byNotFoundAppId)

		byAppName           = "data.huaweicloud_fgs_applications.filter_by_name"
		dcByAppName         = acceptance.InitDataSourceCheck(byAppName)
		byNotFoundAppName   = "data.huaweicloud_fgs_applications.filter_by_not_found_name"
		dcByNotFoundAppName = acceptance.InitDataSourceCheck(byNotFoundAppName)

		byAppStatus           = "data.huaweicloud_fgs_applications.filter_by_status"
		dcByAppStatus         = acceptance.InitDataSourceCheck(byAppStatus)
		byNotFoundAppStatus   = "data.huaweicloud_fgs_applications.filter_by_not_found_status"
		dcByNotFoundAppStatus = acceptance.InitDataSourceCheck(byNotFoundAppStatus)

		byAppDesc           = "data.huaweicloud_fgs_applications.filter_by_description"
		dcByAppDesc         = acceptance.InitDataSourceCheck(byAppDesc)
		byNotFoundAppDesc   = "data.huaweicloud_fgs_applications.filter_by_not_found_description"
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
					// Without filter parameters.
					dcForAllApps.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "applications.#", regexp.MustCompile(`[1-9][0-9]*`)),
					// Filter by application ID.
					dcByAppId.CheckResourceExists(),
					resource.TestCheckOutput("is_application_id_filter_useful", "true"),
					dcByNotFoundAppId.CheckResourceExists(),
					resource.TestCheckOutput("application_id_not_found_validation_pass", "true"),
					// Filter by application name.
					dcByAppName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByNotFoundAppName.CheckResourceExists(),
					resource.TestCheckOutput("name_not_found_validation_pass", "true"),
					// Filter by application status.
					dcByAppStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					dcByNotFoundAppStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_not_found_validation_pass", "true"),
					// Filter by application description.
					dcByAppDesc.CheckResourceExists(),
					resource.TestCheckOutput("is_description_filter_useful", "true"),
					dcByNotFoundAppDesc.CheckResourceExists(),
					resource.TestCheckOutput("description_not_found_validation_pass", "true"),
					// Check the attributes.
					resource.TestCheckResourceAttrPair(byAppId, "applications.0.id", base, "id"),
					resource.TestCheckResourceAttrPair(byAppId, "applications.0.name", base, "name"),
					resource.TestCheckResourceAttrPair(byAppId, "applications.0.status", base, "status"),
					resource.TestCheckResourceAttrPair(byAppId, "applications.0.description", base, "description"),
					resource.TestMatchResourceAttr(byAppId, "applications.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Attribute 'created_at' has been deprecated from the API and the corresponding documentation.
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
  # Query applications after application resource create.
  depends_on = [
    huaweicloud_fgs_application.test,
  ]
}

# Filter by application ID.
locals {
  application_id = huaweicloud_fgs_application.test.id
}

data "huaweicloud_fgs_applications" "filter_by_application_id" {
  application_id = local.application_id
}

data "huaweicloud_fgs_applications" "filter_by_not_found_application_id" {
  # Query applications using a not exist ID after application resource create.
  depends_on = [
    huaweicloud_fgs_application.test,
  ]

  application_id = "%[2]s"
}

locals {
  application_id_filter_result = [for v in data.huaweicloud_fgs_applications.filter_by_application_id.applications[*].id :
    v == local.application_id]
}

output "is_application_id_filter_useful" {
  value = length(local.application_id_filter_result) > 0 && alltrue(local.application_id_filter_result)
}

output "application_id_not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_applications.filter_by_not_found_application_id.applications) == 0
}

# Filter by application name.
locals {
  application_name = huaweicloud_fgs_application.test.name
}

data "huaweicloud_fgs_applications" "filter_by_name" {
  # The behavior of parameter 'name' of the application resource is 'Required', means this parameter does not
  # have 'Know After Apply' behavior.
  depends_on = [
    huaweicloud_fgs_application.test,
  ]

  name = local.application_name
}

data "huaweicloud_fgs_applications" "filter_by_not_found_name" {
  # Query applications using a not exist name after application resource create.
  depends_on = [
    huaweicloud_fgs_application.test,
  ]

  name = "name_not_found"
}

locals {
  name_filter_result = [for v in data.huaweicloud_fgs_applications.filter_by_name.applications[*].name :
    v == local.application_name]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

output "name_not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_applications.filter_by_not_found_name.applications) == 0
}

# Filter by application status.
locals {
  application_status = huaweicloud_fgs_application.test.status
}

data "huaweicloud_fgs_applications" "filter_by_status" {
  status = local.application_status
}

data "huaweicloud_fgs_applications" "filter_by_not_found_status" {
  # Query application using a not exist status after application resource create.
  depends_on = [
    huaweicloud_fgs_application.test,
  ]

  status = "status_not_found"
}

locals {
  status_filter_result = [for v in data.huaweicloud_fgs_applications.filter_by_status.applications[*].status :
    v == local.application_status]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

output "status_not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_applications.filter_by_not_found_status.applications) == 0
}

# Filter by application description.
locals {
  application_description = huaweicloud_fgs_application.test.description
}

data "huaweicloud_fgs_applications" "filter_by_description" {
  # The behavior of parameter 'description' of the resource is not have the 'Computed', means this parameter does not
  # have 'Know After Apply' behavior.
  depends_on = [
    huaweicloud_fgs_application.test,
  ]

  description = local.application_description
}

data "huaweicloud_fgs_applications" "filter_by_not_found_description" {
  # Query applications using a not exist description after application resource create.
  depends_on = [
    huaweicloud_fgs_application.test,
  ]

  description = "description_not_found"
}

locals {
  description_filter_result = [for v in data.huaweicloud_fgs_applications.filter_by_description.applications[*].description :
    v == local.application_description]
}

output "is_description_filter_useful" {
  value = length(local.description_filter_result) > 0 && alltrue(local.description_filter_result)
}

output "description_not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_applications.filter_by_not_found_description.applications) == 0
}
`, testAccApplication_basic(name), randAppId)
}
