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

# Variable definitions for SMN topic
variable "smn_topic_name" {
  description = "The name of the SMN topic used to send alarm notifications"
  type        = string
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project"
  type        = string
  default     = null
}

# Variable definitions for CES alarm rule
variable "alarm_rule_name" {
  description = "The name of the CES alarm rule"
  type        = string
}

variable "alarm_rule_description" {
  description = "The description of the CES alarm rule"
  type        = string
  default     = null
}

variable "alarm_action_enabled" {
  description = "Whether to enable the action to be triggered by an alarm"
  type        = bool
  default     = true
}

variable "alarm_enabled" {
  description = "Whether to enable the alarm"
  type        = bool
  default     = true
}

variable "alarm_type" {
  description = "The type of the alarm"
  type        = string
  default     = "ALL_INSTANCE"
}

variable "alarm_rule_conditions" {
  description = "The list of alarm rule conditions"
  type        = list(object({
    metric_name         = string
    period              = string
    filter              = string
    comparison_operator = string
    value               = string
    count               = string
    unit                = optional(string)
    suppress_duration   = optional(string)
    alarm_level         = optional(string)
  }))

  nullable = false
}

variable "alarm_rule_resource" {
  description = "The list of resource dimensions for specified monitoring scope"
  type        = list(object({
    name  = string
    value = optional(string)
  }))

  default  = []
  nullable = true
}

variable "alarm_rule_notification_begin_time" {
  description = "The alarm notification start time"
  type        = string
  default     = null
}

variable "alarm_rule_notification_end_time" {
  description = "The alarm notification stop time"
  type        = string
  default     = null
}

variable "alarm_rule_effective_timezone" {
  description = "The time zone"
  type        = string
  default     = null
}
