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
variable "smn_topic_name" {
  description = "The name of the SMN topic that used to send the SMN notification"
  type        = string
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project"
  type        = string
  default     = ""
}

variable "alarm_callback_urls" {
  description = "The URLs of the alarm callback"
  type        = list(string)

  validation {
    condition     = length(var.alarm_callback_urls) > 0 && alltrue([for url in var.alarm_callback_urls : length(regexall("^http[s]?://", url)) > 0])
    error_message = "The alarm callback URLs must be provided and must start with http:// or https://"
  }
}

variable "alarm_notification_template_name" {
  description = "The name of the AOM alarm notification template that used to send the SMN notification"
  type        = string
}

variable "alarm_notification_template_locale" {
  description = "The locale of the AOM alarm notification template that used to send the SMN notification"
  type        = string
  default     = "en-us"

  validation {
    condition     = contains(["en-us", "zh-cn"], var.alarm_notification_template_locale)
    error_message = "The alarm notification template locale must be 'en-us' or 'zh-cn'"
  }
}

variable "alarm_notification_template_description" {
  description = "The description of the AOM alarm notification template that used to send the SMN notification"
  type        = string
  default     = ""
}

variable "alarm_notification_template_notification_type" {
  description = "The notification type of the AOM alarm notification template that used to send the SMN notification"
  type        = string
  default     = "email"
}

variable "alarm_notification_template_notification_topic" {
  description = "The notification topic of the AOM alarm notification template that used to send the SMN notification"
  type        = string
  default     = "An alert occurred at time $${starts_at}[$${event_severity}_$${event_type}_$${clear_type}]."
}

variable "alarm_notification_template_content" {
  description = "The content of the AOM alarm notification template that used to send the SMN notification"
  type        = string
  default     = <<EOT
Alarm Name: $${event_name_alias};
Alarm ID: $${id};
Notification Rule: $${action_rule};
Trigger Time: $${starts_at};
Trigger Level: $${event_severity};
Alarm Content: $${alarm_info};
Resource Identifier: $${resources_new};
Remediation Suggestion: $${alarm_fix_suggestion_zh};
  EOT
}

variable "alarm_action_rule_name" {
  description = "The name of the AOM alarm action rule that used to send the SMN notification"
  type        = string
}

variable "alarm_action_rule_user_name" {
  description = "The user name of the AOM alarm action rule that used to send the SMN notification"
  type        = string
}

variable "alarm_action_rule_type" {
  description = "The type of the AOM alarm action rule that used to send the SMN notification"
  type        = string
  default     = "1" # notification
}
