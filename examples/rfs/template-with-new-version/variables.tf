# Variable definitions for authentication
variable "region_name" {
  description = "The region where the RFS resources are located"
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
variable "template_name" {
  description = "The name of the RFS template"
  type        = string
}

variable "template_body" {
  description = "The Terraform template body of the RFS template"
  type        = string
}

variable "template_description" {
  description = "The description of the RFS template"
  type        = string
  default     = ""
}

variable "template_initial_version_description" {
  description = "The initial version description of the RFS template"
  type        = string
  default     = ""
}

variable "template_version_body" {
  description = "The Terraform template body of the RFS template version"
  type        = string
}

variable "template_version_description" {
  description = "The description of the RFS template version"
  type        = string
  default     = ""
}
