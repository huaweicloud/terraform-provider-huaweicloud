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

# Variable definitions for huaweicloud_das_dead_lock_switch
variable "lock_analysis_instance_id" {
  description = "The ID of the database instances, separated by commas"
  type        = string
}

# Variable definitions for huaweicloud_das_full_dead_lock_switch
variable "full_dead_lock_switch_on" {
  description = "Whether to enable the full dead lock switch"
  type        = bool
  default     = false
}

variable "full_dead_lock_retention_hours" {
  description = "The retention hours of the full dead lock data"
  type        = number
  default     = null
}

# Variable definitions for huaweicloud_das_history_transaction_switch
variable "history_transaction_status" {
  description = "The switch status of the history transaction"
  type        = string
  default     = "Enabled"
}

variable "lock_analysis_datastore_type" {
  description = "The database type"
  type        = string
  default     = "MySQL"
}

# Variable definitions for huaweicloud_das_history_transaction_export_task
variable "history_transaction_bucket_name" {
  description = "The OBS bucket name for exporting history transactions"
  type        = string
}

variable "history_transaction_start_time" {
  description = "The start time of the history transactions to export, in RFC3339 format"
  type        = string
  default     = "2000-06-01T00:00:00+08:00"
}

variable "history_transaction_end_time" {
  description = "The end time of the history transactions to export, in RFC3339 format"
  type        = string
  default     = "2099-06-02T00:00:00+08:00"
}

variable "history_transaction_file_path" {
  description = "The OBS file directory for the export task"
  type        = string
  default     = null
}

variable "history_transaction_time_zone" {
  description = "The time zone for the export task"
  type        = string
  default     = "UTC+8"
}

variable "history_transaction_order_field" {
  description = "The sort field for the export task"
  type        = string
  default     = "collectTime"
}

variable "history_transaction_order_by" {
  description = "The sort order for the export task"
  type        = string
  default     = "asc"
}

variable "history_transaction_last_sec_min" {
  description = "The minimum duration for the export task"
  type        = number
  default     = 0
}

variable "history_transaction_last_sec_max" {
  description = "The maximum duration for the export task"
  type        = number
  default     = 100
}
