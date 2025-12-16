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
variable "host_id" {
  description = "The host ID for the host protection."
  type        = string
}

variable "protection_version" {
  description = "The protection version enabled by the host"
  type        = string
}

variable "is_wait_host_available" {
  description = "Whether to wait for the host agent status to become online"
  type        = bool
  default     = false
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project to which the host protection belongs"
  type        = string
  default     = null
}
