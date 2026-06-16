package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataFactoryResources_basic(t *testing.T) {
	var (
		rName      = acceptance.RandomAccResourceName()
		dataSource = "data.huaweicloud_dataarts_factory_resources.all"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckOBSObjectStoragePath(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataFactoryResources_nonExistentWorkspace(),
				ExpectError: regexp.MustCompile("Operation failed, detail msg Workspace does not exists."),
			},
			{
				Config: testAccDataFactoryResources_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "resources.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_resource_queried", "true"),
					resource.TestCheckOutput("is_resource_name_correct", "true"),
					resource.TestCheckOutput("is_resource_type_correct", "true"),
					resource.TestCheckOutput("is_resource_description_correct", "true"),
					resource.TestCheckOutput("is_resource_directory_correct", "true"),
					resource.TestCheckOutput("is_resource_location_correct", "true"),
				),
			},
		},
	})
}

func testAccDataFactoryResources_nonExistentWorkspace() string {
	randUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
data "huaweicloud_dataarts_factory_resources" "test" {
  workspace_id = "%[1]s"
}
`, randUUID.String())
}

func testAccDataFactoryResources_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_factory_resource" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  type         = "jar"
  description  = "Created by terraform script"
  directory    = "/"
  location     = "%[3]s"
}

data "huaweicloud_dataarts_factory_resources" "all" {
  depends_on = [
    huaweicloud_dataarts_factory_resource.test,
  ]

  workspace_id = "%[1]s"
}

locals {
  id_filter_result = [
    for v in data.huaweicloud_dataarts_factory_resources.all.resources : v if v.name == huaweicloud_dataarts_factory_resource.test.name
  ]
}

output "is_resource_queried" {
  value = length(local.id_filter_result) > 0
}

output "is_resource_name_correct" {
  value = local.id_filter_result[0].name == huaweicloud_dataarts_factory_resource.test.name
}

output "is_resource_type_correct" {
  value = local.id_filter_result[0].type == huaweicloud_dataarts_factory_resource.test.type
}

output "is_resource_description_correct" {
  value = local.id_filter_result[0].description == huaweicloud_dataarts_factory_resource.test.description
}

output "is_resource_directory_correct" {
  value = local.id_filter_result[0].directory == huaweicloud_dataarts_factory_resource.test.directory
}

output "is_resource_location_correct" {
  value = local.id_filter_result[0].location == huaweicloud_dataarts_factory_resource.test.location
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name, acceptance.HW_OBS_OBJECT_STORAGE_PATH)
}
