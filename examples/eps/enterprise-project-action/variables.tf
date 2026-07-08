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
variable "enterprise_project_id" {
  description = "The ID of the enterprise project"
  type        = string
  default     = ""
  nullable    = false
}

variable "enterprise_project_name" {
  description = "The name of the enterprise project"
  type        = string
  default     = ""
  nullable    = false

  validation {
    condition     = var.enterprise_project_id != "" || var.enterprise_project_name != ""
    error_message = "enterprise_project_name must be provided if enterprise_project_id is not provided."
  }
}

variable "enterprise_project_description" {
  description = "The description of the enterprise project"
  type        = string
  default     = ""
}

variable "enterprise_project_type" {
  description = "The type of the enterprise project"
  type        = string
  default     = "prod"
}

variable "enterprise_project_enable" {
  description = "Whether to enable the enterprise project"
  type        = bool
  default     = true
}

variable "delete_flag" {
  description = "Whether to delete the enterprise project on destroy"
  type        = bool
  default     = true
}

variable "enterprise_project_action" {
  description = "The action to perform on the enterprise project. Valid values are enable and disable"
  type        = string
  default     = "disable"
}
