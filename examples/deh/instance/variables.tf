# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DEH service is located"
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
variable "availability_zone" {
  description = "The availability zone where the dedicated host will be created"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_host_type" {
  description = "The host type of the dedicated host"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_name" {
  description = "The name of the dedicated host instance"
  type        = string
}

variable "instance_auto_placement" {
  description = "Whether to enable auto placement for the dedicated host"
  type        = string
  default     = "on"
}

variable "instance_metadata" {
  description = "The metadata of the dedicated host"
  type        = map(string)
  default     = {}
}

variable "instance_tags" {
  description = "The tags of the dedicated host"
  type        = map(string)
  default     = {}
}

variable "enterprise_project_id" {
  description = "The enterprise project ID of the dedicated host"
  type        = string
  default     = null
}

variable "instance_charging_mode" {
  description = "The charging mode of the dedicated host"
  type        = string
  default     = "prePaid"
}

variable "instance_period_unit" {
  description = "The unit of the billing period of the dedicated host"
  type        = string
  default     = "month"
}

variable "instance_period" {
  description = "The billing period of the dedicated host"
  type        = string
  default     = "1"
}

variable "instance_auto_renew" {
  description = "Whether to enable auto renew for the dedicated host"
  type        = string
  default     = "false"
}
