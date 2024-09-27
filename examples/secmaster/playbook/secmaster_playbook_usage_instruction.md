# Instructions for using playbook related resources

How to use terraform to orchestrate a playbook

1. Create a playbook
2. Create a playbook version for the playbook
3. Associate a trigger rule for the playbook version
4. Match a process (workflow) for playbook version
5. Submit draft version to approval status
6. Approve draft version to pass status
7. Enable playbook

-> The playbook resources can be updated only after the playbook is disabled.<br/>
   The playbook version and playbook action resources can be updated only after the version is deactivated.

## Example Usage

```hcl
variable "workspace_id" {}
variable "playbook_name" {}

resource "huaweicloud_secmaster_playbook" "test" {
  workspace_id = var.workspace_id
  name         = var.playbook_name
  description  = "created by terraform"
}

data "huaweicloud_secmaster_data_classes" "test" {
  workspace_id = var.workspace_id
}

resource "huaweicloud_secmaster_playbook_version" "test" {
  workspace_id = var.workspace_id
  playbook_id  = huaweicloud_secmaster_playbook.test.id
  dataclass_id = data.huaweicloud_secmaster_data_classes.test.data_classes[0].id
  description  = "created by terraform"
}

resource "huaweicloud_secmaster_playbook_rule" "test" {
  workspace_id    = var.workspace_id
  version_id      = huaweicloud_secmaster_playbook_version.test.id
  expression_type = "common"

  conditions {
    name   = "condition_0"
    detail = "123"
    data   = ["waf.alarm.level", ">", "3"]
  }
  
  logics = ["condition_0"]
}

data "huaweicloud_secmaster_workflows" "test" {
  workspace_id  = var.workspace_id
  data_class_id = data.huaweicloud_secmaster_data_classes.test.data_classes[0].id
}

resource "huaweicloud_secmaster_playbook_action" "test" {
  workspace_id = var.workspace_id
  version_id   = huaweicloud_secmaster_playbook_version.test.id
  action_id    = data.huaweicloud_secmaster_workflows.test.workflows[0].id
  name         = data.huaweicloud_secmaster_workflows.test.workflows[0].name
  description  = "created by terraform"
}

resource "huaweicloud_secmaster_playbook_version_action" "submit_version" {
  workspace_id = var.workspace_id
  version_id   = huaweicloud_secmaster_playbook_version.test.id
  status       = "APPROVING"

  depends_on = [huaweicloud_secmaster_playbook_action.test]
}

resource "huaweicloud_secmaster_playbook_approval" "pass" {
  workspace_id = var.workspace_id
  version_id   = huaweicloud_secmaster_playbook_version.test.id
  result       = "PASS"
  content      = "ok"

  depends_on = [huaweicloud_secmaster_playbook_version_action.submit_version]
}

resource "huaweicloud_secmaster_playbook_enable" "test" {
  workspace_id      = var.workspace_id
  playbook_id       = huaweicloud_secmaster_playbook.test.id
  playbook_name     = huaweicloud_secmaster_playbook.test.name
  active_version_id = huaweicloud_secmaster_playbook_approval.pass.id
}
```
