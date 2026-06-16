package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataArchitectureCodeTables_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dataarts_architecture_code_tables.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_dataarts_architecture_code_tables.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byCode   = "data.huaweicloud_dataarts_architecture_code_tables.filter_by_code"
		dcByCode = acceptance.InitDataSourceCheck(byCode)

		byCreateBy   = "data.huaweicloud_dataarts_architecture_code_tables.filter_by_create_by"
		dcByCreateBy = acceptance.InitDataSourceCheck(byCreateBy)

		byDirectoryId   = "data.huaweicloud_dataarts_architecture_code_tables.filter_by_directory_id"
		dcByDirectoryId = acceptance.InitDataSourceCheck(byDirectoryId)

		byStatus   = "data.huaweicloud_dataarts_architecture_code_tables.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byBeginTime   = "data.huaweicloud_dataarts_architecture_code_tables.filter_by_begin_time"
		dcByBeginTime = acceptance.InitDataSourceCheck(byBeginTime)

		byEndTime   = "data.huaweicloud_dataarts_architecture_code_tables.filter_by_end_time"
		dcByEndTime = acceptance.InitDataSourceCheck(byEndTime)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataArchitectureCodeTables_nonExistentWorkspace(),
				ExpectError: regexp.MustCompile("error querying DataArts Architecture code tables"),
			},
			{
				Config: testAccDataArchitectureCodeTables_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "code_tables.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "code_tables.0.id"),
					resource.TestCheckResourceAttrSet(all, "code_tables.0.name"),
					resource.TestCheckResourceAttrSet(all, "code_tables.0.code"),
					resource.TestCheckResourceAttrSet(all, "code_tables.0.directory_id"),
					resource.TestCheckResourceAttrSet(all, "code_tables.0.directory_path"),
					resource.TestCheckResourceAttrSet(all, "code_tables.0.description"),
					resource.TestCheckResourceAttrSet(all, "code_tables.0.created_by"),
					resource.TestMatchResourceAttr(all, "code_tables.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(all, "code_tables.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(all, "code_tables.0.status"),
					resource.TestMatchResourceAttr(all, "code_tables.0.fields.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "code_tables.0.fields.0.id"),
					resource.TestCheckResourceAttrSet(all, "code_tables.0.fields.0.name"),
					resource.TestCheckResourceAttrSet(all, "code_tables.0.fields.0.code"),
					resource.TestCheckResourceAttrSet(all, "code_tables.0.fields.0.type"),
					resource.TestCheckResourceAttrSet(all, "code_tables.0.fields.0.ordinal"),
					// Filter by name.
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),

					// Filter by code.
					dcByCode.CheckResourceExists(),
					resource.TestCheckOutput("is_code_filter_useful", "true"),

					// Filter by create by.
					dcByCreateBy.CheckResourceExists(),
					resource.TestCheckOutput("is_create_by_filter_useful", "true"),

					// Filter by directory ID.
					dcByDirectoryId.CheckResourceExists(),
					resource.TestCheckOutput("is_directory_id_filter_useful", "true"),

					// Filter by status.
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),

					// Filter by begin time.
					dcByBeginTime.CheckResourceExists(),
					resource.TestCheckOutput("is_begin_time_filter_useful", "true"),

					// Filter by end time.
					dcByEndTime.CheckResourceExists(),
					resource.TestCheckOutput("is_end_time_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataArchitectureCodeTables_nonExistentWorkspace() string {
	randUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
data "huaweicloud_dataarts_architecture_code_tables" "test" {
  workspace_id = "%[1]s"
}
`, randUUID.String())
}

func testAccDataArchitectureCodeTables_basic_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_directory" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  type         = "CODE"
}

resource "huaweicloud_dataarts_architecture_code_table" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  code         = "%[2]s_code"
  directory_id = huaweicloud_dataarts_architecture_directory.test.id
  description  = "Created by acceptance test"

  fields {
    name        = "field"
    code        = "field_code"
    type        = "BIGINT"
    description = "Created by acceptance test"
  }
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testAccDataArchitectureCodeTables_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dataarts_architecture_code_tables" "all" {
  depends_on = [huaweicloud_dataarts_architecture_code_table.test]

  workspace_id = "%[2]s"
}

# Filter by name (exact match).
locals {
  code_table_name = huaweicloud_dataarts_architecture_code_table.test.name
}

data "huaweicloud_dataarts_architecture_code_tables" "filter_by_name" {
  depends_on = [huaweicloud_dataarts_architecture_code_table.test]

  workspace_id = "%[2]s"
  name         = local.code_table_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_code_tables.filter_by_name.code_tables :
      v.name == local.code_table_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by code (exact match).
locals {
  code_table_code = huaweicloud_dataarts_architecture_code_table.test.code
}

data "huaweicloud_dataarts_architecture_code_tables" "filter_by_code" {
  depends_on = [huaweicloud_dataarts_architecture_code_table.test]

  workspace_id = "%[2]s"
  code         = local.code_table_code

}

locals {
  code_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_code_tables.filter_by_code.code_tables :
      v.code == local.code_table_code
  ]
}

output "is_code_filter_useful" {
  value = contains(data.huaweicloud_dataarts_architecture_code_tables.filter_by_code.code_tables[*].id,
    huaweicloud_dataarts_architecture_code_table.test.id)
}

# Filter by create by.
locals {
  code_table_creator = huaweicloud_dataarts_architecture_code_table.test.created_by
}

data "huaweicloud_dataarts_architecture_code_tables" "filter_by_create_by" {
  depends_on = [huaweicloud_dataarts_architecture_code_table.test]

  workspace_id = "%[2]s"
  create_by    = local.code_table_creator
}

locals {
  create_by_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_code_tables.filter_by_create_by.code_tables :
      v.created_by == local.code_table_creator
  ]
}

output "is_create_by_filter_useful" {
  value = length(local.create_by_filter_result) > 0 && alltrue(local.create_by_filter_result)
}

# Filter by directory ID.
locals {
  code_table_directory_id = huaweicloud_dataarts_architecture_code_table.test.directory_id
}

data "huaweicloud_dataarts_architecture_code_tables" "filter_by_directory_id" {
  depends_on = [huaweicloud_dataarts_architecture_code_table.test]

  workspace_id = "%[2]s"
  directory_id = local.code_table_directory_id
}

locals {
  directory_id_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_code_tables.filter_by_directory_id.code_tables :
      v.directory_id == local.code_table_directory_id
  ]
}

output "is_directory_id_filter_useful" {
  value = length(local.directory_id_filter_result) > 0 && alltrue(local.directory_id_filter_result)
}

# Filter by status.
locals {
  code_table_status = huaweicloud_dataarts_architecture_code_table.test.status
}

data "huaweicloud_dataarts_architecture_code_tables" "filter_by_status" {
  depends_on = [huaweicloud_dataarts_architecture_code_table.test]

  workspace_id = "%[2]s"
  status       = local.code_table_status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_code_tables.filter_by_status.code_tables :
      v.status == local.code_table_status
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

# Filter by begin time.
locals {
  code_table_begin_time = timeadd(huaweicloud_dataarts_architecture_directory.test.created_at, "-10m")
}

data "huaweicloud_dataarts_architecture_code_tables" "filter_by_begin_time" {
  depends_on = [huaweicloud_dataarts_architecture_code_table.test]

  workspace_id = "%[2]s"
  begin_time   = local.code_table_begin_time
}

output "is_begin_time_filter_useful" {
  value = contains(data.huaweicloud_dataarts_architecture_code_tables.filter_by_begin_time.code_tables[*].id,
    huaweicloud_dataarts_architecture_code_table.test.id)
}

# Filter by end time.
locals {
  code_table_end_time = timeadd(huaweicloud_dataarts_architecture_directory.test.created_at, "10m")
}

data "huaweicloud_dataarts_architecture_code_tables" "filter_by_end_time" {
  depends_on = [huaweicloud_dataarts_architecture_code_table.test]

  workspace_id = "%[2]s"
  end_time     = local.code_table_end_time
}

output "is_end_time_filter_useful" {
  value = contains(data.huaweicloud_dataarts_architecture_code_tables.filter_by_end_time.code_tables[*].id,
    huaweicloud_dataarts_architecture_code_table.test.id)
}
`, testAccDataArchitectureCodeTables_basic_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}
