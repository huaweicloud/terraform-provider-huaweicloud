# Variable definitions for authentication
variable "region_name" {
  description = "The region where the LTS service is located"
  type        = string
}

variable "access_key" {
  description = "The access key of the IAM user"
  type        = string
}

variable "secret_key" {
  description = "The secret key of the IAM user"
  type        = string
}

# Variable definitions for SMN resources
variable "topic_name" {
  description = "The name of the SMN topic"
  type        = string
}

# Variable definitions for LTS resources
variable "group_name" {
  description = "The name of the log group"
  type        = string
}

variable "group_log_expiration_days" {
  description = "The log expiration days of the log group"
  type        = number
  default     = 14
}

variable "stream_name" {
  description = "The name of the log stream"
  type        = string
}

variable "stream_log_expiration_days" {
  description = "The log expiration days of the log stream"
  type        = number
  default     = null
}

variable "notification_template_name" {
  description = "The name of the notification template"
  type        = string
  default     = ""
  nullable    = false
}

variable "domain_id" {
  description = "The domain ID"
  type        = string
  default     = null
}

variable "alarm_rule_name" {
  description = "The name of the SQL alarm rule"
  type        = string
}

variable "alarm_rule_condition_expression" {
  description = "The condition expression of the SQL alarm rule"
  type        = string
}

variable "alarm_rule_alarm_level" {
  description = "The alarm level of the SQL alarm rule"
  type        = string
  default     = "MINOR"
}

variable "alarm_rule_trigger_condition_count" {
  description = "The trigger condition count of the SQL alarm rule"
  type        = number
  default     = 2
}

variable "alarm_rule_trigger_condition_frequency" {
  description = "The trigger condition frequency of the SQL alarm rule"
  type        = number
  default     = 3
}

variable "alarm_rule_send_recovery_notifications" {
  description = "The send recovery notifications of the SQL alarm rule"
  type        = bool
  default     = true
}

variable "alarm_rule_recovery_frequency" {
  description = "The recovery frequency of the SQL alarm rule"
  type        = number
  default     = 4
}

variable "alarm_rule_notification_frequency" {
  description = "The notification frequency of the SQL alarm rule"
  type        = number
  default     = 15
}

variable "alarm_rule_alias" {
  description = "The alias of the SQL alarm rule"
  type        = string
  default     = ""
}

variable "alarm_rule_request_title" {
  description = "The request title of the SQL alarm rule"
  type        = string
}

variable "alarm_rule_request_sql" {
  description = "The request SQL of the SQL alarm rule"
  type        = string
}

variable "alarm_rule_request_search_time_range_unit" {
  description = "The request search time range unit of the SQL alarm rule"
  type        = string
  default     = "minute"
}

variable "alarm_rule_request_search_time_range" {
  description = "The request search time range of the SQL alarm rule"
  type        = number
  default     = 5
}

variable "alarm_rule_frequency_type" {
  description = "The frequency type of the SQL alarm rule"
  type        = string
  default     = "HOURLY"
}

variable "alarm_rule_notification_user_name" {
  description = "The notification user name of the SQL alarm rule"
  type        = string
}

variable "alarm_rule_notification_language" {
  description = "The notification language of the SQL alarm rule"
  type        = string
  default     = "en-us"
}
