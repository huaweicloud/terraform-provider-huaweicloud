# Variable definitions for authentication
variable "region_name" {
  description = "The region where the RFS private module is located"
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
variable "module_name" {
  description = "The name of the RFS private module"
  type        = string
}

variable "module_description" {
  description = "The description of the RFS private module"
  type        = string
  default     = ""
}

variable "module_version" {
  description = "The version number of the RFS private module"
  type        = string
}

variable "module_uri" {
  description = "The OBS address of the private module package"
  type        = string
}

variable "version_description" {
  description = "The description of the private module version"
  type        = string
  default     = ""
}
