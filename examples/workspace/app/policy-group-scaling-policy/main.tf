# Get Workspace service information
data "huaweicloud_workspace_service" "test" {}

# Create APP server group
resource "huaweicloud_workspace_app_server_group" "test" {
  name             = var.app_server_group_name
  app_type         = var.app_server_group_app_type
  os_type          = var.app_server_group_os_type
  flavor_id        = var.app_server_group_flavor_id
  image_type       = "gold"
  image_id         = var.app_server_group_image_id
  image_product_id = var.app_server_group_image_product_id
  vpc_id           = data.huaweicloud_workspace_service.test.vpc_id
  subnet_id        = try(data.huaweicloud_workspace_service.test.network_ids[0], null)
  system_disk_type = var.app_server_group_system_disk_type
  system_disk_size = var.app_server_group_system_disk_size
  is_vdi           = true
}

# Create APP group
resource "huaweicloud_workspace_app_group" "test" {
  depends_on = [huaweicloud_workspace_app_server_group.test]

  server_group_id = huaweicloud_workspace_app_server_group.test.id
  name            = var.app_group_name
  type            = "SESSION_DESKTOP_APP"
  description     = "Created APP group by Terraform"
}

# Create APP policy group
resource "huaweicloud_workspace_app_policy_group" "test" {
  depends_on = [huaweicloud_workspace_app_group.test]

  name        = var.policy_group_name
  priority    = var.policy_group_priority
  description = var.policy_group_description

  targets {
    id   = var.target_type == "APPGROUP" ? huaweicloud_workspace_app_group.test.id : "default-apply-all-targets"
    name = var.target_type == "APPGROUP" ? huaweicloud_workspace_app_group.test.name : "All-Targets"
    type = var.target_type
  }

  policies = jsonencode({
    "client": {
      "automatic_reconnection_interval": var.automatic_reconnection_interval,
      "session_persistence_time":        var.session_persistence_time,
      "forbid_screen_capture":           var.forbid_screen_capture
    }
  })
}

# Create APP server group scaling policy
resource "huaweicloud_workspace_app_server_group_scaling_policy" "test" {
  depends_on = [huaweicloud_workspace_app_server_group.test]

  server_group_id        = huaweicloud_workspace_app_server_group.test.id
  max_scaling_amount     = var.max_scaling_amount
  single_expansion_count = var.single_expansion_count

  scaling_policy_by_session {
    session_usage_threshold           = var.session_usage_threshold
    shrink_after_session_idle_minutes = var.shrink_after_session_idle_minutes
  }
}
