# Variable definitions for authentication
variable "region_name" {
  description = "The region where the LTS service is located"
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
variable "is_instance_create" {
  description = "Whether to create the identity center instance"
  type        = bool
  default     = true
}

variable "is_region_need_register" {
  description = "Whether to register the region"
  type        = bool
  default     = true
}

variable "instance_store_id_alias" {
  description = "The alias of the identity center instance"
  type        = string
  default     = ""
}

# SC.005 Disable
variable "policy_max_password_age" {
  description = "The max password age of the identity center password policy, unit in days"
  type        = number
  default     = 10

  validation {
    condition     = var.policy_max_password_age >= 1 && var.policy_max_password_age <= 1095
    error_message = "The valid value of the max password age is range from 1 to 1095"
  }
}

variable "policy_minimum_password_length" {
  description = "The minimum password length of the identity center password policy"
  type        = number
  default     = 10
}

variable "policy_password_reuse_prevention" {
  description = "The password reuse prevention feature of the identity center's password policy indicates whether to prohibit the use of the same password as the previous one"
  type        = bool
  default     = true
}
# SC.005 Enable

variable "policy_require_uppercase_characters" {
  description = "The require uppercase characters of the identity center password policy"
  type        = bool
  default     = true
}

variable "policy_require_lowercase_characters" {
  description = "The require lowercase characters of the identity center password policy"
  type        = bool
  default     = true
}

variable "policy_require_numbers" {
  description = "The require numbers of the identity center password policy"
  type        = bool
  default     = true
}

variable "policy_require_symbols" {
  description = "The require symbols of the identity center password policy"
  type        = bool
  default     = true
}
