vpc_name            = "tf_test_kafka_instance"
subnet_name         = "tf_test_kafka_instance"
security_group_name = "tf_test_kafka_instance"
task_name           = "tf_test_kafka_task"
topic_name          = "tf_test_kafka_topic"

instance_configurations = [
  {
    name = "tf_test_instance"
  },
  {
    name               = "tf_test_peer_instance"
    access_user        = "admin"
    password           = "YourKafkaInstancePassword!"
    enabled_mechanisms = ["SCRAM-SHA-512"]
    port_protocol      = {
      private_plain_enable    = false
      private_sasl_ssl_enable = true
    }
  }
]
