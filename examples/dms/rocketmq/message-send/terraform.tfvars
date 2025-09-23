vpc_name             = "tf_test_vpc"
subnet_name          = "tf_test_subnet"
security_group_name  = "tf_test_security_group"
instance_name        = "tf_test_instance"
instance_broker_num  = 1
instance_description = "Created by terraform script"
topic_name           = "tf_test_topic"
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
