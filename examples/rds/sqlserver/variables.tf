variable "vpc_name" {
  description = "The name of the HuaweiCloud VPC"
  default     = "tf_vpc_demo"
}

variable "subnet_name" {
  description = "The name of the HuaweiCloud Subnet"
  default     = "tf_subnet_demo"
}

variable "security_group_name" {
  description = "The name of the HuaweiCloud Security Group"
  default     = "tf_secgroup_demo"
}

variable "single_instance_name" {
  description = "The name of the HuaweiCloud RDS single instance"
  default     = "tf_sqlserver_single_instance_demo"
}

variable "ha_instance_name" {
  description = "The name of the HuaweiCloud RDS HA instance"
  default     = "tf_sqlserver_ha_instance_demo"
}
