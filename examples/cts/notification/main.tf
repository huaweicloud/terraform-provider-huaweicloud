resource "huaweicloud_smn_topic" "test" {
  name = var.topic_name
}

resource "huaweicloud_cts_notification" "test" {
  name           = var.notification_name
  operation_type = var.notification_operation_type
  smn_topic      = huaweicloud_smn_topic.test.id
  agency_name    = var.notification_agency_name

  dynamic "filter" {
    for_each = var.notification_filter

    content {
      condition = filter.value.condition
      rule      = filter.value.rule
    }
  }

  dynamic "operations" {
    for_each = var.notification_operations

    content {
      service     = operations.value.service
      resource    = operations.value.resource
      trace_names = operations.value.trace_names
    }
  }

  dynamic "operation_users" {
    for_each = var.notification_operation_users

    content {
      group = operation_users.value.group
      users = operation_users.value.users
    }
  }
}
