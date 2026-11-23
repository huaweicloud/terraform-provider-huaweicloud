# Variables definitions for authorization
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

# Variables definitions for resource/data source
variable "bandwidth_name" {
  description = "The name of the bandwidth"
  type        = string
}

variable "bandwidth_size" {
  description = "The size of the bandwidth"
  type        = number
  default     = 5
}

variable "bandwidth_type" {
  description = "The type of the bandwidth"
  type        = string
  default     = "share"
}

variable "enterprise_project_id" {
  description = "The enterprise project ID"
  type        = string
  default     = null
}

variable "charging_mode" {
  description = "The charging mode of the bandwidth"
  type        = string
  default     = "postPaid"
}

variable "period_unit" {
  description = "The period unit of the bandwidth"
  type        = string
  default     = null
}

variable "period" {
  description = "The period of the bandwidth"
  type        = number
  default     = null
}

variable "auto_renew" {
  description = "Whether the auto renew is enabled"
  type        = string
  default     = "false"
}

variable "alarm_rule_name" {
  description = "The name of the CES alarm rule"
  type        = string
}

variable "alarm_rule_description" {
  description = "The description of the CES alarm rule"
  type        = string
  default     = ""
}

variable "alarm_enabled" {
  description = "Whether the CES alarm rule is enabled"
  type        = bool
  default     = true
}

variable "alarm_action_enabled" {
  description = "Whether to enable the action to be triggered by an alarm"
  type        = bool
  default     = true
}

variable "rule_conditions" {
  description = "The conditions of the alarm rule"
  type        = list(object({
    period              = number
    filter              = string
    comparison_operator = string
    value               = number
    count               = number
    unit                = optional(string, null)
    alarm_level         = optional(number, 2)
    suppress_duration   = optional(number, 0)
  }))

  nullable = false

  validation {
    condition     = length(var.rule_conditions) > 0
    error_message = "The 'rule_conditions' is not allowed to be empty"
  }
}

variable "bandwidth_policy_name" {
  description = "The name of the AS bandwidth policy"
  type        = string
}

variable "bandwidth_policy_type" {
  description = "The type of the AS bandwidth policy"
  type        = string
  default     = "INTERVAL_ALARM"
}

variable "bandwidth_policy_description" {
  description = "The description of the AS bandwidth policy"
  type        = string
  default     = ""
}

variable "cool_down_time" {
  description = "The cooldown period of the AS bandwidth policy"
  type        = number
  default     = 300
}

variable "interval_alarm_conditions" {
  description = "The alarm interval of the bandwidth policy"
  type        = list(object({
    lower_bound = optional(string, null)
    upper_bound = optional(string, null)
    operation   = optional(string, null)
    size        = optional(number, 2)
    limits      = optional(number, null)
  }))

  nullable = false

  validation {
    condition     = length(var.interval_alarm_conditions) > 0
    error_message = "The 'interval_alarm_conditions' is not allowed to be empty"
  }
}
