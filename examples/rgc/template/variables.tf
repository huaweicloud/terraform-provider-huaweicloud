# Variable definitions for authentication
variable "region_name" {
  description = "The region where the RGC template is located"
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

# Variable definitions for template
variable "template_name" {
  description = "The name of the template"
  type        = string
}

variable "template_type" {
  description = "The type of the template"
  type        = string
  default     = "predefined"
}

variable "template_description" {
  description = "The description of the customized template"
  type        = string
  default     = null
}

variable "template_body" {
  description = "The content of the customized template"
  type        = string
  default     = null
}
