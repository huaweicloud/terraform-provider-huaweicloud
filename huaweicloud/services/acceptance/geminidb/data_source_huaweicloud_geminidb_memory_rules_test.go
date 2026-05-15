package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceMemoryRules_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_geminidb_memory_rules.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGeminiDBMappingId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMemoryRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "rules.#"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.source_db_schema"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.source_db_table"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.storage_type"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.target_database"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.key_columns.#"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.value_columns.#"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.ttl"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.key_separator"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.key_prefix"),

					resource.TestCheckOutput("rule_id_filter_useful", "true"),
					resource.TestCheckOutput("rule_name_filter_useful", "true"),
					resource.TestCheckOutput("source_db_schema_filter_useful", "true"),
					resource.TestCheckOutput("source_db_table_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceMemoryRules_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_geminidb_memory_rules" "test" {
  dbcache_mapping_id = "%[1]s"
}

locals {
  rule_id          = data.huaweicloud_geminidb_memory_rules.test.rules[0].id
  rule_name        = data.huaweicloud_geminidb_memory_rules.test.rules[0].name
  source_db_schema = data.huaweicloud_geminidb_memory_rules.test.rules[0].source_db_schema
  source_db_table  = data.huaweicloud_geminidb_memory_rules.test.rules[0].source_db_table
}

data "huaweicloud_geminidb_memory_rules" "rule_id_filter" {
  dbcache_mapping_id = "%[1]s"
  rule_id            = local.rule_id
}

output "rule_id_filter_useful" {
  value = length(data.huaweicloud_geminidb_memory_rules.rule_id_filter.rules) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_memory_rules.rule_id_filter.rules[*].id : v == local.rule_id]
  )
}

data "huaweicloud_geminidb_memory_rules" "rule_name_filter" {	
  dbcache_mapping_id = "%[1]s"
  rule_name          = local.rule_name
}

output "rule_name_filter_useful" {
  value = length(data.huaweicloud_geminidb_memory_rules.rule_name_filter.rules) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_memory_rules.rule_name_filter.rules[*].name : v == local.rule_name]
  )
}

data "huaweicloud_geminidb_memory_rules" "source_db_schema_filter" {
  dbcache_mapping_id = "%[1]s"
  source_db_schema   = local.source_db_schema
}

output "source_db_schema_filter_useful" {
  value = length(data.huaweicloud_geminidb_memory_rules.source_db_schema_filter.rules) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_memory_rules.source_db_schema_filter.rules[*].source_db_schema : v == local.source_db_schema]
  )
}

data "huaweicloud_geminidb_memory_rules" "source_db_table_filter" {
  dbcache_mapping_id = "%[1]s"
  source_db_table    = local.source_db_table
}

output "source_db_table_filter_useful" {
  value = length(data.huaweicloud_geminidb_memory_rules.source_db_table_filter.rules) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_memory_rules.source_db_table_filter.rules[*].source_db_table : v == local.source_db_table]
  )
}
`, acceptance.HW_GEMINIDB_MAPPING_ID)
}
