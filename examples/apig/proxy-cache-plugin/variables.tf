variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
}

variable "subnet_name" {
  description = "The name of the subnet"
  type        = string
}

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "instance_name" {
  description = "The name of the APIG instance"
  type        = string
}

variable "plugin_name" {
  description = "The name of the APIG plugin"
  type        = string
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project"
  type        = string
}
