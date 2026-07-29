# Variable definitions for authentication
variable "region_name" {
  description = "The region where the GaussDB instance is located"
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

# Variable definitions for resources
variable "lts_group_name" {
  description = "The LTS log group name"
  type        = string
}

variable "lts_stream_name" {
  description = "The LTS log stream name"
  type        = string
}

variable "gaussdb_instance_id" {
  description = "The ID of the GaussDB instance"
  type        = string
}

variable "lts_log_type" {
  description = "The LTS log type"
  type        = string
  default     = "audit_log"
}
