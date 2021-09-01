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

variable "waf_dedicated_instance_name" {
  description = "The name of the WAF Dedicated Instance"
  default     = "tf_instance_demo"
}

variable "waf_policy_name" {
  description = "The name of the WAF Policy"
  default     = "tf_policy_demo"
}
