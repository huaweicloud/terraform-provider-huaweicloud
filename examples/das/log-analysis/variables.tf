# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DAS resources are located"
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

# Variable definitions for huaweicloud_das_slow_log_export_task
variable "log_analysis_instance_id" {
  description = "The ID of the database instance"
  type        = string
  default     = ""
}

variable "slow_log_bucket_name" {
  description = "The OBS bucket name for exporting slow logs"
  type        = string
}

variable "slow_log_start_time" {
  description = "The start time of the slow logs to export, in RFC3339 format"
  type        = string
}

variable "slow_log_end_time" {
  description = "The end time of the slow logs to export, in RFC3339 format"
  type        = string
}

variable "slow_log_file_path" {
  description = "The OBS file directory for the export task"
  type        = string
  default     = null
}

variable "slow_log_export_type" {
  description = "The export type for the slow log export task"
  type        = string
  default     = null
}

variable "slow_log_sort_field" {
  description = "The sort field for the slow log export task"
  type        = string
  default     = null
}

variable "slow_log_sort_asc" {
  description = "Whether to sort in ascending order"
  type        = bool
  default     = null
}

variable "slow_log_time_zone" {
  description = "The time zone for the slow log export task"
  type        = string
  default     = null
}

variable "binlog_binlog_type" {
  description = "The binlog type"
  type        = string
}

variable "binlog_file_name" {
  description = "The binlog file name"
  type        = string
}

variable "binlog_backup_id" {
  description = "The backup ID"
  type        = string
  default     = null
}

# Variable definitions for huaweicloud_das_binlog_parse_task_export
variable "binlog_export_bucket_name" {
  description = "The OBS bucket name for exporting binlog parse results"
  type        = string
}

variable "binlog_filter_db_names" {
  description = "The list of database names to filter"
  type        = list(string)
  default     = []
  nullable    = false
}

variable "binlog_filter_tb_names" {
  description = "The list of table names to filter"
  type        = list(string)
  default     = []
  nullable    = false
}

variable "binlog_filter_start_time" {
  description = "The start time of the export range, in RFC3339 format"
  type        = string
  default     = null
}

variable "binlog_filter_end_time" {
  description = "The end time of the export range, in RFC3339 format"
  type        = string
  default     = null
}

variable "binlog_filter_types" {
  description = "The list of SQL types to filter"
  type        = list(string)
  default     = []
  nullable    = false
}

variable "binlog_filter_parse_double_insert" {
  description = "Whether to export UPDATE statements as two INSERT statements"
  type        = bool
  default     = null
}
