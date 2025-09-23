package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOperationConnections_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_operation_connections.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceOperationConnections_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.workspace_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.component_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.component_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.component_version_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.config"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.creator_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.creator_name"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_component_name_filter_useful", "true"),
					resource.TestCheckOutput("is_creator_name_filter_useful", "true"),
					resource.TestCheckOutput("is_description_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testDataSourceOperationConnections_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_operation_connections" "test" {
  workspace_id = "%[1]s"
}

# Filter using name.
locals {
  name = data.huaweicloud_secmaster_operation_connections.test.data[0].name
}

data "huaweicloud_secmaster_operation_connections" "name_filter" {
  workspace_id = "%[1]s"
  name         = local.name
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_secmaster_operation_connections.name_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_operation_connections.name_filter.data[*].name : v == local.name]
  )
}

# Filter using component_name.
locals {
  component_name = data.huaweicloud_secmaster_operation_connections.test.data[0].component_name
}

data "huaweicloud_secmaster_operation_connections" "component_name_filter" {
  workspace_id   = "%[1]s"
  component_name = local.component_name
}

output "is_component_name_filter_useful" {
  value = length(data.huaweicloud_secmaster_operation_connections.component_name_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_operation_connections.component_name_filter.data[*].component_name : v == local.component_name]
  )
}

# Filter using creator_name.
locals {
  creator_name = data.huaweicloud_secmaster_operation_connections.test.data[0].creator_name
}

data "huaweicloud_secmaster_operation_connections" "creator_name_filter" {
  workspace_id = "%[1]s"
  creator_name = local.creator_name
}

output "is_creator_name_filter_useful" {
  value = length(data.huaweicloud_secmaster_operation_connections.creator_name_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_operation_connections.creator_name_filter.data[*].creator_name : v == local.creator_name]
  )
}

# Filter using description.
locals {
  description = data.huaweicloud_secmaster_operation_connections.test.data[0].description
}

data "huaweicloud_secmaster_operation_connections" "description_filter" {
  workspace_id = "%[1]s"
  description  = local.description
}

output "is_description_filter_useful" {
  value = length(data.huaweicloud_secmaster_operation_connections.description_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_operation_connections.description_filter.data[*].description : v == local.description]
  )
}

# Filter using non-existent name.
data "huaweicloud_secmaster_operation_connections" "not_found" {
  workspace_id = "%[1]s"
  name         = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_secmaster_operation_connections.not_found.data) == 0
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
