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
variable "server_id" {
  description = "The BMS instance ID"
  type        = string
}

variable "volume_id" {
  description = "The ID of the disk to be attached to a BMS instance"
  type        = string
}

variable "device" {
  description = "The mount point"
  type        = string
  default     = null
}
