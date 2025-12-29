# Variable definitions for authentication
variable "region_name" {
  description = "The region where resources will be created"
  type        = string
}

# Variable definitions for sharer account
variable "access_key" {
  description = "The access key of the IAM user in sharer account"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "The secret key of the IAM user in sharer account"
  type        = string
  sensitive   = true
}

# Variable definitions for accepter account
variable "accepter_access_key" {
  description = "The access key of the IAM user in accepter account"
  type        = string
  sensitive   = true
}

variable "accepter_secret_key" {
  description = "The secret key of the IAM user in accepter account"
  type        = string
  sensitive   = true
}

# Variable definitions for instance flavor
variable "instance_flavor_id" {
  description = "The ID of the ECS instance flavor"
  type        = string
  default     = ""
  nullable    = true
}

variable "instance_flavor_performance_type" {
  description = "The performance type of the ECS instance flavor"
  type        = string
  default     = "normal"
}

variable "instance_flavor_cpu_core_count" {
  description = "The CPU core count of the ECS instance flavor"
  type        = number
  default     = 2
}

variable "instance_flavor_memory_size" {
  description = "The memory size of the ECS instance flavor (GB)"
  type        = number
  default     = 4
}

# Variable definitions for instance image
variable "instance_image_id" {
  description = "The ID of the ECS instance image"
  type        = string
  default     = ""
  nullable    = true
}

variable "instance_image_visibility" {
  description = "The visibility of the ECS instance image"
  type        = string
  default     = "public"
}

variable "instance_image_os" {
  description = "The OS of the ECS instance image"
  type        = string
  default     = "Ubuntu"
}

# Variable definitions for network
variable "vpc_name" {
  description = "The name of the VPC in sharer account"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC in sharer account"
  type        = string
  default     = "192.168.0.0/16"
}

variable "subnet_name" {
  description = "The name of the VPC subnet in sharer account"
  type        = string
}

variable "subnet_cidr" {
  description = "The CIDR block of the VPC subnet in sharer account"
  type        = string
  default     = ""
  nullable    = true
}

variable "subnet_gateway_ip" {
  description = "The gateway IP of the VPC subnet in sharer account"
  type        = string
  default     = ""
  nullable    = true
}

# Variable definitions for security group
variable "security_group_name" {
  description = "The name of the security group in sharer account"
  type        = string
}

# Variable definitions for ECS instance
variable "instance_name" {
  description = "The name of the ECS instance to be created in sharer account"
  type        = string
}

variable "administrator_password" {
  description = "The password of the administrator for the ECS instance"
  type        = string
  sensitive   = true
}

variable "instance_data_disks" {
  description = "The data disks of the ECS instance"
  type        = list(object({
    type = string
    size = number
  }))

  default  = []
  nullable = true
}

variable "vault_name" {
  description = "The name of the CBR vault in sharer account"
  type        = string
}

variable "vault_type" {
  description = "The type of the CBR vault"
  type        = string
  default     = "server"
}

variable "vault_consistent_level" {
  description = "The consistent level of the CBR vault"
  type        = string
  default     = "crash_consistent"
}

variable "vault_protection_type" {
  description = "The protection type of the CBR vault"
  type        = string
  default     = "backup"
}

variable "vault_size" {
  description = "The size of the CBR vault in GB"
  type        = number
  default     = 200
}

variable "whole_image_name" {
  description = "The name of the whole image to be created"
  type        = string
}

variable "whole_image_description" {
  description = "The description of the whole image"
  type        = string
  default     = ""
}

# Variable definitions for cross-account migration
variable "accepter_project_ids" {
  description = "The project IDs of accepter account for image sharing"
  type        = list(string)
}

variable "accepter_vault_name" {
  description = "The name of the CBR vault in accepter account"
  type        = string
}

variable "accepter_vault_type" {
  description = "The type of the CBR vault in accepter account"
  type        = string
  default     = "server"
}

variable "accepter_vault_consistent_level" {
  description = "The consistent level of the CBR vault in accepter account"
  type        = string
  default     = "crash_consistent"
}

variable "accepter_vault_protection_type" {
  description = "The protection type of the CBR vault in accepter account"
  type        = string
  default     = "backup"
}

variable "accepter_vault_size" {
  description = "The size of the CBR vault in accepter account in GB"
  type        = number
  default     = 200
}

variable "accepter_instance_name" {
  description = "The name of the new ECS instance to be created in accepter account"
  type        = string
}
