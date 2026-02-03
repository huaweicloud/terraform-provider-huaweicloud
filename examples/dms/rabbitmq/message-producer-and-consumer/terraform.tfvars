# Network variables
vpc_name                  = "tf_test_vpc_rabbitmq"
subnet_name               = "tf_test_subnet_rabbitmq"
security_group_name       = "tf_test_sg_rabbitmq"
# RabbitMQ instance variables
instance_name             = "tf_test_rabbitmq_instance"
instance_access_user_name = "admin"
instance_password         = "YourPassword@123"
# ECS instance variables
producer_instance_name    = "tf_test_producer"
consumer_instance_name    = "tf_test_consumer"
