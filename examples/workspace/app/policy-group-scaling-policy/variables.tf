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

# Variable definitions for APP group
variable "app_group_name" {
  description = "The name of the APP group"
  type        = string
}

# Variable definitions for policy group
variable "policy_group_name" {
  description = "The name of the policy group"
  type        = string
}

variable "policy_group_priority" {
  description = "The priority of the policy group"
  type        = number
  default     = 1
}

variable "policy_group_description" {
  description = "The description of the policy group"
  type        = string
  default     = "Created APP policy group by Terraform"
}

variable "target_type" {
  description = "The type of target for the policy group"
  type        = string
  default     = "APPGROUP"

  validation {
    condition     = contains(["APPGROUP", "ALL"], var.target_type)
    error_message = "The target_type must be either 'APPGROUP' or 'ALL'."
  }
}

# Variable definitions for policy settings
variable "automatic_reconnection_interval" {
  description = "The automatic reconnection interval in minutes"
  type        = number
  default     = 10

  validation {
    condition     = var.automatic_reconnection_interval >= 1 && var.automatic_reconnection_interval <= 60
    error_message = "The automatic_reconnection_interval must be between 1 and 60 minutes."
  }
}

variable "session_persistence_time" {
  description = "The session persistence time in minutes"
  type        = number
  default     = 120

  validation {
    condition     = var.session_persistence_time >= 1 && var.session_persistence_time <= 1440
    error_message = "The session_persistence_time must be between 1 and 1440 minutes."
  }
}

variable "forbid_screen_capture" {
  description = "Whether to forbid screen capture"
  type        = bool
  default     = true
}

# Variable definitions for scaling policy
variable "max_scaling_amount" {
  description = "The maximum number of instances that can be scaled out"
  type        = number

  validation {
    condition     = var.max_scaling_amount >= 1 && var.max_scaling_amount <= 100
    error_message = "The max_scaling_amount must be between 1 and 100."
  }
}

variable "single_expansion_count" {
  description = "The number of instances to scale out in a single scaling operation"
  type        = number

  validation {
    condition     = var.single_expansion_count >= 1 && var.single_expansion_count <= 10
    error_message = "The single_expansion_count must be between 1 and 10."
  }
}

variable "session_usage_threshold" {
  description = "The session usage threshold percentage"
  type        = number
  default     = 80

  validation {
    condition     = var.session_usage_threshold >= 1 && var.session_usage_threshold <= 100
    error_message = "The session_usage_threshold must be between 1 and 100."
  }
}

variable "shrink_after_session_idle_minutes" {
  description = "The number of minutes to wait before shrinking idle instances"
  type        = number
  default     = 30

  validation {
    condition     = var.shrink_after_session_idle_minutes >= 1 && var.shrink_after_session_idle_minutes <= 1440
    error_message = "The shrink_after_session_idle_minutes must be between 1 and 1440 minutes."
  }
}
