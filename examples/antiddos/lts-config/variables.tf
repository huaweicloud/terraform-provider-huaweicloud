# Variable definitions for authentication
variable "region_name" {
  description = "The region where the Anti-DDoS cloud domain is located"
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

# Variable definitions for module
variable "lts_group_name" {
  description = "The name of the LTS group"
  type        = string
}

variable "lts_ttl_in_days" {
  description = "The log expiration time(days)"
  type        = number
}

variable "enterprise_project_id" {
  description = "The enterprise project ID"
  type        = string
  default     = null
}

variable "lts_stream_name" {
  description = "The name of the LTS stream"
  type        = string
}

variable "lts_is_favorite" {
  description = "Whether to favorite the log stream"
  type        = bool
  default     = false
}
