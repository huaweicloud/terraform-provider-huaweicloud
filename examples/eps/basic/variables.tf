# Variable definitions for authentication
variable "region_name" {
  description = "The region where the EPS service is located"
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

# Variable definitions for resources
variable "enterprise_project_name" {
  description = "The name of the enterprise project"
  type        = string
}

variable "enterprise_project_description" {
  description = "The description of the enterprise project"
  type        = string
  default     = ""
}

variable "enterprise_project_type" {
  description = "The type of the enterprise project. Valid values are poc and prod"
  type        = string
  default     = "prod"
}

variable "enterprise_project_enable" {
  description = "Whether to enable the enterprise project"
  type        = bool
  default     = true
}

variable "skip_disable_on_destroy" {
  description = "Whether to skip disabling the enterprise project on destroy"
  type        = bool
  default     = false
}

variable "delete_flag" {
  description = "Whether to delete the enterprise project on destroy"
  type        = bool
  default     = true
}
