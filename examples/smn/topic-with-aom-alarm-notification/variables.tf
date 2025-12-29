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
  description = "The name of the SMN topic used to send AOM alarm notifications"
  type        = string
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project"
  type        = string
  default     = null
}

variable "lts_group_name" {
  description = "The name of the LTS group"
  type        = string
}

variable "lts_group_ttl_in_days" {
  description = "The TTL in days of the LTS group"
  type        = number
  default     = 30
}

variable "lts_stream_name" {
  description = "The name of the LTS stream"
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

variable "alarm_action_rule_description" {
  description = "The description of the AOM alarm action rule"
  type        = string
  default     = null
}
