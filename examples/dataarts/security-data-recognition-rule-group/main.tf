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

locals {
  workspace_id = var.workspace_id != "" ? var.workspace_id : try(data.huaweicloud_dataarts_studio_workspaces.test[0].workspaces[0].id, "")
}

# Query data categories in the workspace to resolve category IDs for data recognition rules.
data "huaweicloud_dataarts_security_data_categories" "test" {
  workspace_id = local.workspace_id
}

locals {
  category_ids = length(var.category_ids) > 0 ? var.category_ids : slice([
    for category in data.huaweicloud_dataarts_security_data_categories.test.categories : category.category_id
  ], 0, var.data_recognition_rule_count)
}

# Create data secrecy levels for data recognition rules.
resource "huaweicloud_dataarts_security_data_secrecy_level" "test" {
  count = var.data_recognition_rule_count

  workspace_id = local.workspace_id
  name         = format("%s_secrecy_level_%d", var.rule_group_name, count.index)
}

# Create data recognition rules that will be grouped together.
resource "huaweicloud_dataarts_security_data_recognition_rule" "test" {
  count = var.data_recognition_rule_count

  workspace_id     = local.workspace_id
  rule_type        = "CUSTOM"
  name             = format("%s_rule_%d", var.rule_group_name, count.index)
  secrecy_level_id = huaweicloud_dataarts_security_data_secrecy_level.test[count.index].id
  category_id      = local.category_ids[count.index]
  method           = "NONE"

  lifecycle {
    precondition {
      condition     = length(local.category_ids) >= var.data_recognition_rule_count
      error_message = "The number of available category IDs must be greater than or equal to data_recognition_rule_count."
    }
  }
}

# Create a data recognition rule group that contains the created rules.
resource "huaweicloud_dataarts_security_data_recognition_rule_group" "test" {
  depends_on = [
    huaweicloud_dataarts_security_data_recognition_rule.test,
  ]

  workspace_id = local.workspace_id
  name         = var.rule_group_name
  description  = var.rule_group_description
  rule_ids     = huaweicloud_dataarts_security_data_recognition_rule.test[*].id
}

# Query the created rule group to verify the group metadata after creation.
data "huaweicloud_dataarts_security_data_recognition_rule_groups" "test" {
  depends_on = [
    huaweicloud_dataarts_security_data_recognition_rule_group.test,
  ]

  workspace_id = local.workspace_id
  name         = var.rule_group_name
}
