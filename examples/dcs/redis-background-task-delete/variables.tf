# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DCS Redis single instance is located"
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
variable "dcs_instance_id" {
  description = "The ID of the DCS instance that owns the background task"
  type        = string
  default     = ""
}

variable "background_task_id" {
  description = "The ID of the background task to delete"
  type        = string
  default     = ""
}
