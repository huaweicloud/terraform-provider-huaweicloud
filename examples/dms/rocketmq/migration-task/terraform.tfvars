vpc_name                 = "tf_test_vpc"
subnet_name              = "tf_test_subnet"
security_group_name      = "tf_test_security_group"
instance_name            = "tf_test_instance"
instance_broker_num      = 1
instance_description     = "Created by terraform script"
migration_task_name      = "tf_test_migration_task"
migration_task_overwrite = "true"
migration_task_type      = "rocketmq"

migration_task_topic_configs       = [
  {
    topic_name        = "tf_test_task_topic"
    topic_filter_type = "SINGLE_TAG"
    perm              = 6
    read_queue_nums   = 3
    write_queue_nums  = 3
  }
]
migration_task_subscription_groups = [
  {
    group_name      = "tf_test_task_group"
    consume_enable  = true
    retry_max_times = 16
  }
]
