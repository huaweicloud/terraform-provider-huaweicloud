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
variable "script_name" {
  description = "The name of the script"
  type        = string
}

variable "script_description" {
  description = "The description of the script"
  type        = string
}

variable "script_risk_level" {
  description = "The risk level of the script"
  type        = string
}

variable "script_version" {
  description = "The version of the script"
  type        = string
}

variable "script_type" {
  description = "The type of the script"
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

  nullable = false
}
