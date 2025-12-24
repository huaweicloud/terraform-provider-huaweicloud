# Variable definitions for authentication
variable "region_name" {
  description = "The region where the volume is located"
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
variable "volume_availability_zone" {
  description = "The availability zone for the volume"
  type        = string
  default     = ""
  nullable    = false
}

variable "volume_name" {
  description = "The name of the volume"
  type        = string
}

variable "volume_type" {
  description = "The type of the volume"
  type        = string
  default     = "SSD"
}

variable "voulme_size" {
  description = "The size of the volume"
  type        = number
  default     = 40
}

variable "volume_description" {
  description = "The description of the volume"
  type        = string
  default     = ""
}

variable "volume_multiattach" {
  description = "The volume is shared volume or not"
  type        = bool
  default     = false
}

variable "volume_iops" {
  description = "The IOPS for the volume"
  type        = number
  default     = null
}

variable "volume_throughput" {
  description = "The throughput for the volume"
  type        = number
  default     = null
}

variable "volume_device_type" {
  description = "The device type of disk to create"
  type        = string
  default     = "VBD"
}

variable "enterprise_project_id" {
  description = "The enterprise project ID of the volume"
  type        = string
  default     = null
}

variable "volume_tags" {
  description = "The tags of the volume"
  type        = map(string)
  default     = {}
}

variable "charging_mode" {
  description = "The charging mode of the volume"
  type        = string
  default     = "postPaid"
}

variable "period_unit" {
  description = "The period unit of the volume"
  type        = string
  default     = null
}

variable "period" {
  description = "The period of the volume"
  type        = number
  default     = null
}

variable "auto_renew" {
  description = "The auto renew of the volume"
  type        = string
  default     = "false"
}
