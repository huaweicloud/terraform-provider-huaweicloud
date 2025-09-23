package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceArchitectureTableModels_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dataarts_architecture_table_models.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		name       = acceptance.RandomAccResourceName()

		byName   = "data.huaweicloud_dataarts_architecture_table_models.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byEnName   = "data.huaweicloud_dataarts_architecture_table_models.filter_by_en_name"
		dcByEnName = acceptance.InitDataSourceCheck(byEnName)

		byNotName   = "data.huaweicloud_dataarts_architecture_table_models.filter_by_not_found_name"
		dcByNotName = acceptance.InitDataSourceCheck(byNotName)

		bySubjectId   = "data.huaweicloud_dataarts_architecture_table_models.test"
		dcBySubjectId = acceptance.InitDataSourceCheck(bySubjectId)

		byStatus   = "data.huaweicloud_dataarts_architecture_table_models.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byCreator   = "data.huaweicloud_dataarts_architecture_table_models.filter_by_creator"
		dcByCreator = acceptance.InitDataSourceCheck(byCreator)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsSubjectID(t)
			acceptance.TestAccPreCheckDataArtsConnectionID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDatasourceArchitectureTableModels_workspace_id_not_found,
				ExpectError: regexp.MustCompile("DLG.0818"),
			},
			{
				Config:      testAccDatasourceArchitectureTableModels_model_id_not_found(),
				ExpectError: regexp.MustCompile("DLG.3908"),
			},
			{
				Config: testAccDatasourceArchitectureTableModels_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tables.0.model_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tables.0.table_name"),
					resource.TestCheckResourceAttrSet(dataSource, "tables.0.physical_table_name"),
					resource.TestCheckResourceAttrSet(dataSource, "tables.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "tables.0.dw_type"),
					resource.TestCheckResourceAttrSet(dataSource, "tables.0.subject_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tables.0.catalog_path"),
					resource.TestCheckResourceAttrSet(dataSource, "tables.0.attributes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tables.0.attributes.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "tables.0.attributes.0.name_en"),
					resource.TestCheckResourceAttrSet(dataSource, "tables.0.attributes.0.data_type"),
					resource.TestMatchResourceAttr(dataSource, "tables.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(dataSource, "tables.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(dataSource, "tables.0.created_by"),
					resource.TestCheckResourceAttrSet(dataSource, "tables.0.updated_by"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByEnName.CheckResourceExists(),
					resource.TestCheckOutput("is_en_name_filter_useful", "true"),
					resource.TestCheckResourceAttr(byEnName, "tables.0.dw_id", acceptance.HW_DATAARTS_CONNECTION_ID),
					resource.TestCheckResourceAttrSet(byEnName, "tables.0.db_name"),
					resource.TestCheckResourceAttrSet(byEnName, "tables.0.queue_name"),
					resource.TestCheckResourceAttrSet(byEnName, "tables.0.table_type"),
					resource.TestCheckResourceAttrSet(byEnName, "tables.0.data_format"),
					resource.TestCheckResourceAttrSet(byEnName, "tables.0.configs"),
					resource.TestCheckResourceAttrSet(byEnName, "tables.0.obs_location"),
					resource.TestCheckResourceAttrSet(byEnName, "tables.0.distribute"),
					resource.TestCheckResourceAttrSet(byEnName, "tables.0.distribute_column"),
					resource.TestCheckResourceAttrSet(byEnName, "tables.0.dirty_out_database"),
					resource.TestCheckResourceAttr(byEnName, "tables.0.dirty_out_switch", "true"),
					resource.TestCheckResourceAttrSet(byEnName, "tables.0.dirty_out_prefix"),
					resource.TestCheckResourceAttrSet(byEnName, "tables.0.dirty_out_suffix"),
					resource.TestCheckResourceAttrPair(byEnName, "tables.0.related_logic_table_id",
						"huaweicloud_dataarts_architecture_table_model.logical_table_model", "id"),
					resource.TestCheckResourceAttrSet(byEnName, "tables.0.attributes.0.data_type_extend"),
					resource.TestCheckResourceAttrSet(byEnName, "tables.0.attributes.0.description"),
					resource.TestCheckResourceAttr(byEnName, "tables.0.attributes.0.is_primary_key", "true"),
					resource.TestCheckResourceAttr(byEnName, "tables.0.attributes.0.not_null", "true"),
					resource.TestCheckResourceAttr(byEnName, "tables.0.attributes.0.is_partition_key", "false"),
					resource.TestMatchResourceAttr(dataSource, "tables.0.attributes.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(dataSource, "tables.0.attributes.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					dcByNotName.CheckResourceExists(),
					resource.TestCheckOutput("not_found_name", "true"),
					dcBySubjectId.CheckResourceExists(),
					resource.TestCheckOutput("is_subject_id_filter_useful", "true"),
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					dcByCreator.CheckResourceExists(),
					resource.TestCheckOutput("is_creator_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDatasourceArchitectureTableModels_workspace_id_not_found = `
data "huaweicloud_dataarts_architecture_table_models" "test" {
  workspace_id = "not_found"
  model_id     = "1268862656148799488"
}
`

func testAccDatasourceArchitectureTableModels_model_id_not_found() string {
	return fmt.Sprintf(`
data "huaweicloud_dataarts_architecture_table_models" "test" {
  workspace_id = "%s"
  model_id     = "not_found"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID)
}

func testAccDatasourceArchitectureTableModels_basic(name string) string {
	return fmt.Sprintf(`
# Creating a logical model
resource "huaweicloud_dataarts_architecture_model" "logical_model" {
  workspace_id = "%[1]s"
  name         = "%[2]s_logical"
  type         = "THIRD_NF"
  physical     = false
}

resource "huaweicloud_dataarts_architecture_table_model" "logical_table_model" {
  workspace_id        = "%[1]s"
  model_id            = huaweicloud_dataarts_architecture_model.logical_model.id
  subject_id          = "%[3]s"
  table_name          = "%[2]s_logical"
  physical_table_name = "%[2]s_logical_en"
  description         = "Logical table model"
  dw_type             = "UNSPECIFIED"

  attributes {
    name      = "logical_key"
    name_en   = "logical_key_en"
    data_type = "STRING"
    ordinal   = "1"
  }
}

# Creating a physical model
resource "huaweicloud_dataarts_architecture_model" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  type         = "THIRD_NF"
  physical     = true
  dw_type      = "DLI"
}

resource "huaweicloud_dataarts_architecture_table_model" "test" {
  workspace_id                 = "%[1]s"
  model_id                     = huaweicloud_dataarts_architecture_model.test.id
  subject_id                   = "%[3]s"
  table_name                   = "%[2]s"
  physical_table_name          = "%[2]s_en_%[2]s"
  description                  = "Created DLI table model by terraform scrit"
  dw_type                      = huaweicloud_dataarts_architecture_model.test.dw_type
  dw_id                        = "%[4]s"
  db_name                      = "default"
  queue_name                   = "default"
  table_type                   = "EXTERNAL"
  data_format                  = "Parquet"
  related_logic_table_model_id = huaweicloud_dataarts_architecture_table_model.logical_table_model.id
  related_logic_model_id       = huaweicloud_dataarts_architecture_model.logical_model.id

  configs = jsonencode({
    "advanced" : "terraform"
  })

  obs_location      = "%[2]s"
  distribute        = "HASH"
  distribute_column = "key"

  attributes {
    name             = "key"
    name_en          = "key_en"
    data_type        = "STRING"
    ordinal          = "1"
    data_type_extend = "extend"
    description      = "the attr of physical model table"
    is_primary_key   = true
    not_null         = true
    is_partition_key = false
  }

  dirty_out_database = "%[2]s"
  dirty_out_switch   = true
  dirty_out_prefix   = "prefix_test"
  dirty_out_suffix   = "suffix_test"
}

data "huaweicloud_dataarts_architecture_table_models" "test" {
  depends_on = [
    huaweicloud_dataarts_architecture_table_model.test
  ]

  workspace_id = "%[1]s"
  model_id     = huaweicloud_dataarts_architecture_model.test.id
}

# Filter by Chinese name
locals {
  name = huaweicloud_dataarts_architecture_table_model.test.table_name
}

data "huaweicloud_dataarts_architecture_table_models" "filter_by_name" {
  depends_on = [
    huaweicloud_dataarts_architecture_table_model.test
  ]

  workspace_id = "%[1]s"
  model_id     = huaweicloud_dataarts_architecture_model.test.id
  name         = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_table_models.filter_by_name.tables[*].table_name: v == local.name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by English name
locals {
  en_name = huaweicloud_dataarts_architecture_table_model.test.physical_table_name
}

data "huaweicloud_dataarts_architecture_table_models" "filter_by_en_name" {
  depends_on = [
    huaweicloud_dataarts_architecture_table_model.test
  ]

  workspace_id = "%[1]s"
  model_id     = huaweicloud_dataarts_architecture_model.test.id
  name         = local.en_name
}

locals {
  en_name_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_table_models.filter_by_en_name.tables[*].physical_table_name: v == local.en_name
  ]
}

output "is_en_name_filter_useful" {
  value = length(local.en_name_filter_result) > 0 && alltrue(local.en_name_filter_result)
}

# Not found name
data "huaweicloud_dataarts_architecture_table_models" "filter_by_not_found_name" {
  depends_on = [
    huaweicloud_dataarts_architecture_table_model.test
  ]

  workspace_id = "%[1]s"
  model_id     = huaweicloud_dataarts_architecture_model.test.id
  name         = "not_found"
}

output "not_found_name" {
  value = length(data.huaweicloud_dataarts_architecture_table_models.filter_by_not_found_name.tables) == 0
}

# Filter by subject ID
data "huaweicloud_dataarts_architecture_table_models" "filter_by_subject_id" {
  depends_on = [
    huaweicloud_dataarts_architecture_table_model.test
  ]

  workspace_id = "%[1]s"
  model_id     = huaweicloud_dataarts_architecture_model.test.id
  subject_id   = "%[3]s"
}

locals {
  subject_id_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_table_models.filter_by_subject_id.tables[*].subject_id: v == "%[3]s"
  ]
}

output "is_subject_id_filter_useful" {
  value = length(local.subject_id_filter_result) > 0 && alltrue(local.subject_id_filter_result)
}

# Filter by status
locals {
  status = huaweicloud_dataarts_architecture_table_model.test.status
}

data "huaweicloud_dataarts_architecture_table_models" "filter_by_status" {
  workspace_id = "%[1]s"
  model_id     = huaweicloud_dataarts_architecture_model.test.id
  status       = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_table_models.filter_by_status.tables[*].status: v == local.status
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

# Filter by creator
locals {
  created_by = huaweicloud_dataarts_architecture_table_model.test.created_by
}

data "huaweicloud_dataarts_architecture_table_models" "filter_by_creator" {
  workspace_id = "%[1]s"
  model_id     = huaweicloud_dataarts_architecture_model.test.id
  created_by   = local.created_by
}

locals {
  creator_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_table_models.filter_by_creator.tables[*].created_by: v == local.created_by
  ]
}

output "is_creator_filter_useful" {
  value = length(local.creator_filter_result) > 0 && alltrue(local.creator_filter_result)
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name, acceptance.HW_DATAARTS_SUBJECT_ID, acceptance.HW_DATAARTS_CONNECTION_ID)
}
