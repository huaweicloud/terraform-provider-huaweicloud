# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DSC security level is located"
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
variable "security_level_name" {
  description = "The name of the security level"
  type        = string
}

variable "security_level_color_number" {
  description = "The color number of the security level displayed on the console"
  type        = number
  default     = 6
}

variable "security_level_description" {
  description = "The description of the security level"
  type        = string
  default     = ""
}
