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
variable "lts_group_name" {
  description = "The name of the LTS group that used to store the SMN notification logs"
  type        = string
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project"
  type        = string
  default     = ""
}

variable "lts_stream_name" {
  description = "The name of the LTS stream that used to store the SMN notification logs"
  type        = string
}

variable "smn_topic_name" {
  description = "The name of the SMN topic that used to send the SMN notification"
  type        = string
}

variable "alarm_action_rule_name" {
  description = "The name of the AOM alarm action rule that used to send the SMN notification"
  type        = string
  default     = "apm"
}

variable "alarm_action_rule_user_name" {
  description = "The user name of the AOM alarm action rule"
  type        = string
}

variable "alarm_action_rule_type" {
  description = "The type of the AOM alarm action rule"
  type        = string
  default     = "1" # notification
}

variable "alarm_group_rule_name" {
  description = "The name of the AOM alarm group rule"
  type        = string
}

variable "alarm_group_rule_group_interval" {
  description = "The group interval of the AOM alarm group rule"
  type        = number
  default     = 60
}

variable "alarm_group_rule_group_repeat_waiting" {
  description = "The group repeat waiting of the AOM alarm group rule"
  type        = number
  default     = 3600
}

variable "alarm_group_rule_group_wait" {
  description = "The group wait of the AOM alarm group rule"
  type        = number
  default     = 15
}

variable "alarm_group_rule_description" {
  description = "The description of the AOM alarm group rule"
  type        = string
  default     = ""
}

variable "alarm_group_rule_condition_matching_rules" {
  description = "The condition matching rules of the AOM alarm group rule"
  type        = list(object({
    key     = string
    operate = string
    value   = list(string)
  }))
  default     = [
    {
      key     = "event_severity"
      operate = "EXIST"
      value   = ["Critical", "Major"]
    },
    {
      key     = "resource_provider"
      operate = "EQUALS"
      value   = ["AOM"]
    }
  ]
}

variable "alarm_rule_name" {
  description = "The name of the AOM alarm rule"
  type        = string
}

variable "prometheus_instance_id" {
  description = "The ID of the Prometheus instance"
  type        = string
  default     = "0" # The default prometheus instance is 'Prometheus_AOM_Default'.
}

variable "alarm_rule_trigger_conditions" {
  description = "The trigger conditions of the AOM alarm rule"
  type        = list(object({
    metric_name             = string
    promql                  = string
    promql_for              = string
    aggregate_type          = optional(string, "by")
    aggregation_type        = string
    aggregation_window      = string
    metric_statistic_method = string
    thresholds              = map(any) # key is the alarm level, value is the alarm threshold, e.g. "{\"Critical\": 1}"
    trigger_type            = string
    trigger_interval        = string
    trigger_times           = string
    query_param             = string # Query parameters in JSON format
    query_match             = string # Query match conditions in JSON format
  }))
}
