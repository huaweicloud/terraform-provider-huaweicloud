# Variable definitions for authentication
variable "region_name" {
  description = "The region where the ECS instance is located"
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
  description = "The availability zone to which the ECS instance belongs"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_flavor_id" {
  description = "The flavor ID of the ECS instance"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_performance_type" {
  description = "The performance type of the ECS instance flavor"
  type        = string
  default     = "normal"
}

variable "instance_cpu_core_count" {
  description = "The number of CPU cores of the ECS instance"
  type        = number
  default     = 2
}

variable "instance_memory_size" {
  description = "The memory size in GB of the ECS instance"
  type        = number
  default     = 4
}

variable "instance_image_id" {
  description = "The image ID of the ECS instance"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_image_visibility" {
  description = "The visibility of the ECS instance image"
  type        = string
  default     = "public"
}

variable "instance_image_os" {
  description = "The operating system of the ECS instance image"
  type        = string
  default     = "Ubuntu"
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

variable "subnet_configurations" {
  description = "The list of subnet configurations for ECS instance."
  type        = list(object({
    subnet_name       = string
    subnet_cidr       = optional(string)
    subnet_gateway_ip = optional(string)
  }))
}

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "instance_name" {
  description = "The name of the ECS instance"
  type        = string
}

variable "instance_admin_password" {
  description = "The login password of the ECS instance"
  type        = string
  sensitive   = true
}

variable "attached_network_id" {
  description = "The ID of the network to which the ECS instance to be attached"
  type        = string
  default     = ""
  nullable    = false

  validation {
    condition     = var.attached_network_id != "" || length(var.subnet_configurations) == 2
    error_message = "When attached_network_id is not provided, subnet_configurations must have exactly 2 elements."
  }
}

variable "attached_interface_fixed_ip" {
  description = "The fixed IP address of the ECS instance to be attached"
  type        = string
  default     = null
}

variable "attached_security_group_ids" {
  description = "The list of security group IDs of the ECS instance to be attached"
  type        = list(string)
  default     = null
}
