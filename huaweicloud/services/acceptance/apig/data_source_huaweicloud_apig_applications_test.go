package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceApplications_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_apig_applications.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		rName      = acceptance.RandomAccResourceName()

		byId   = "data.huaweicloud_apig_applications.filter_by_id"
		dcById = acceptance.InitDataSourceCheck(byId)

		byName   = "data.huaweicloud_apig_applications.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byAppKey   = "data.huaweicloud_apig_applications.filter_by_app_key"
		dcByAppKey = acceptance.InitDataSourceCheck(byAppKey)

		byCreatedBy   = "data.huaweicloud_apig_applications.filter_by_created_by"
		dcByCreatedBy = acceptance.InitDataSourceCheck(byCreatedBy)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceApplications_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "applications.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.app_key"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.created_by"),
					dcById.CheckResourceExists(),
					resource.TestCheckOutput("application_id_filter_is_useful", "true"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					dcByAppKey.CheckResourceExists(),
					resource.TestCheckOutput("app_key_filter_is_useful", "true"),
					dcByCreatedBy.CheckResourceExists(),
					resource.TestCheckOutput("created_by_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceApplications_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_apig_instances" "test" {
  instance_id = "%[1]s"
}

locals {
  instance_id = data.huaweicloud_apig_instances.test.instances[0].id
}

resource "huaweicloud_apig_application" "test" {
  instance_id = local.instance_id
  name        = "%[2]s"
}

data "huaweicloud_apig_applications" "test" {
  depends_on = [
    huaweicloud_apig_application.test
  ]

  instance_id = local.instance_id
}

# Filter by ID
locals {
  application_id = huaweicloud_apig_application.test.id
}

data "huaweicloud_apig_applications" "filter_by_id" {
  instance_id    = local.instance_id
  application_id = local.application_id
}

locals {
  id_filter_result = [
    for v in data.huaweicloud_apig_applications.filter_by_id.applications[*].id : v == local.application_id
  ]
}

output "application_id_filter_is_useful" {
  value = length(local.id_filter_result) > 0 && alltrue(local.id_filter_result)
}

# Filter by name
locals {
  name = huaweicloud_apig_application.test.name
}

data "huaweicloud_apig_applications" "filter_by_name" {
  depends_on = [
    huaweicloud_apig_application.test
  ]

  instance_id = local.instance_id
  name        = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_apig_applications.filter_by_name.applications[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by app_key
locals {
  app_key = huaweicloud_apig_application.test.app_key
}

data "huaweicloud_apig_applications" "filter_by_app_key" {
  instance_id = local.instance_id
  app_key     = local.app_key
}

locals {
  app_key_filter_result = [
    for v in data.huaweicloud_apig_applications.filter_by_app_key.applications[*].app_key : v == local.app_key
  ]
}

output "app_key_filter_is_useful" {
  value = length(local.app_key_filter_result) > 0 && alltrue(local.app_key_filter_result)
}

# Filter by created_by
locals {
  created_by = data.huaweicloud_apig_applications.filter_by_id.applications[0].created_by
}

data "huaweicloud_apig_applications" "filter_by_created_by" {
  instance_id = local.instance_id
  created_by  = local.created_by
}

locals {
  created_by_filter_result = [
    for v in data.huaweicloud_apig_applications.filter_by_created_by.applications[*].created_by : v == local.created_by
  ]
}

output "created_by_filter_is_useful" {
  value = length(local.created_by_filter_result) > 0 && alltrue(local.created_by_filter_result)
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}
