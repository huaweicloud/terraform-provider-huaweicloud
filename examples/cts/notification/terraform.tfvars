topic_name               = "tf_test_topic"
notification_name        = "tf_test_notification"
notification_agency_name = "cts_admin_trust"

notification_filter = [
  {
    condition = "OR"
    rule      = [
      "code = 400",
      "resource_name = name",
      "api_version = 1.0",
    ]
  }
]

notification_operations = [
  {
    service     = "ECS"
    resource    = "ecs"
    trace_names = [
      "createServer",
      "deleteServer",
    ]
  }
]

notification_operation_users = [
  {
    group = "devops"
    users = [
      "your_operation_user_name",
    ]
  }
]
