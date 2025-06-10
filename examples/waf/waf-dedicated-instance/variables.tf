variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
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

variable "waf_dedicated_instance_name" {
  description = "The name of the dedicated instance"
  type        = string
}

variable "waf_dedicated_instance_specification_code" {
  description = "The specification code of the dedicated instance"
  type        = string
}

variable "waf_policy_name" {
  description = "The name of the WAF policy"
  type        = string
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project"
  type        = string
  default     = null
}
