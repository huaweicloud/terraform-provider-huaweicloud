# Variable definitions for authentication
variable "region_name" {
  description = "The region where the VOD service is located"
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

# Variable definitions for VOD resources
variable "media_category_name" {
  description = "The name of the media category"
  type        = string
}

variable "template_group_name" {
  description = "The name of the transcoding template group"
  type        = string
}

variable "template_group_description" {
  description = "The description of the transcoding template group"
  type        = string
  default     = ""
}
