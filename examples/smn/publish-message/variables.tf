# Variable definitions for authentication
variable "region_name" {
  description = "The region where resources will be created"
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
variable "topic_name" {
  description = "The name of the SMN topic"
  type        = string
}

variable "topic_display_name" {
  description = "The display name of the SMN topic"
  type        = string
  default     = ""
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project"
  type        = string
  default     = null
}

variable "subscription_protocol" {
  description = "The protocol of the subscription"
  type        = string
}

variable "subscription_endpoint" {
  description = "The endpoint of the subscription"
  type        = string
}

variable "subscription_description" {
  description = "The remark for SMN subscriptions"
  type        = string
  default     = ""
}

variable "template_name" {
  description = "The name of the message template"
  type        = string
  default     = ""
  nullable    = false
}

variable "template_protocol" {
  description = "The protocol of the message template"
  type        = string
  default     = ""
  nullable    = false

  validation {
    condition     = var.template_name == "" || var.template_protocol != ""
    error_message = "The template_protocol is required if template_name is provided."
  }
}

variable "template_content" {
  description = "The content of the message template"
  type        = string
  default     = ""
  nullable    = false

  validation {
    condition     = var.template_name == "" || var.template_content != ""
    error_message = "The template_content is required if template_name is provided."
  }
}

variable "pulblish_subject" {
  description = "The subject of the message"
  type        = string
}

variable "pulblish_message" {
  description = "The message content (mutually exclusive with message_structure)"
  type        = string
  default     = ""
  nullable    = false
}

variable "pulblish_message_structure" {
  description = "The JSON message structure that allows sending different content to different protocol subscribers (mutually exclusive with message)"
  type        = string
  default     = ""
  nullable    = false

  validation {
    condition     = !(var.template_name == "" && var.pulblish_message == "") || var.pulblish_message_structure != ""
    error_message = "The pulblish_message_structure is required if both template_name and pulblish_message are not provided."
  }
}

variable "pulblish_time_to_live" {
  description = "The maximum retention time of the message within the SMN system in seconds (default: 3600, max: 86400)"
  type        = string
  default     = null
}

variable "pulblish_tags" {
  description = "The tags of the message"
  type        = map(string)
  default     = {}
}

variable "pulblish_message_attributes" {
  description = "The message attributes of the message"
  type        = list(object({
    name   = string
    type   = string
    value  = optional(string)
    values = optional(list(string))
  }))

  default  = []
  nullable = false
}
