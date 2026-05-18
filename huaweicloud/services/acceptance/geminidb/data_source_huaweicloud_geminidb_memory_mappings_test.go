package geminidb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceMemoryMappings_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_geminidb_memory_mappings.test"
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
				Config: testAccDataSourceMemoryMappings_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "dbcache_mappings.#"),
					resource.TestCheckResourceAttrSet(dataSource, "dbcache_mappings.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "dbcache_mappings.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "dbcache_mappings.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "dbcache_mappings.0.source_instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "dbcache_mappings.0.source_instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "dbcache_mappings.0.target_instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "dbcache_mappings.0.target_instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "dbcache_mappings.0.created"),
					resource.TestCheckResourceAttrSet(dataSource, "dbcache_mappings.0.updated"),

					resource.TestCheckOutput("mapping_id_filter_useful", "true"),
					resource.TestCheckOutput("name_filter_useful", "true"),
					resource.TestCheckOutput("source_instance_id_filter_useful", "true"),
					resource.TestCheckOutput("source_instance_name_filter_useful", "true"),
					resource.TestCheckOutput("target_instance_id_filter_useful", "true"),
					resource.TestCheckOutput("target_instance_name_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceMemoryMappings_basic = `
data "huaweicloud_geminidb_memory_mappings" "test" {}

locals {
  mapping_id           = data.huaweicloud_geminidb_memory_mappings.test.dbcache_mappings[0].id
  name                 = data.huaweicloud_geminidb_memory_mappings.test.dbcache_mappings[0].name
  source_instance_id   = data.huaweicloud_geminidb_memory_mappings.test.dbcache_mappings[0].source_instance_id
  source_instance_name = data.huaweicloud_geminidb_memory_mappings.test.dbcache_mappings[0].source_instance_name
  target_instance_id   = data.huaweicloud_geminidb_memory_mappings.test.dbcache_mappings[0].target_instance_id
  target_instance_name = data.huaweicloud_geminidb_memory_mappings.test.dbcache_mappings[0].target_instance_name
}

data "huaweicloud_geminidb_memory_mappings" "mapping_id_filter" {
  mapping_id = local.mapping_id
}

output "mapping_id_filter_useful" {
  value = length(data.huaweicloud_geminidb_memory_mappings.mapping_id_filter.dbcache_mappings) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_memory_mappings.mapping_id_filter.dbcache_mappings[*].id : v == local.mapping_id]
  )
}

data "huaweicloud_geminidb_memory_mappings" "name_filter" {	
  name = local.name
}

output "name_filter_useful" {
  value = length(data.huaweicloud_geminidb_memory_mappings.name_filter.dbcache_mappings) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_memory_mappings.name_filter.dbcache_mappings[*].name : v == local.name]
  )
}

data "huaweicloud_geminidb_memory_mappings" "source_instance_id_filter" {
  source_instance_id = local.source_instance_id
}

output "source_instance_id_filter_useful" {
  value = length(data.huaweicloud_geminidb_memory_mappings.source_instance_id_filter.dbcache_mappings) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_memory_mappings.source_instance_id_filter.dbcache_mappings[*].source_instance_id : v
    == local.source_instance_id]
  )
}

data "huaweicloud_geminidb_memory_mappings" "source_instance_name_filter" {
  source_instance_name = local.source_instance_name
}

output "source_instance_name_filter_useful" {
  value = length(data.huaweicloud_geminidb_memory_mappings.source_instance_name_filter.dbcache_mappings) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_memory_mappings.source_instance_name_filter.dbcache_mappings[*].source_instance_name : v
    == local.source_instance_name]
  )
}

data "huaweicloud_geminidb_memory_mappings" "target_instance_id_filter" {
  target_instance_id = local.target_instance_id
}

output "target_instance_id_filter_useful" {
  value = length(data.huaweicloud_geminidb_memory_mappings.target_instance_id_filter.dbcache_mappings) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_memory_mappings.target_instance_id_filter.dbcache_mappings[*].target_instance_id : v
    == local.target_instance_id]
  )
}

data "huaweicloud_geminidb_memory_mappings" "target_instance_name_filter" {
  target_instance_name = local.target_instance_name
}

output "target_instance_name_filter_useful" {
  value = length(data.huaweicloud_geminidb_memory_mappings.target_instance_name_filter.dbcache_mappings) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_memory_mappings.target_instance_name_filter.dbcache_mappings[*].target_instance_name : v
    == local.target_instance_name]
  )
}
`
