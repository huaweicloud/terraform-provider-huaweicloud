variable "enterprise_project_id" {
  description = "The ID of the enterprise project"
  type        = string
  default     = null
}

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

variable "ecs_instance_name" {
  type        = string
  description = "The name of the ECS instance"
}

variable "ecs_instance_user_data" {
  type        = string
  description = "The user data for installing UniAgent on the ECS instance"
}

variable "coc_script_name" {
  description = "The name of the script"
  type        = string
}

variable "coc_script_description" {
  description = "The description of the script"
  type        = string
}

variable "coc_script_risk_level" {
  description = "The description of the script"
  type        = string
}

variable "coc_script_version" {
  description = "The description of the script"
  type        = string
}

variable "coc_script_type" {
  description = "The risk level of the script"
  type        = string
}

variable "coc_script_content" {
  description = "The content of the script"
  type        = string
}

variable "coc_script_parameters" {
  description = "The parameter list of the script"
  type = list(object({
    name        = string
    value       = string
    description = string
    sensitive   = optional(bool)
  }))
}

variable "coc_script_execute_timeout" {
  description = "The maximum time to execute the script in seconds"
  type        = number
}

variable "coc_script_execute_user" {
  description = "The user to execute the script"
  type        = string
}

variable "coc_script_execute_parameters" {
  description = "The parameter list of the script execution"
  type = list(object({
    name  = string
    value = string
  }))
}

variable "coc_script_order_operation_batch_index" {
  description = "The batch index for the script order"
  type        = number
}

variable "coc_script_order_operation_instance_id" {
  description = "The instance ID for the script order"
  type        = number
}

variable "coc_script_order_operation_type" {
  description = "The operation type for the script order"
  type        = string
  default     = "CANCEL_ORDER"
}

