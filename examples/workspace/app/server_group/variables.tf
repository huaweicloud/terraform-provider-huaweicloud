# Variable definitions for authentication
variable "region_name" {
  description = "The region where the Workspace service is located"
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

# Variable definitions for APP server group
variable "app_server_group_name" {
  description = "The name of the APP server group"
  type        = string
}

variable "app_server_group_app_type" {
  description = "The application type of the APP server group"
  type        = string
  default     = "SESSION_DESKTOP_APP"
}

variable "app_server_group_os_type" {
  description = "The operating system type of the APP server group"
  type        = string
  default     = "Windows"
}

variable "app_server_group_flavor_id" {
  description = "The flavor ID of the APP server group"
  type        = string
}

variable "app_server_group_image_id" {
  description = "The image ID of the APP server group"
  type        = string
}

variable "app_server_group_image_product_id" {
  description = "The image product ID of the APP server group"
  type        = string
}

variable "app_server_group_system_disk_type" {
  description = "The system disk type of the APP server group"
  type        = string
  default     = "SAS"
}

variable "app_server_group_system_disk_size" {
  description = "The system disk size of the APP server group in GB"
  type        = number
  default     = 80

  validation {
    condition     = var.app_server_group_system_disk_size >= 40 && var.app_server_group_system_disk_size <= 2048
    error_message = "The system disk size must be between 40 and 2048 GB."
  }
}
