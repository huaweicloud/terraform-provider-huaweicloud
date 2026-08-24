# Variable definitions for authentication
variable "region_name" {
  description = "The region where the ServiceStage service is located"
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

# Variable definitions for ServiceStage resources
variable "configuration_group_name" {
  description = "The name of the configuration group"
  type        = string
}

variable "configuration_group_description" {
  description = "The description of the configuration group"
  type        = string
  default     = ""
}

variable "configuration_name" {
  description = "The name of the configuration file"
  type        = string
}

variable "configuration_content" {
  description = "The content of the configuration file"
  type        = string
}

variable "configuration_description" {
  description = "The description of the configuration file"
  type        = string
  default     = ""
}
