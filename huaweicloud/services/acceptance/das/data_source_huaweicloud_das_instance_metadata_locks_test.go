package das

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccInstanceMetadataLocks_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_das_instance_metadata_locks.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byDatabase   = "data.huaweicloud_das_instance_metadata_locks.filter_by_database"
		dcByDatabase = acceptance.InitDataSourceCheck(byDatabase)

		byThreadId   = "data.huaweicloud_das_instance_metadata_locks.filter_by_thread_id"
		dcByThreadId = acceptance.InitDataSourceCheck(byThreadId)

		byTable   = "data.huaweicloud_das_instance_metadata_locks.filter_by_table"
		dcByTable = acceptance.InitDataSourceCheck(byTable)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDasInstanceIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceMetadataLocks_basic(),
				Check: resource.ComposeTestCheckFunc(
					// Query all metadata locks
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "metadata_locks.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "metadata_locks.0.thread_id"),
					resource.TestCheckResourceAttrSet(all, "metadata_locks.0.lock_status"),
					resource.TestCheckResourceAttrSet(all, "metadata_locks.0.lock_mode"),
					resource.TestCheckResourceAttrSet(all, "metadata_locks.0.lock_type"),
					resource.TestCheckResourceAttrSet(all, "metadata_locks.0.lock_duration"),
					resource.TestCheckResourceAttrSet(all, "metadata_locks.0.table_schema"),
					resource.TestCheckResourceAttrSet(all, "metadata_locks.0.table_name"),
					resource.TestCheckResourceAttrSet(all, "metadata_locks.0.user"),
					resource.TestCheckResourceAttrSet(all, "metadata_locks.0.time"),
					resource.TestCheckResourceAttrSet(all, "metadata_locks.0.host"),
					resource.TestCheckResourceAttrSet(all, "metadata_locks.0.database"),
					resource.TestCheckResourceAttrSet(all, "metadata_locks.0.command"),
					resource.TestCheckResourceAttrSet(all, "metadata_locks.0.trx_exec_time"),

					// Filter by database
					dcByDatabase.CheckResourceExists(),
					resource.TestCheckOutput("is_database_filter_useful", "true"),

					// Filter by thread_id
					dcByThreadId.CheckResourceExists(),
					resource.TestCheckOutput("is_thread_id_filter_useful", "true"),

					// Filter by table
					dcByTable.CheckResourceExists(),
					resource.TestCheckOutput("is_table_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccInstanceMetadataLocks_base() string {
	return fmt.Sprintf(`
locals {
  instance_ids = split(",", "%[1]s")
}

data "huaweicloud_das_database_users" "all" {
  instance_id = local.instance_ids[0]
}

locals {
  db_user_id = try(data.huaweicloud_das_database_users.all.users[0].id, "")
}

`, acceptance.HW_DAS_INSTANCE_IDS)
}

func testAccInstanceMetadataLocks_basic() string {
	return fmt.Sprintf(`
%[1]s

// Query all metadata locks
data "huaweicloud_das_instance_metadata_locks" "all" {
  instance_id = local.instance_ids[0]
  db_user_id  = local.db_user_id
}

locals {
  all_databases = [for v in data.huaweicloud_das_instance_metadata_locks.all.metadata_locks : v.database if v.database != ""]
  all_thread_ids = [for v in data.huaweicloud_das_instance_metadata_locks.all.metadata_locks : v.thread_id if v.thread_id != ""]
  all_tables     = [for v in data.huaweicloud_das_instance_metadata_locks.all.metadata_locks : v.table_name if v.table_name != ""]
}

// Filter by database
data "huaweicloud_das_instance_metadata_locks" "filter_by_database" {
  instance_id = local.instance_ids[0]
  db_user_id  = local.db_user_id
  database    = length(local.all_databases) > 0 ? local.all_databases[0] : "non_exist"
}

locals {
  database_filter_result = [
    for v in data.huaweicloud_das_instance_metadata_locks.filter_by_database.metadata_locks :
    v.database == (length(local.all_databases) > 0 ? local.all_databases[0] : "non_exist")
  ]
}

output "is_database_filter_useful" {
  value = length(local.database_filter_result) == 0 || alltrue(local.database_filter_result)
}

// Filter by thread_id
data "huaweicloud_das_instance_metadata_locks" "filter_by_thread_id" {
  instance_id = local.instance_ids[0]
  db_user_id  = local.db_user_id
  thread_id   = length(local.all_thread_ids) > 0 ? local.all_thread_ids[0] : "non_exist"
}

locals {
  thread_id_filter_result = [
    for v in data.huaweicloud_das_instance_metadata_locks.filter_by_thread_id.metadata_locks :
    v.thread_id == (length(local.all_thread_ids) > 0 ? local.all_thread_ids[0] : "non_exist")
  ]
}

output "is_thread_id_filter_useful" {
  value = length(local.thread_id_filter_result) == 0 || alltrue(local.thread_id_filter_result)
}

// Filter by table
data "huaweicloud_das_instance_metadata_locks" "filter_by_table" {
  instance_id = local.instance_ids[0]
  db_user_id  = local.db_user_id
  table       = length(local.all_tables) > 0 ? local.all_tables[0] : "non_exist"
}

locals {
  table_filter_result = [
    for v in data.huaweicloud_das_instance_metadata_locks.filter_by_table.metadata_locks :
    v.table_name == (length(local.all_tables) > 0 ? local.all_tables[0] : "non_exist")
  ]
}

output "is_table_filter_useful" {
  value = length(local.table_filter_result) == 0 || alltrue(local.table_filter_result)
}
`, testAccInstanceMetadataLocks_base())
}
