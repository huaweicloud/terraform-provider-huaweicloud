variable "vpc_name" {
  description = "The name of the Huaweicloud VPC"
  default     = "tf_vpc_demo"
}

variable "subnet_name" {
  description = "The name of the Huaweicloud Subnet"
  default     = "tf_subnet_demo"
}

variable "security_group_name" {
  description = "The name of the Huaweicloud Security Group"
  default     = "tf_secgroup_demo"
}

variable "access_user_name" {
  description = "The access user of the Kafka instance"
  default     = "user"
}

variable "access_user_password" {
  description = "The access password of the Kafka instance"
  default     = "Kafkatest@123"
}

variable "manager_user" {
  description = "The manager user of the Kafka instance"
  default     = "kafka-user"
}

variable "manager_password" {
  description = "The manager password of the Kafka instance"
  default     = "Kafkatest@123"
}
