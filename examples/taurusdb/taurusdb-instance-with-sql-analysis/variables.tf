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
variable "vpc_name" {
  description = "The VPC name"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "192.168.0.0/16"
}

variable "subnet_name" {
  description = "The subnet name"
  type        = string
}

variable "subnet_cidr" {
  description = "The CIDR block of the subnet"
  type        = string
  default     = ""
}

variable "gateway_ip" {
  description = "The gateway IP address of the subnet"
  type        = string
  default     = ""
}

variable "availability_zone_mode" {
  description = "The availability zone mode. Valid values are single, multi"
  type        = string
  default     = "multi"
}

variable "master_availability_zone" {
  description = "The master availability zone of the TaurusDB instance. If not specified, the first available AZ from flavors will be used"
  type        = string
  default     = ""
}

variable "security_group_name" {
  description = "The security group name"
  type        = string
}

variable "instance_db_port" {
  description = "The database port"
  type        = number
  default     = 3306
}

variable "instance_password" {
  description = "The password for the TaurusDB instance"
  type        = string
  default     = ""
  sensitive   = true
}

variable "instance_name" {
  description = "The TaurusDB instance name"
  type        = string
}

variable "instance_flavor_ref" {
  description = "The flavor code of the TaurusDB instance. If not specified, the first available flavor will be used"
  type        = string
  default     = ""
}

variable "instance_mode" {
  description = "The instance mode. Valid values are Cluster, StandSingle"
  type        = string
  default     = "Cluster"
}

variable "read_replicas" {
  description = "The number of read replicas"
  type        = number
  default     = 2
}

variable "enterprise_project_id" {
  description = "The enterprise project ID"
  type        = string
  default     = "0"
}

variable "volume_type" {
  description = "The storage type of the instance. Valid values are DL6, DL5"
  type        = string
  default     = "DL6"
}

variable "time_zone" {
  description = "The time zone of the instance"
  type        = string
  default     = "UTC+08:00"
}

variable "ssl_option" {
  description = "Whether to enable SSL. Valid values are true, false"
  type        = string
  default     = "true"
}

variable "sql_filter_enabled" {
  description = "Whether to enable SQL filter. Required for SQL control rules and auto throttling"
  type        = bool
  default     = true
}

variable "slow_log_show_original_switch" {
  description = "Whether to enable slow log show original switch"
  type        = bool
  default     = true
}

variable "table_name_case_sensitivity" {
  description = "Whether the kernel table name is case sensitive"
  type        = bool
  default     = true
}

variable "multi_tenant_switch" {
  description = "Whether to enable multi-tenancy switch. Valid values are true, false"
  type        = string
  default     = "true"
}

variable "maintain_begin" {
  description = "The start time of the maintenance window in HH:MM format"
  type        = string
  default     = "02:00"
}

variable "maintain_end" {
  description = "The end time of the maintenance window in HH:MM format"
  type        = string
  default     = "06:00"
}

variable "description" {
  description = "The description of the TaurusDB instance"
  type        = string
  default     = ""
}

variable "seconds_level_monitoring_enabled" {
  description = "Whether to enable seconds level monitoring"
  type        = bool
  default     = true
}

variable "seconds_level_monitoring_period" {
  description = "The seconds level collection period. Valid values are 1, 5"
  type        = number
  default     = 5
}

variable "audit_log_enabled" {
  description = "Whether to enable audit log"
  type        = bool
  default     = true
}

variable "audit_log_keep_days" {
  description = "The number of days for storing audit logs"
  type        = number
  default     = 7
}

variable "reserve_audit_logs" {
  description = "Whether to reserve historical audit logs when SQL audit is disabled. Valid values are true, false"
  type        = string
  default     = "true"
}

variable "instance_backup_time_window" {
  description = "The backup time window in HH:MM-HH:MM format"
  type        = string
}

variable "instance_backup_keep_days" {
  description = "The number of days to retain backups"
  type        = number
}

variable "tags" {
  description = "The tags of the TaurusDB instance"
  type        = map(string)
  default     = {}
}

# SQL concurrency control rule variables
variable "sql_control_rule_sql_type" {
  description = "The SQL statement type for the concurrency control rule. Valid values are SELECT, UPDATE, DELETE"
  type        = string
  default     = "SELECT"
}

variable "sql_control_rule_pattern" {
  description = "The concurrency control rule pattern of SQL statements. Keywords are separated by tildes (~)"
  type        = string
  default     = "select~from~t1"
}

variable "sql_control_rule_max_concurrency" {
  description = "The maximum number of concurrent SQL statements for the control rule"
  type        = number
  default     = 20
}

# SQL auto throttling variables
variable "sql_auto_throttling_start_time" {
  description = "The start time of the SQL auto throttling time window in HH:MM format"
  type        = string
  default     = "00:00"
}

variable "sql_auto_throttling_end_time" {
  description = "The end time of the SQL auto throttling time window in HH:MM format"
  type        = string
  default     = "01:00"
}

variable "sql_auto_throttling_condition" {
  description = "The relationship between CPU usage and active sessions conditions. Valid values are and, or"
  type        = string
  default     = "and"
}

variable "sql_auto_throttling_cpu_usage" {
  description = "The CPU usage threshold (percentage) to trigger auto throttling. Value ranges from 70 to 100"
  type        = number
  default     = 70
}

variable "sql_auto_throttling_active_sessions" {
  description = "The active sessions threshold to trigger auto throttling. Value ranges from 1 to 5000"
  type        = number
  default     = 3
}

variable "sql_auto_throttling_clear_time" {
  description = "The maximum throttling duration each time in minutes. Value ranges from 1 to 1440"
  type        = number
  default     = 3
}

variable "sql_auto_throttling_duration" {
  description = "The duration for which the limiting conditions are met in minutes. Value ranges from 2 to 60"
  type        = number
  default     = 2
}

variable "sql_auto_throttling_max_concurrency" {
  description = "The maximum concurrency to limit concurrent SQL executions. Value ranges from 0 to 1000000000"
  type        = number
  default     = 1000
}

variable "sql_auto_throttling_retain_sql_rule" {
  description = "Whether to retain existing SQL limiting rules. Valid values are true, false"
  type        = string
  default     = "true"
}
