vpc_name             = "tf_test_rocketmq_message_send"
subnet_name          = "tf_test_rocketmq_message_send"
security_group_name  = "tf_test_rocketmq_message_send"
instance_name        = "tf_test_rocketmq_message_send"
instance_broker_num  = 1
instance_description = "Created by terraform script"
topic_name           = "tf_test_rocketmq_message_send"
topic_queue_num      = 3
topic_permission     = "all"
message_body         = "tf terraform script test"

message_properties = [
  {
    name  = "KEYS"
    value = "owner"
  },
  {
    name  = "TAGS"
    value = "terraform"
  }
]
