data "huaweicloud_secmaster_workspaces" "test" {
  count = var.workspace_id == "" ? 1 : 0

  name = var.workspace_name
}

resource "huaweicloud_secmaster_playbook" "test" {
  workspace_id = var.workspace_id != "" ? var.workspace_id : try(data.huaweicloud_secmaster_workspaces.test[0].workspaces[0].id, null)
  name         = var.playbook_name

  lifecycle {
    ignore_changes = [
      workspace_id
    ]
  }
}

data "huaweicloud_secmaster_data_classes" "test" {
  workspace_id = var.workspace_id != "" ? var.workspace_id : try(data.huaweicloud_secmaster_workspaces.test[0].workspaces[0].id, null)
}

resource "huaweicloud_secmaster_playbook_version" "test" {
  workspace_id      = var.workspace_id != "" ? var.workspace_id : try(data.huaweicloud_secmaster_workspaces.test[0].workspaces[0].id, null)
  playbook_id       = huaweicloud_secmaster_playbook.test.id
  dataclass_id      = try(data.huaweicloud_secmaster_data_classes.test.data_classes[0].id, null)
  rule_enable       = true
  trigger_type      = "EVENT"
  dataobject_create = true
  action_strategy   = "ASYNC"

  lifecycle {
    ignore_changes = [
      workspace_id,
      dataclass_id
    ]
  }
}

resource "huaweicloud_secmaster_playbook_rule" "test" {
  workspace_id    = var.workspace_id != "" ? var.workspace_id : try(data.huaweicloud_secmaster_workspaces.test[0].workspaces[0].id, null)
  version_id      = huaweicloud_secmaster_playbook_version.test.id
  expression_type = var.rule_expression_type

  dynamic "conditions" {
    for_each = var.rule_conditions
    content {
      name   = conditions.value.name
      detail = conditions.value.detail
      data   = conditions.value.data
    }
  }

  logics = split(",", join(",AND,", [
    for condition in var.rule_conditions : condition.name
  ])) # Using AND logic to combine conditions

  lifecycle {
    ignore_changes = [
      workspace_id
    ]
  }
}

data "huaweicloud_secmaster_workflows" "test" {
  workspace_id  = var.workspace_id != "" ? var.workspace_id : try(data.huaweicloud_secmaster_workspaces.test[0].workspaces[0].id, null)
  data_class_id = try(data.huaweicloud_secmaster_data_classes.test.data_classes[0].id, null)
}

resource "huaweicloud_secmaster_playbook_action" "test" {
  workspace_id = var.workspace_id != "" ? var.workspace_id : try(data.huaweicloud_secmaster_workspaces.test[0].workspaces[0].id, null)
  version_id   = huaweicloud_secmaster_playbook_version.test.id
  action_id    = try(data.huaweicloud_secmaster_workflows.test.workflows[0].id, null)
  name         = try(data.huaweicloud_secmaster_workflows.test.workflows[0].name, null)

  lifecycle {
    ignore_changes = [
      workspace_id,
      action_id,
      name
    ]
  }

  depends_on = [
    huaweicloud_secmaster_playbook_rule.test
  ]
}

resource "huaweicloud_secmaster_playbook_version_action" "test" {
  workspace_id = var.workspace_id != "" ? var.workspace_id : try(data.huaweicloud_secmaster_workspaces.test[0].workspaces[0].id, null)
  version_id   = huaweicloud_secmaster_playbook_version.test.id
  status       = "APPROVING"

  depends_on = [huaweicloud_secmaster_playbook_action.test]

  lifecycle {
    ignore_changes = [
      status,
      enabled
    ]
  }
}

resource "huaweicloud_secmaster_playbook_approval" "test" {
  workspace_id = var.workspace_id != "" ? var.workspace_id : try(data.huaweicloud_secmaster_workspaces.test[0].workspaces[0].id, null)
  version_id   = huaweicloud_secmaster_playbook_version.test.id
  result       = "PASS"
  content      = var.approval_content

  lifecycle {
    ignore_changes = [
      workspace_id
    ]
  }

  depends_on = [
    huaweicloud_secmaster_playbook_version_action.test
  ]
}

resource "huaweicloud_secmaster_playbook_enable" "test" {
  workspace_id      = var.workspace_id != "" ? var.workspace_id : try(data.huaweicloud_secmaster_workspaces.test[0].workspaces[0].id, null)
  playbook_id       = huaweicloud_secmaster_playbook.test.id
  playbook_name     = huaweicloud_secmaster_playbook.test.name
  active_version_id = huaweicloud_secmaster_playbook_approval.test.id

  lifecycle {
    ignore_changes = [
      workspace_id
    ]
  }
}
