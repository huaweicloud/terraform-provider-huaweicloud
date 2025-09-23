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

variable "group_tags" {
  description = "The tags of the log group"
  type        = map(string)
  default     = {}
}

variable "enterprise_project_id" {
  description = "The enterprise project ID of the log group"
  type        = string
  default     = null
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

variable "stream_tags" {
  description = "The tags of the log stream"
  type        = map(string)
  default     = {}
}

variable "stream_is_favorite" {
  description = "Whether to favorite the log stream"
  type        = bool
  default     = false
}
