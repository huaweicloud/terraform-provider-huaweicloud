# Variable definitions for authentication
variable "region_name" {
  description = "The region where the SMS task is located"
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
variable "source_server_name" {
  description = "The name of the SMS source server"
  type        = string
  default     = null
}

variable "server_template_name" {
  description = "The name of the SMS server template"
  type        = string
}

variable "migrate_task_type" {
  description = "The type of the SMS task"
  type        = string
}

variable "server_os_type" {
  description = "The OS type of the server"
  type        = string
}
