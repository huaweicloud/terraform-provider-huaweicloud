resource "huaweicloud_iotda_data_flow_control_policy" "test" {
  name        = var.flow_control_policy_name
  description = var.flow_control_policy_description
  scope       = "USER"
  limit       = var.flow_control_policy_limit
}

resource "huaweicloud_iotda_data_backlog_policy" "test" {
  name         = var.backlog_policy_name
  description  = var.backlog_policy_description
  backlog_size = var.backlog_policy_size
  backlog_time = var.backlog_policy_time
}
