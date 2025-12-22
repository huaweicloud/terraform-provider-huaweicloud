# Network configuration
vpc_name    = "tf_test_vpc"
subnet_name = "tf_test_subnet"

# Security group
security_group_name = "tf_test_security_group"

# APIG instance configuration
instance_name         = "tf_test_apig_instance"
instance_edition      = "BASIC"
enterprise_project_id = "0"

# Plugin configuration
plugin_name        = "tf_test_kafka_forward_plugin"
plugin_description = "Kafka forward plugin created by Terraform script"

# Kafka instance configuration
kafka_instance_name              = "tf_test_kafka_instance"
kafka_instance_description       = "Kafka instance for testing"
kafka_instance_flavor_type       = "cluster"
kafka_instance_storage_spec_code = "dms.physical.storage.high.v2"
kafka_instance_engine_version    = "2.7"
kafka_instance_storage_space     = 600
kafka_instance_broker_num        = 3
kafka_instance_ssl_enable        = false
kafka_instance_user_name         = "user"
kafka_instance_user_password     = "Kafkatest@123"

# Kafka charging configuration
kafka_charging_mode = "prePaid"
kafka_period_unit   = "month"
kafka_period        = 1
kafka_auto_new      = "false"

# Kafka topic configuration
kafka_topic_name       = "tf_test_kafka_topic"
kafka_topic_partitions = 1

# Kafka plugin configuration
kafka_message_key     = "terraform-test"  # Use request ID as message key
kafka_max_retry_count = 3
kafka_retry_backoff   = 10

# Kafka security configuration
kafka_security_protocol = "PLAINTEXT"  # Options: PLAINTEXT, SASL_PLAINTEXT, SASL_SSL, SSL
kafka_sasl_mechanisms   = "PLAIN"      # Options: PLAIN, SCRAM-SHA-256, SCRAM-SHA-512
kafka_access_user       = "user"
kafka_password          = "Kafkatest@123"
