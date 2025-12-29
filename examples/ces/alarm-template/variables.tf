# Provider Configuration
variable "region_name" {
  description = "The name of the region to deploy resources"
  type        = string
}

variable "access_key" {
  description = "The access key of HuaweiCloud"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "The secret key of HuaweiCloud"
  type        = string
  sensitive   = true
}

# CES alarm template Configuration
variable "alarm_template_name" {
  description = "The name of the alarm template"
  type        = string
}

variable "alarm_template_description" {
  description = "The description of the alarm template"
  type        = string
  default     = ""
}

variable "alarm_template_policies" {
  description = "The policy list of the CES alarm template"
  type        = list(object({
    namespace           = string
    metric_name         = string
    period              = number
    filter              = string
    comparison_operator = string
    count               = number
    suppress_duration   = number
    value               = number
    alarm_level         = number
    unit                = string
    dimension_name      = string
    hierarchical_value  = list(object({
      critical = number
      major    = number
      minor    = number
      info     = number
    }))
  }))
  default     = []
}
