# Query workspaces under the specified DataArts Studio instance.
data "huaweicloud_dataarts_studio_workspaces" "test" {
  count = var.workspace_id == "" ? 1 : 0

  instance_id  = var.instance_id
  name         = var.workspace_name != "" ? var.workspace_name : null
  workspace_id = var.workspace_id != "" ? var.workspace_id : null

  lifecycle {
    precondition {
      condition     = var.instance_id != ""
      error_message = "instance_id must be provided if workspace_id is omitted."
    }
  }
}

# Query data connections in the workspace to obtain the DLI connection name.
data "huaweicloud_dataarts_studio_data_connections" "test" {
  count = var.connection_name == "" ? 1 : 0

  workspace_id = var.workspace_id != "" ? var.workspace_id : try(data.huaweicloud_dataarts_studio_workspaces.test[0].workspaces[0].id, "")
  type         = "DLI"
}

# Create a DLI database as the data source for the script.
resource "huaweicloud_dli_database" "test" {
  name        = var.dli_database_name
  description = var.dli_database_description
}

# Create a DLI table as the data source for the script.
resource "huaweicloud_dli_table" "test" {
  database_name = huaweicloud_dli_database.test.name
  name          = var.dli_table_name
  data_location = "DLI"
  description   = var.dli_table_description

  dynamic "columns" {
    for_each = var.dli_table_columns

    content {
      name = columns.value.name
      type = columns.value.type
    }
  }
}

# Create a DataArts Factory script that queries the DLI table.
resource "huaweicloud_dataarts_factory_script" "test" {
  depends_on = [
    huaweicloud_dli_database.test,
    huaweicloud_dli_table.test,
  ]

  workspace_id    = var.workspace_id != "" ? var.workspace_id : try(data.huaweicloud_dataarts_studio_workspaces.test[0].workspaces[0].id, "")
  name            = var.script_name
  type            = var.script_type
  connection_name = var.connection_name != "" ? var.connection_name : try(data.huaweicloud_dataarts_studio_data_connections.test[0].connections[0].name, "")
  directory       = var.script_directory
  database        = huaweicloud_dli_database.test.name
  queue_name      = var.queue_name
  description     = var.script_description
  configuration   = var.script_configuration
  content         = var.script_content != "" ? var.script_content : format("SELECT * FROM %s.%s", var.dli_database_name, var.dli_table_name)
}

# Execute the DataArts Factory script once.
resource "huaweicloud_dataarts_factory_script_execute" "test" {
  depends_on = [
    huaweicloud_dataarts_factory_script.test,
  ]

  workspace_id = var.workspace_id != "" ? var.workspace_id : try(data.huaweicloud_dataarts_studio_workspaces.test[0].workspaces[0].id, "")
  script_name  = huaweicloud_dataarts_factory_script.test.name
  params       = length(var.script_execute_params) > 0 ? jsonencode(var.script_execute_params) : null
}
