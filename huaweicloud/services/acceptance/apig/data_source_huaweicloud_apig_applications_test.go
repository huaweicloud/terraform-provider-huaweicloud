package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataApplications_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_apig_applications.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byId   = "data.huaweicloud_apig_applications.filter_by_id"
		dcById = acceptance.InitDataSourceCheck(byId)

		byName   = "data.huaweicloud_apig_applications.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byAppKey   = "data.huaweicloud_apig_applications.filter_by_app_key"
		dcByAppKey = acceptance.InitDataSourceCheck(byAppKey)

		notFound   = "data.huaweicloud_apig_applications.not_found"
		dcNotFound = acceptance.InitDataSourceCheck(notFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataApplications_nonExistentInstance(),
				ExpectError: regexp.MustCompile(`error querying applications`),
			},
			{
				Config: testAccDataApplications_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "applications.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrPair(dataSource, "applications.0.id", "huaweicloud_apig_application.test", "id"),
					resource.TestCheckResourceAttrPair(dataSource, "applications.0.name", "huaweicloud_apig_application.test", "name"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.app_key"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.app_secret"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.app_type"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.created_by"),
					resource.TestMatchResourceAttr(dataSource, "applications.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(dataSource, "applications.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					dcById.CheckResourceExists(),
					resource.TestCheckResourceAttr(byId, "applications.#", "1"),
					resource.TestCheckResourceAttrPair(byId, "applications.0.id", "huaweicloud_apig_application.test", "id"),
					resource.TestCheckOutput("application_id_filter_is_useful", "true"),
					dcByName.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(byName, "applications.0.name", "huaweicloud_apig_application.test", "name"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					dcByAppKey.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(byAppKey, "applications.0.app_key", "huaweicloud_apig_application.test", "app_key"),
					resource.TestCheckOutput("app_key_filter_is_useful", "true"),
					dcNotFound.CheckResourceExists(),
					resource.TestCheckResourceAttr(notFound, "applications.#", "0"),
				),
			},
		},
	})
}

func testAccDataApplications_nonExistentInstance() string {
	randomUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
data "huaweicloud_apig_applications" "test" {
  instance_id = "%[1]s"
}
`, randomUUID.String())
}

func testAccDataApplications_basic_base() string {
	name := acceptance.RandomAccResourceName()

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
  description = "Created by acceptance test"
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccDataApplications_basic() string {
	randomUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
%s

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

output "app_key_filter_is_useful" {
  value = length(data.huaweicloud_apig_applications.filter_by_app_key.applications) > 0
}

# Filter by non-existent application ID
data "huaweicloud_apig_applications" "not_found" {
  instance_id    = local.instance_id
  application_id = "%[2]s"
}
`, testAccDataApplications_basic_base(), randomUUID.String())
}
