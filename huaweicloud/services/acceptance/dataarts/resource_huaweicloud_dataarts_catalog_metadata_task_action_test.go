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

func TestAccCatalogMetadataTaskAction_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_dataarts_catalog_metadata_task.test"
		rcActionTask = acceptance.InitResourceCheck(resourceName, &obj, getCatalogMetadataTaskResourceFunc)

		name        = acceptance.RandomAccResourceName()
		currentTime = time.Now().Local().Format(time.RFC3339)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.12.1",
			},
		},
		// This resource is a one-time action resource and the instance record will not delete even if resource is deleted.
		// So we need to ignore the check destroy.
		// lintignore:AT001
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcActionTask.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config:      testAccCatalogMetadataTaskAction_nonExistentMetadataTask(),
				ExpectError: regexp.MustCompile(`error action metadata task`),
			},
			{
				Config: testAccCatalogMetadataTaskAction_basic(name, currentTime),
				Check: resource.ComposeTestCheckFunc(
					rcActionTask.CheckResourceExists(),
				),
			},
		},
	})
}

func testAccCatalogMetadataTaskAction_nonExistentMetadataTask() string {
	randUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
resource "huaweicloud_dataarts_catalog_metadata_task_action" "non_existent_metadata_task" {
  workspace_id = "%[1]s"
  task_id      = "%[2]s"
  action       = "run"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, randUUID.String())
}

func testAccCatalogMetadataTaskAction_basic_base(name, currentTime string) string {
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
`, acceptance.HW_DATAARTS_CONNECTION_ID, acceptance.HW_DATAARTS_WORKSPACE_ID, name, currentTime)
}

func testAccCatalogMetadataTaskAction_basic(name, currentTime string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_catalog_metadata_task_action" "run" {
  workspace_id = "%[2]s"
  task_id      = huaweicloud_dataarts_catalog_metadata_task.test.id
  action       = "run"
}

resource "time_sleep" "wait_30_seconds" {
  depends_on = [huaweicloud_dataarts_catalog_metadata_task_action.run]

  create_duration = "30s"
}

resource "huaweicloud_dataarts_catalog_metadata_task_action" "stop" {
  depends_on  = [time_sleep.wait_30_seconds]

  workspace_id = "%[2]s"
  task_id      = huaweicloud_dataarts_catalog_metadata_task.test.id
  action       = "stop"
}
`, testAccCatalogMetadataTaskAction_basic_base(name, currentTime), acceptance.HW_DATAARTS_WORKSPACE_ID)
}
