vpc_name            = "tf_test_kafka_instance"
subnet_name         = "tf_test_kafka_instance"
security_group_name = "tf_test_kafka_instance"
instance_name       = "tf_test_kafka_instance"
topic_name          = "tf_test_topic"
message_body        = "Hello Kafka!"

message_properties = [
  {
    name  = "KEY"
    value = "testKey"
  },
  {
    name  = "PARTITION"
    value = "1"
  }
]
