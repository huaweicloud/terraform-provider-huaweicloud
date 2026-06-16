package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataFactoryScripts_basic(t *testing.T) {
	var (
		rName = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dataarts_factory_scripts.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByName   = "data.huaweicloud_dataarts_factory_scripts.filter_by_name"
		dcFilterByName = acceptance.InitDataSourceCheck(filterByName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsConnectionName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataFactoryScripts_nonExistentWorkspace(),
				ExpectError: regexp.MustCompile("Workspace does not exists"),
			},
			{
				Config: testAccDataFactoryScripts_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "scripts.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcFilterByName.CheckResourceExists(),
					resource.TestCheckOutput("is_script_name_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(filterByName, "scripts.0.id"),
					resource.TestCheckResourceAttrPair(filterByName, "scripts.0.name",
						"huaweicloud_dataarts_factory_script.test", "name"),
					resource.TestCheckResourceAttrPair(filterByName, "scripts.0.type",
						"huaweicloud_dataarts_factory_script.test", "type"),
					resource.TestCheckResourceAttrPair(filterByName, "scripts.0.directory",
						"huaweicloud_dataarts_factory_script.test", "directory"),
					resource.TestCheckResourceAttrPair(filterByName, "scripts.0.create_user",
						"huaweicloud_dataarts_factory_script.test", "created_by"),
					resource.TestCheckResourceAttrPair(filterByName, "scripts.0.connection_name",
						"huaweicloud_dataarts_factory_script.test", "connection_name"),
					resource.TestCheckResourceAttrPair(filterByName, "scripts.0.database",
						"huaweicloud_dataarts_factory_script.test", "database"),
					resource.TestCheckResourceAttrPair(filterByName, "scripts.0.queue_name",
						"huaweicloud_dataarts_factory_script.test", "queue_name"),
					resource.TestCheckResourceAttrPair(filterByName, "scripts.0.configuration",
						"huaweicloud_dataarts_factory_script.test", "configuration"),
					resource.TestCheckResourceAttrPair(filterByName, "scripts.0.description",
						"huaweicloud_dataarts_factory_script.test", "description"),
					resource.TestCheckResourceAttrSet(filterByName, "scripts.0.modify_time"),
					resource.TestCheckResourceAttrSet(filterByName, "scripts.0.owner"),
					resource.TestCheckResourceAttrSet(filterByName, "scripts.0.version"),
				),
			},
		},
	})
}

func testAccDataFactoryScripts_nonExistentWorkspace() string {
	randUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
data "huaweicloud_dataarts_factory_scripts" "test" {
  workspace_id = "%[1]s"
}
`, randUUID.String())
}

func testAccDataFactoryScripts_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_database" "test" {
  name        = "%[1]s"
  description = "Created by Terraform Script"
}

resource "huaweicloud_dli_table" "test" {
  database_name = huaweicloud_dli_database.test.name
  name          = "%[1]s"
  data_location = "DLI"
  description   = "Created by Terraform Script"

  columns {
    name = "name"
    type = "string"
  }

  columns {
    name = "age"
    type = "int"
  }
}

resource "huaweicloud_dataarts_factory_script" "test" {
  depends_on = [
    huaweicloud_dli_database.test,
    huaweicloud_dli_table.test,
  ]

  workspace_id    = "%[2]s"
  name            = "%[1]s"
  type            = "DLISQL"
  content         = "SELECT * FROM ${huaweicloud_dli_database.test.name}.${huaweicloud_dli_table.test.name}"
  connection_name = "%[3]s"
  directory       = "/basic"
  database        = huaweicloud_dli_database.test.name
  queue_name      = "default"
  description     = "Created by Terraform Script"
}

data "huaweicloud_dataarts_factory_scripts" "all" {
  depends_on = [
    huaweicloud_dataarts_factory_script.test,
  ]

  workspace_id = "%[2]s"
}

locals {
  script_name = huaweicloud_dataarts_factory_script.test.name
}

data "huaweicloud_dataarts_factory_scripts" "filter_by_name" {
  # The behavior of parameter 'name' of the resource is 'Required', means this parameter does not 
  # have 'Know After Apply' behavior.
  depends_on = [
    huaweicloud_dataarts_factory_script.test,
  ]

  workspace_id = "%[2]s"
  script_name  = local.script_name
}

locals {
  script_name_filter_result = [
    for v in data.huaweicloud_dataarts_factory_scripts.filter_by_name.scripts : v.name == local.script_name
  ]
}

output "is_script_name_filter_useful" {
  value = length(local.script_name_filter_result) > 0 && alltrue(local.script_name_filter_result)
}
`, name, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_CONNECTION_NAME)
}
