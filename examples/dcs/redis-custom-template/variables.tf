# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DCS Redis single instance is located"
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
variable "source_template_id" {
  description = "The ID of the source template to create the custom template from"
  type        = string
}

variable "source_type" {
  description = "The type of the source template. Valid values: sys, user"
  type        = string
  default     = "sys"
}

variable "template_name" {
  description = "The name of the custom template"
  type        = string
}

variable "template_description" {
  description = "The description of the custom template"
  type        = string
  default     = ""
}

variable "template_params" {
  description = "The template params to override, mapping param names to values"
  type        = map(string)
  default     = {
    "timeout" = "200"
  }
}
