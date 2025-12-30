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
variable "policy_max_consecutive_identical_chars" {
  description = "The maximum number of times that a character is allowed to consecutively present in a password"
  type        = number

  validation {
    condition     = var.policy_max_consecutive_identical_chars >= 0 && var.policy_max_consecutive_identical_chars <= 32
    error_message = "The valid value of the maximum number of times that a character is allowed to consecutively present in a password is range from 0 to 32"
  }
}

# SC.005 Disable
variable "policy_min_password_age" {
  description = "The minimum period (minutes) after which users are allowed to make a password change"
  type        = number

  validation {
    condition     = var.policy_min_password_age >= 0 && var.policy_min_password_age <= 1440
    error_message = "The valid value of the minimum period (minutes) after which users are allowed to make a password change is range from 0 to 1440"
  }
}

variable "policy_min_password_length" {
  description = "The minimum number of characters that a password must contain"
  type        = number

  validation {
    condition     = var.policy_min_password_length >= 8 && var.policy_min_password_length <= 32
    error_message = "The valid value of the minimum number of characters that a password must contain is range from 8 to 32"
  }
}

variable "policy_password_reuse_prevention" {
  description = "The password reuse prevention feature of the identity center's password policy indicates whether to prohibit the use of the same password as the previous one"
  type        = number
  default     = 3

  validation {
    condition     = var.policy_password_reuse_prevention >= 0 && var.policy_password_reuse_prevention <= 24
    error_message = "The valid value of the password reuse prevention feature of the identity center's password policy is range from 0 to 24"
  }
}

variable "policy_password_not_username_or_invert" {
  description = "Whether the password can be the username or the username spelled backwards"
  type        = bool
  default     = false
}

variable "policy_password_validity_period" {
  description = "The password validity period (days)"
  type        = number
  default     = 7

  validation {
    condition     = var.policy_password_validity_period >= 0 && var.policy_password_validity_period <= 180
    error_message = "The valid value of the password validity period (days) is range from 0 to 180"
  }
}

variable "policy_password_char_combination" {
  description = "The minimum number of character types that a password must contain"
  type        = number

  validation {
    condition     = var.policy_password_char_combination >= 2 && var.policy_password_char_combination <= 4
    error_message = "The valid value of the minimum number of character types that a password must contain is range from 2 to 4"
  }
}

variable "policy_allow_user_to_change_password" {
  description = "Whether IAM users are allowed to change their own passwords"
  type        = bool
  default     = true
}
# SC.005 Enable
