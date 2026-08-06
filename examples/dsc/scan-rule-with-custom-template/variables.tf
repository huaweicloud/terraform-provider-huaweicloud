# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DSC scan rule is located"
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

variable "scan_template_name" {
  description = "The name of the scan template"
  type        = string
}

variable "scan_template_description" {
  description = "The description of the scan template"
  type        = string
  default     = "Created_by_terraform_script"
}

variable "classification_name" {
  description = "The name of the scan template classification"
  type        = string
}

variable "scan_rule_name" {
  description = "The name of the scan rule"
  type        = string
}

variable "scan_rule_match_rate" {
  description = "The match rate of the scan rule"
  type        = number
  default     = 1
}

variable "scan_rule_description" {
  description = "The description of the scan rule"
  type        = string
  default     = ""
}
