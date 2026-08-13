# SecMaster workspace
resource "huaweicloud_secmaster_workspace" "test" {
  name         = var.workspace_name
  project_name = var.region_name
  description  = var.workspace_description
}

# SecMaster dataspace
resource "huaweicloud_secmaster_dataspace" "test" {
  workspace_id   = huaweicloud_secmaster_workspace.test.id
  dataspace_name = var.dataspace_name
  description    = var.dataspace_description
}

# SecMaster pipe
resource "huaweicloud_secmaster_pipe" "test" {
  workspace_id    = huaweicloud_secmaster_workspace.test.id
  dataspace_id    = huaweicloud_secmaster_dataspace.test.id
  pipe_name       = var.pipe_name
  shards          = var.shards
  storage_period  = var.storage_period
  description     = var.pipe_description
  timestamp_field = var.timestamp_field
  status          = var.pipe_status

  mapping = jsonencode({
    id   = {
      is_chinese_exist = true
      properties       = {}
      type             = "text"
    }
    name = {
      is_chinese_exist = false
      properties       = {}
      type             = "text"
    }
  })
}
