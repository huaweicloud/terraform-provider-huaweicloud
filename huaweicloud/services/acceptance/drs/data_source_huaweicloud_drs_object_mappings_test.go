package drs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsObjectMappings_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_drs_object_mappings.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byDbName   = "data.huaweicloud_drs_object_mappings.filter_by_db_name"
		dcByDbName = acceptance.InitDataSourceCheck(byDbName)

		bySchemaName   = "data.huaweicloud_drs_object_mappings.filter_by_schema_name"
		dcBySchemaName = acceptance.InitDataSourceCheck(bySchemaName)

		byTableName   = "data.huaweicloud_drs_object_mappings.filter_by_table_name"
		dcByTableName = acceptance.InitDataSourceCheck(byTableName)

		byHasColumnInfo   = "data.huaweicloud_drs_object_mappings.filter_by_has_column_info"
		dcByHasColumnInfo = acceptance.InitDataSourceCheck(byHasColumnInfo)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDrsObjectMappings_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "object_mapping_list.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestCheckResourceAttrSet(all, "object_mapping_list.0.source_db_name"),
					resource.TestCheckResourceAttrSet(all, "object_mapping_list.0.source_table_name"),
					resource.TestCheckResourceAttrSet(all, "object_mapping_list.0.target_db_name"),
					resource.TestCheckResourceAttrSet(all, "object_mapping_list.0.target_table_name"),
					resource.TestCheckResourceAttrSet(all, "object_mapping_list.0.has_column_info"),
					dcByDbName.CheckResourceExists(),
					resource.TestCheckOutput("is_db_name_filter_useful", "true"),
					dcBySchemaName.CheckResourceExists(),
					resource.TestCheckOutput("is_schema_name_filter_useful", "true"),
					dcByTableName.CheckResourceExists(),
					resource.TestCheckOutput("is_table_name_filter_useful", "true"),
					dcByHasColumnInfo.CheckResourceExists(),
					resource.TestCheckOutput("is_has_column_info_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceDrsObjectMappings_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_object_mappings" "all" {
  job_id = "%[1]s"
}

# Filter subjects by source_db_name.
locals {
  db_name = data.huaweicloud_drs_object_mappings.all.object_mapping_list[0].source_db_name
}

data "huaweicloud_drs_object_mappings" "filter_by_db_name" {
  job_id  = "%[1]s"
  db_name = local.db_name
}

locals {
  db_name_filter_result = [
    for v in data.huaweicloud_drs_object_mappings.filter_by_db_name.object_mapping_list
    : v.source_db_name == local.db_name
  ]
}

output "is_db_name_filter_useful" {
  value = length(local.db_name_filter_result) > 0 && alltrue(local.db_name_filter_result)
}

# Filter subjects by source_schema_name.
locals {
  schema_name = data.huaweicloud_drs_object_mappings.all.object_mapping_list[0].source_schema_name
}

data "huaweicloud_drs_object_mappings" "filter_by_schema_name" {
  job_id      = "%[1]s"
  schema_name = local.schema_name
}

locals {
	schema_name_filter_result = [
	  for v in data.huaweicloud_drs_object_mappings.filter_by_schema_name.object_mapping_list :
	  v.source_schema_name == local.schema_name
  ]
}
  
output "is_schema_name_filter_useful" {
	value = length(local.schema_name_filter_result) > 0 && alltrue(local.schema_name_filter_result)
}

# Filter subjects by source_table_name.
locals {
  table_name = data.huaweicloud_drs_object_mappings.all.object_mapping_list[0].source_table_name
}

data "huaweicloud_drs_object_mappings" "filter_by_table_name" {
  job_id     = "%[1]s"
  table_name = local.table_name
}

locals {
  table_name_filter_result = [
    for v in data.huaweicloud_drs_object_mappings.filter_by_table_name.object_mapping_list
    : v.source_table_name == local.table_name
  ]
}

output "is_table_name_filter_useful" {
  value = length(local.table_name_filter_result) > 0 && alltrue(local.table_name_filter_result)
}

# Filter subjects by has_column_info.
locals {
	has_column_info = data.huaweicloud_drs_object_mappings.all.object_mapping_list[0].has_column_info
}

data "huaweicloud_drs_object_mappings" "filter_by_has_column_info" {
  job_id          = "%[1]s"
  has_column_info = tostring(local.has_column_info)
}

locals {
	has_column_info_filter_result = [
	  for v in data.huaweicloud_drs_object_mappings.filter_by_has_column_info.object_mapping_list :
	  v.has_column_info == local.has_column_info
  ]
}

output "is_has_column_info_filter_useful" {
  value = length(local.has_column_info_filter_result) > 0 && alltrue(local.has_column_info_filter_result)
}
`, acceptance.HW_DRS_JOB_ID)
}
