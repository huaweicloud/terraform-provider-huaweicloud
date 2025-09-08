# Variable definitions for authentication
variable "region_name" {
  description = "The region where the Workspace service is located"
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

variable "bucket_name" {
  description = "The name of the OBS bucket"
  type        = string
}

variable "transfer_type" {
  description = "The type of the log transfer"
  type        = string
  default     = "OBS"
}

variable "transfer_mode" {
  description = "The mode of the log transfer"
  type        = string
  default     = "cycle"
}

variable "transfer_storage_format" {
  description = "The storage format of the log transfer"
  type        = string
  default     = "JSON"
}

variable "transfer_status" {
  description = "The status of the log transfer"
  type        = string
  default     = "ENABLE"
}

variable "bucket_dir_prefix_name" {
  description = "The prefix path of the OBS transfer task"
  type        = string
  default     = "LTS-test/%GroupName/%StreamName/%Y/%m/%d/%H/%M"
}

variable "bucket_time_zone" {
  description = "The time zone of the OBS bucket"
  type        = string
  default     = "UTC"
}

variable "bucket_time_zone_id" {
  description = "The time zone ID of the OBS bucket"
  type        = string
  default     = "Etc/GMT"
}
