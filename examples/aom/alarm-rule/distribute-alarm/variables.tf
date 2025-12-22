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
variable "dcs_instance_name" {
  description = "The name of the DCS instance that used to store the alarm data"
  type        = string
  default     = ""
}

variable "lts_group_name" {
  description = "The name of the LTS group that used to store the SMN notification logs"
  type        = string
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

variable "alarm_rule_matric_dimension_tags" {
  description = "The custom tag of the DCS instance that used to distribute the alarm rule"
  type        = map(string)
}

variable "prometheus_instance_name" {
  description = "The name of the Prometheus instance that used to store the alarm data"
  type        = string
}

variable "alarm_rule_name" {
  description = "The name of the AOM alarm rule"
  type        = string
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
    metric_unit             = string
    metric_query_mode       = string
    metric_namespace        = string
    operator                = string
    metric_statistic_method = string
    thresholds              = map(any) # key is the alarm level, value is the alarm threshold, e.g. "{\"Critical\": 1}"
    trigger_type            = string
    trigger_interval        = string
    trigger_times           = string
    query_param             = string # Query parameters in JSON format
    query_match             = string # Query match conditions in JSON format
  }))
}
