# Variable definitions for authentication
variable "region_name" {
  description = "The region where the COC script is located"
  type        = string
}

variable "access_key" {
  description = "The access key of the IAM user"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "The secret key of the IAM user"
  type        = string
  sensitive   = true
}

# Variable definitions for resources/data sources
variable "availability_zone" {
  description = "The availability zone to which the ECS instance and network belong"
  type        = string
  default     = ""
}

variable "instance_flavor_id" {
  description = "The flavor ID of the ECS instance"
  type        = string
  default     = ""
}

variable "instance_flavor_performance_type" {
  description = "The performance type of the ECS instance flavor"
  type        = string
  default     = "normal"
}

variable "instance_flavor_cpu_core_count" {
  description = "The number of the ECS instance flavor CPU cores"
  type        = number
  default     = 2
}

variable "instance_flavor_memory_size" {
  description = "The number of the ECS instance flavor memories"
  type        = number
  default     = 4
}

variable "instance_image_id" {
  description = "The image ID of the ECS instance"
  type        = string
  default     = ""
}

variable "instance_image_os_type" {
  description = "The OS type of the ECS instance flavor"
  type        = string
  default     = "Ubuntu"
}

variable "instance_image_visibility" {
  description = "The visibility of the ECS instance flavor"
  type        = string
  default     = "public"
}

variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "192.168.0.0/16"
}

variable "subnet_name" {
  description = "The name of the subnet"
  type        = string
}

variable "subnet_cidr" {
  description = "The CIDR block of the subnet"
  type        = string
  default     = ""
}

variable "subnet_gateway_ip" {
  description = "The gateway IP of the subnet"
  type        = string
  default     = ""
}

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "instance_name" {
  type        = string
  description = "The name of the ECS instance"
}

variable "instance_user_data" {
  type        = string
  description = "The user data for installing UniAgent on the ECS instance"
}

variable "script_name" {
  description = "The name of the script"
  type        = string
}

variable "script_description" {
  description = "The description of the script"
  type        = string
}

variable "script_risk_level" {
  description = "The description of the script"
  type        = string
}

variable "script_version" {
  description = "The description of the script"
  type        = string
}

variable "script_type" {
  description = "The risk level of the script"
  type        = string
}

variable "script_content" {
  description = "The content of the script"
  type        = string
}

variable "script_parameters" {
  description = "The parameter list of the script"
  type = list(object({
    name        = string
    value       = string
    description = string
    sensitive   = optional(bool)
  }))
}

variable "script_execute_timeout" {
  description = "The maximum time to execute the script in seconds"
  type        = number
}

variable "script_execute_user" {
  description = "The user to execute the script"
  type        = string
}

variable "script_execute_parameters" {
  description = "The parameter list of the script execution"
  type = list(object({
    name  = string
    value = string
  }))
}

variable "script_order_operation_batch_index" {
  description = "The batch index for the script order"
  type        = number
}

variable "script_order_operation_instance_id" {
  description = "The instance ID for the script order"
  type        = number
}

variable "script_order_operation_type" {
  description = "The operation type for the script order"
  type        = string
  default     = "CANCEL_ORDER"
}
