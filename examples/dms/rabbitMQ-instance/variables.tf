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
