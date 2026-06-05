package drs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBackupMigrationJobs_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_backup_migration_jobs.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBackupMigrationJobs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "jobs.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "jobs.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "jobs.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "jobs.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "jobs.0.engine_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "jobs.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "jobs.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "jobs.0.finish_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "jobs.0.enterprise_project_id"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_description_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceBackupMigrationJobs_basic() string {
	return `
data "huaweicloud_drs_backup_migration_jobs" "test" {
  sort_key              = "name"
  sort_dir              = "desc"
  enterprise_project_id = "0"
}

# Filter using name.
locals {
  name = data.huaweicloud_drs_backup_migration_jobs.test.jobs[0].name
}

data "huaweicloud_drs_backup_migration_jobs" "name_filter" {
  name                  = local.name
  enterprise_project_id = "0"
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_drs_backup_migration_jobs.name_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_drs_backup_migration_jobs.name_filter.jobs[*].name : v == local.name]
  )
}

# Filter using status.
locals {
  status = data.huaweicloud_drs_backup_migration_jobs.test.jobs[0].status
}

data "huaweicloud_drs_backup_migration_jobs" "status_filter" {
  status                = local.status
  enterprise_project_id = "0"
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_drs_backup_migration_jobs.status_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_drs_backup_migration_jobs.status_filter.jobs[*].status : v == local.status]
  )
}

# Filter using description.
locals {
  description = data.huaweicloud_drs_backup_migration_jobs.test.jobs[0].description
}

data "huaweicloud_drs_backup_migration_jobs" "description_filter" {
  description           = local.description
  enterprise_project_id = "0"
}

output "is_description_filter_useful" {
  value = length(data.huaweicloud_drs_backup_migration_jobs.description_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_drs_backup_migration_jobs.description_filter.jobs[*].description : v == local.description]
  )
}
`
}
