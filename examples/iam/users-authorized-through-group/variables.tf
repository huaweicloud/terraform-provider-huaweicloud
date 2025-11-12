# Variables definitions for authorization
variable "region_name" {
  description = "The region where the scaling configuration is located"
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

# Variables definitions for resource/data source
variable "role_id" {
  description = "The ID of the IAM role"
  type        = string
  default     = ""
  nullable    = false
}

variable "role_policy" {
  description = "The policy of the IAM role"
  type        = string
  default     = ""
  nullable    = false
}

variable "role_name" {
  description = "The name of the IAM role"
  type        = string
  default     = ""
  nullable    = false

  validation {
    condition     = !(var.role_name == "" && var.role_id == "")
    error_message = "The role_name must be provided when role_id is not provided."
  }
}

variable "role_type" {
  description = "The type of the IAM role"
  type        = string
  default     = "XA"
}

variable "role_description" {
  description = "The description of the IAM role"
  type        = string
  default     = ""

  validation {
    condition     = !(var.role_description == "" && var.role_policy != "")
    error_message = "The role_description must be provided when role_policy is provided."
  }
}

variable "group_id" {
  description = "The ID of the IAM group"
  type        = string
  default     = ""
  nullable    = false
}

variable "group_name" {
  description = "The name of the IAM group"
  type        = string
  default     = ""
  nullable    = false

  validation {
    condition     = !(var.group_name == "" && var.group_id == "")
    error_message = "The group_name must be provided when group_id is not provided."
  }
}

variable "group_description" {
  description = "The description of the IAM group"
  type        = string
  default     = ""
}

variable "authorized_project_id" {
  description = "The ID of the IAM project"
  type        = string
  default     = ""
  nullable    = false
}

variable "authorized_project_name" {
  description = "The name of the IAM project"
  type        = string
  default     = ""
  nullable    = true

  validation {
    condition     = !(var.authorized_project_name == "" && var.authorized_project_id == "")
    error_message = "The authorized_project_name must be provided when authorized_project_id is not provided."
  }
}

variable "authorized_domain_id" {
  description = "The ID of the IAM domain"
  type        = string
  default     = ""
  nullable    = false
}

variable "users_configuration" {
  description = "The configuration of the IAM users"
  type        = list(object({
    name     = string
    password = optional(string, "")
  }))
  nullable    = false
}
