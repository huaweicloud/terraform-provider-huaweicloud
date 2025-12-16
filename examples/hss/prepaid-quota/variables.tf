# Variable definitions for authentication
variable "region_name" {
  description = "The region where the resources are located"
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
variable "quota_version" {
  description = "The protection quota version"
  type        = string
}

variable "period_unit" {
  description = "The charging period unit of the quota"
  type        = string
}

variable "period" {
  description = "The charging period of the quota"
  type        = number
}

variable "is_auto_renew" {
  description = "Whether auto-renew is enabled"
  type        = bool
  default     = false
}

variable "enterprise_project_id" {
  description = "The enterprise project ID to which the HSS quota belongs"
  type        = string
  default     = null
}

variable "quota_tags" {
  description = "The key/value pairs to associate with the HSS quota"
  type        = map(string)
  default     = null
}
