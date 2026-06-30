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

# Query available workspace user roles under the target workspace.
data "huaweicloud_dataarts_studio_workspace_user_roles" "test" {
  workspace_id = var.workspace_id != "" ? var.workspace_id : try(data.huaweicloud_dataarts_studio_workspaces.test[0].workspaces[0].id, "")
}

# Query the IAM user by name when user_id is omitted.
data "huaweicloud_identity_users" "test" {
  count = var.user_id == "" ? 1 : 0

  name = var.user_name

  lifecycle {
    precondition {
      condition     = var.user_name != ""
      error_message = "user_name must be provided if user_id is omitted."
    }
  }
}

# Create a workspace user and assign roles.
resource "huaweicloud_dataarts_studio_workspace_user" "test" {
  workspace_id = var.workspace_id != "" ? var.workspace_id : try(data.huaweicloud_dataarts_studio_workspaces.test[0].workspaces[0].id, "")
  user_id      = var.user_id != "" ? var.user_id : try(data.huaweicloud_identity_users.test[0].users[0].id, "")

  dynamic "roles" {
    for_each = var.role_ids

    content {
      id = roles.value
    }
  }

  # ST.003 Disable
  lifecycle {
    precondition {
      condition     = length([
        for role_id in var.role_ids : role_id
        if contains([for role in data.huaweicloud_dataarts_studio_workspace_user_roles.test.roles : role.id], role_id)
      ]) == length(var.role_ids)
      error_message = "All role_ids must exist in the workspace user roles returned by the data source."
    }
  }
  # ST.003 Enable
}
