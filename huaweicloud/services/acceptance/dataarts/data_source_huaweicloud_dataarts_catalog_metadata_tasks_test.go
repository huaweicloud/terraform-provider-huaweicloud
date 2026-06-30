package dataarts

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataCatalogMetadataTasks_basic(t *testing.T) {
	var (
		name        = acceptance.RandomAccResourceName()
		currentTime = time.Now().Local().Format(time.RFC3339)

		all = "data.huaweicloud_dataarts_catalog_metadata_tasks.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByUserName   = "data.huaweicloud_dataarts_catalog_metadata_tasks.filter_by_user_name"
		dcFilterByUserName = acceptance.InitDataSourceCheck(filterByUserName)

		filterByName   = "data.huaweicloud_dataarts_catalog_metadata_tasks.filter_by_name"
		dcFilterByName = acceptance.InitDataSourceCheck(filterByName)

		filterByDataSourceType   = "data.huaweicloud_dataarts_catalog_metadata_tasks.filter_by_data_source_type"
		dcFilterByDataSourceType = acceptance.InitDataSourceCheck(filterByDataSourceType)

		filterByDataConnectionId   = "data.huaweicloud_dataarts_catalog_metadata_tasks.filter_by_data_connection_id"
		dcFilterByDataConnectionId = acceptance.InitDataSourceCheck(filterByDataConnectionId)

		filterByStartTime   = "data.huaweicloud_dataarts_catalog_metadata_tasks.filter_by_start_time"
		dcFilterByStartTime = acceptance.InitDataSourceCheck(filterByStartTime)

		filterByEndTime   = "data.huaweicloud_dataarts_catalog_metadata_tasks.filter_by_end_time"
		dcFilterByEndTime = acceptance.InitDataSourceCheck(filterByEndTime)

		filterByDirectoryId   = "data.huaweicloud_dataarts_catalog_metadata_tasks.filter_by_directory_id"
		dcFilterByDirectoryId = acceptance.InitDataSourceCheck(filterByDirectoryId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsConnectionName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.12.1",
			},
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccDataCatalogMetadataTasks_nonExistentWorkspace(),
				ExpectError: regexp.MustCompile("Processing failed."),
			},
			{
				Config: testAccDataSourceCatalogMetadataTasks_basic_step1(name, currentTime),
			},
			{
				Config: testAccDataSourceCatalogMetadataTasks_basic_step2(name, currentTime),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "metadata_tasks.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "metadata_tasks.0.id"),
					resource.TestCheckResourceAttrPair(all, "metadata_tasks.0.name",
						"huaweicloud_dataarts_catalog_metadata_task.test", "name"),
					resource.TestCheckResourceAttrPair(all, "metadata_tasks.0.dir_id",
						"huaweicloud_dataarts_catalog_metadata_task.test", "dir_id"),
					resource.TestCheckResourceAttrPair(all, "metadata_tasks.0.schedule_config.0.cron_expression",
						"huaweicloud_dataarts_catalog_metadata_task.test", "schedule_config.0.cron_expression"),
					resource.TestCheckResourceAttrPair(all, "metadata_tasks.0.schedule_config.0.max_time_out",
						"huaweicloud_dataarts_catalog_metadata_task.test", "schedule_config.0.max_time_out"),
					resource.TestCheckResourceAttrPair(all, "metadata_tasks.0.schedule_config.0.end_time",
						"huaweicloud_dataarts_catalog_metadata_task.test", "schedule_config.0.end_time"),
					resource.TestCheckResourceAttrPair(all, "metadata_tasks.0.schedule_config.0.interval",
						"huaweicloud_dataarts_catalog_metadata_task.test", "schedule_config.0.interval"),
					resource.TestCheckResourceAttrPair(all, "metadata_tasks.0.schedule_config.0.schedule_type",
						"huaweicloud_dataarts_catalog_metadata_task.test", "schedule_config.0.schedule_type"),
					resource.TestCheckResourceAttrPair(all, "metadata_tasks.0.schedule_config.0.start_time",
						"huaweicloud_dataarts_catalog_metadata_task.test", "schedule_config.0.start_time"),
					resource.TestCheckResourceAttrPair(all, "metadata_tasks.0.schedule_config.0.job_id",
						"huaweicloud_dataarts_catalog_metadata_task.test", "schedule_config.0.job_id"),
					resource.TestCheckResourceAttrPair(all, "metadata_tasks.0.schedule_config.0.enabled",
						"huaweicloud_dataarts_catalog_metadata_task.test", "schedule_config.0.enabled"),
					resource.TestCheckResourceAttrPair(all, "metadata_tasks.0.data_source_type",
						"huaweicloud_dataarts_catalog_metadata_task.test", "data_source_type"),
					resource.TestCheckResourceAttrPair(all, "metadata_tasks.0.task_config",
						"huaweicloud_dataarts_catalog_metadata_task.test", "task_config"),
					resource.TestCheckResourceAttrPair(all, "metadata_tasks.0.description",
						"huaweicloud_dataarts_catalog_metadata_task.test", "description"),
					resource.TestCheckResourceAttrPair(all, "metadata_tasks.0.create_time",
						"huaweicloud_dataarts_catalog_metadata_task.test", "create_time"),
					resource.TestCheckResourceAttrPair(all, "metadata_tasks.0.update_time",
						"huaweicloud_dataarts_catalog_metadata_task.test", "update_time"),
					resource.TestCheckResourceAttrPair(all, "metadata_tasks.0.user_id",
						"huaweicloud_dataarts_catalog_metadata_task.test", "user_id"),
					resource.TestCheckResourceAttrPair(all, "metadata_tasks.0.user_name",
						"huaweicloud_dataarts_catalog_metadata_task.test", "user_name"),
					resource.TestCheckResourceAttrPair(all, "metadata_tasks.0.path",
						"huaweicloud_dataarts_catalog_metadata_task.test", "path"),
					resource.TestCheckResourceAttrPair(all, "metadata_tasks.0.duty_person",
						"huaweicloud_dataarts_catalog_metadata_task.test", "duty_person"),

					// filter by user name
					dcFilterByUserName.CheckResourceExists(),
					resource.TestCheckOutput("is_user_name_filter_useful", "true"),

					// filter by name
					dcFilterByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),

					// filter by data source type
					dcFilterByDataSourceType.CheckResourceExists(),
					resource.TestCheckOutput("is_data_source_type_filter_useful", "true"),

					// filter by data connection ID
					dcFilterByDataConnectionId.CheckResourceExists(),
					resource.TestCheckOutput("is_data_connection_id_filter_useful", "true"),

					// filter by start time
					dcFilterByStartTime.CheckResourceExists(),
					resource.TestCheckOutput("is_start_time_filter_useful", "true"),

					// filter by end time
					dcFilterByEndTime.CheckResourceExists(),
					resource.TestCheckOutput("is_end_time_filter_useful", "true"),

					// filter by directory ID
					dcFilterByDirectoryId.CheckResourceExists(),
					resource.TestCheckOutput("is_directory_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataCatalogMetadataTasks_nonExistentWorkspace() string {
	randUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
data "huaweicloud_dataarts_catalog_metadata_tasks" "test" {
  workspace_id = "%[1]s"
}
`, randUUID.String())
}

func testAccDataSourceCatalogMetadataTasks_basic_step1(name, currentTime string) string {
	return fmt.Sprintf(` 
variable "data_connection_id" {
  type    = string
  default = "%[1]s"
}

resource "huaweicloud_dataarts_studio_data_connection" "test" {
  count = var.data_connection_id == "" ? 1 : 0
  
  workspace_id = "%[2]s"
  type         = "DLI"
  name         = "%[3]s"
  env_type     = 0
  config       = jsonencode({
    "cdmPropertyEnable"             = false
    "metadata.collectionScope"      = ""
    "metadata.enableAutoCollection" = false
    "metadata.enableRealtimeSync"   = false
  })
  
  lifecycle {
    ignore_changes = [
      config,
    ]
  }
}
  
data "huaweicloud_dataarts_studio_data_connections" "test" {
  workspace_id  = "%[2]s"
  connection_id = var.data_connection_id != "" ? var.data_connection_id : huaweicloud_dataarts_studio_data_connection.test[0].id
}

resource "huaweicloud_dli_database" "test" {
  name        = "%[3]s"
  description = "Created by Terraform Script"
}
  
resource "huaweicloud_dli_table" "test" {
  database_name = huaweicloud_dli_database.test.name
  name          = "%[3]s"
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
  
resource "huaweicloud_dataarts_catalog_metadata_task" "test" {
  depends_on = [
    huaweicloud_dli_database.test,
    huaweicloud_dli_table.test,
  ]

  workspace_id = "%[2]s"
  name         = "%[3]s"
  dir_id       = "0"
  
  schedule_config {
    cron_expression = "00 */15 9-23 * * ?"
    max_time_out    = 60
    end_time        = format("%%s +08", split("+", timeadd("%[4]s", "24h"))[0])
    interval        = "15 minutes"
    schedule_type   = "CRON"
    start_time      = format("%%s +08", split("+", timeadd("%[4]s", "1h"))[0])
    enabled         = 0
  }
  
  data_source_type = "DLI"
  task_config      = jsonencode({
    data_connection_name        = try(data.huaweicloud_dataarts_studio_data_connections.test.connections[0].name, "")
    data_connection_id          = try(data.huaweicloud_dataarts_studio_data_connections.test.connections[0].id, "")
    data_connection_create_time = try(data.huaweicloud_dataarts_studio_data_connections.test.connections[0].create_timestamp, "")
    databaseName                = [
      huaweicloud_dli_database.test.name,
    ]
    tableName                   = [
      format("%%s.%%s", huaweicloud_dli_database.test.name, huaweicloud_dli_table.test.name),
    ]
    alive_object_policy         = "3"
    deleted_obkect_policy       = "1"
    deleted_object_policy       = "10"
    enableDataProfile           = true
    enableDataClassification    = false
    sampling                    = "10"
    queue                       = "default"
    unique                      = true
  })
  description      = "Created by terraform script"
}

resource "huaweicloud_dataarts_catalog_metadata_task_action" "runimmediate" {
  workspace_id = "%[2]s"
  task_id      = huaweicloud_dataarts_catalog_metadata_task.test.id
  action       = "runimmediate"
}
	
resource "time_sleep" "wait_120_seconds" {
  depends_on = [
    huaweicloud_dataarts_catalog_metadata_task_action.runimmediate,
  ]

  create_duration = "120s"
}

resource "huaweicloud_dataarts_catalog_metadata_task_action" "stop" {
  depends_on  = [
    time_sleep.wait_120_seconds,
  ]
  
  workspace_id = "%[2]s"
  task_id      = huaweicloud_dataarts_catalog_metadata_task.test.id
  action       = "stop"
}
`, acceptance.HW_DATAARTS_CONNECTION_ID, acceptance.HW_DATAARTS_WORKSPACE_ID, name, currentTime)
}

func testAccDataSourceCatalogMetadataTasks_basic_step2(name, currentTime string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dataarts_catalog_metadata_tasks" "all" {
  workspace_id = "%[2]s"
}

# Filter by user name
locals {
  user_name = huaweicloud_dataarts_catalog_metadata_task.test.user_name
}

data "huaweicloud_dataarts_catalog_metadata_tasks" "filter_by_user_name" {
  depends_on = [
    huaweicloud_dataarts_catalog_metadata_task.test,
  ]

  workspace_id = "%[2]s"
  user_name = local.user_name
}

locals {
  user_name_filter_result = [
    for v in data.huaweicloud_dataarts_catalog_metadata_tasks.filter_by_user_name.metadata_tasks : v.user_name == local.user_name
  ]
}
  
output "is_user_name_filter_useful" {
  value = length(local.user_name_filter_result) > 0 && alltrue(local.user_name_filter_result)
}

# Filter by name
locals {
  name = huaweicloud_dataarts_catalog_metadata_task.test.name
}

data "huaweicloud_dataarts_catalog_metadata_tasks" "filter_by_name" {
  depends_on = [
    huaweicloud_dataarts_catalog_metadata_task.test,
  ]

  workspace_id = "%[2]s"
  name 		   = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dataarts_catalog_metadata_tasks.filter_by_name.metadata_tasks : v.name == local.name
  ]
}
  
output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by data source type
locals {
  data_source_type = huaweicloud_dataarts_catalog_metadata_task.test.data_source_type
}

data "huaweicloud_dataarts_catalog_metadata_tasks" "filter_by_data_source_type" {
  depends_on = [
    huaweicloud_dataarts_catalog_metadata_task.test,
  ]

  workspace_id 	   = "%[2]s"
  data_source_type = local.data_source_type
}

locals {
  data_source_type_filter_result = [
    for v in data.huaweicloud_dataarts_catalog_metadata_tasks.filter_by_data_source_type.metadata_tasks : v.data_source_type == local.data_source_type
  ]
}
  
output "is_data_source_type_filter_useful" {
  value = length(local.data_source_type_filter_result) > 0 && alltrue(local.data_source_type_filter_result)
}

# Filter by data connection ID
locals {
  data_connection_id = try(data.huaweicloud_dataarts_studio_data_connections.test.connections[0].id, "")
  task_config        = huaweicloud_dataarts_catalog_metadata_task.test.task_config
}

data "huaweicloud_dataarts_catalog_metadata_tasks" "filter_by_data_connection_id" {
  depends_on = [
    huaweicloud_dataarts_catalog_metadata_task.test,
  ]

  workspace_id 	     = "%[2]s"
  data_connection_id = local.data_connection_id
}

locals {
  data_connection_id_filter_result = [
    for v in data.huaweicloud_dataarts_catalog_metadata_tasks.
	filter_by_data_connection_id.metadata_tasks : jsondecode(v.task_config).data_connection_id == local.data_connection_id
  ]
}

output "is_data_connection_id_filter_useful" {
  value = length(local.data_connection_id_filter_result) > 0 && alltrue(local.data_connection_id_filter_result)
}

# Filter by start time
locals {
  start_time = huaweicloud_dataarts_catalog_metadata_task.test.last_run_time
}

data "huaweicloud_dataarts_catalog_metadata_tasks" "filter_by_start_time" {
  depends_on = [
    huaweicloud_dataarts_catalog_metadata_task.test,
  ]

  workspace_id = "%[2]s"
  start_time   = local.start_time
}

locals {
  start_time_filter_result = [
    for v in data.huaweicloud_dataarts_catalog_metadata_tasks.filter_by_start_time.metadata_tasks : v.last_run_time >= local.start_time
  ]
}
  
output "is_start_time_filter_useful" {
  value = length(local.start_time_filter_result) > 0 && alltrue(local.start_time_filter_result)
}

# Filter by end time
locals {
  end_time = huaweicloud_dataarts_catalog_metadata_task.test.last_run_time
}

data "huaweicloud_dataarts_catalog_metadata_tasks" "filter_by_end_time" {
  depends_on = [
    huaweicloud_dataarts_catalog_metadata_task.test,
  ]

  workspace_id = "%[2]s"
  end_time     = local.end_time
}

locals {
  end_time_filter_result = [
    for v in data.huaweicloud_dataarts_catalog_metadata_tasks.filter_by_end_time.metadata_tasks : v.last_run_time <= local.end_time
  ]
}
  
output "is_end_time_filter_useful" {
  value = length(local.end_time_filter_result) > 0 && alltrue(local.end_time_filter_result)
}

# Filter by directory ID
locals {
  directory_id = huaweicloud_dataarts_catalog_metadata_task.test.dir_id
}

data "huaweicloud_dataarts_catalog_metadata_tasks" "filter_by_directory_id" {
  depends_on = [
    huaweicloud_dataarts_catalog_metadata_task.test,
  ]

  workspace_id = "%[2]s"
  directory_id = local.directory_id
}

locals {
  directory_id_filter_result = [
    for v in data.huaweicloud_dataarts_catalog_metadata_tasks.filter_by_directory_id.metadata_tasks : v.dir_id == local.directory_id
  ]
}
  
output "is_directory_id_filter_useful" {
  value = length(local.directory_id_filter_result) > 0 && alltrue(local.directory_id_filter_result)
}
`, testAccDataSourceCatalogMetadataTasks_basic_step1(name, currentTime), acceptance.HW_DATAARTS_WORKSPACE_ID)
}
