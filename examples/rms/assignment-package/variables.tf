# Variable definitions for authentication
variable "region_name" {
  description = "The region where the resources are located"
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

variable "template_key" {
  description = "The name of a built-in assignment package template"
  type        = string
  default     = ""
}

# Variable definitions for assignment package
variable "assignment_package_name" {
  description = "The name of the assignment package"
  type        = string
}
