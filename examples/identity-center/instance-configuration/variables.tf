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

variable "configuration_type" {
  description = "The type of the identity center instance configuration"
  type        = string
  default     = "APP_AUTHENTICATION_CONFIGURATION"
}

variable "configuration_mfa_mode" {
  description = "The mfa mode of the identity center instance configuration"
  type        = string
  default     = null

  validation {
    condition     = contains(["ALWAYS_ON", "CONTEXT_AWARE", "DISABLED"], var.configuration_mfa_mode)
    error_message = "The valid values of the mfa mode are: ALWAYS_ON, CONTEXT_AWARE, DISABLED"
  }
}

variable "configuration_allowed_mfa_types" {
  description = "The allowed mfa types of the identity center instance configuration"
  type        = list(string)
  default     = null

  validation {
    condition     = contains(["TOTP", "WEBAUTHN_SECURITY_KEY"], var.configuration_allowed_mfa_types)
    error_message = "The valid values of the allowed mfa types are: TOTP, WEBAUTHN_SECURITY_KEY"
  }
}

variable "configuration_no_mfa_signin_behavior" {
  description = "The no mfa signin behavior of the identity center instance configuration"
  type        = string
  default     = null

  validation {
    condition     = contains(["ALLOWED_WITH_ENROLLMENT", "ALLOWED", "EMAIL_OTP", "BLOCKED"], var.configuration_no_mfa_signin_behavior)
    error_message = "The valid values of the no mfa signin behavior are: ALLOWED_WITH_ENROLLMENT, ALLOWED, EMAIL_OTP, BLOCKED"
  }
}

# SC.005 Disable
variable "configuration_no_password_signin_behavior" {
  description = "The no password signin behavior of the identity center instance configuration"
  type        = string
  default     = null

  validation {
    condition     = contains(["BLOCKED"], var.configuration_no_password_signin_behavior)
    error_message = "The valid value of the no password signin behavior is: BLOCKED"
  }
}
# SC.005 Enable

variable "configuration_max_authentication_age" {
  description = "The max authentication age of the identity center instance configuration"
  type        = string
  default     = null
}
