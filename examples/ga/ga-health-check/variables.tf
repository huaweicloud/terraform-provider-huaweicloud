# Variable definitions for authentication
variable "region_name" {
  description = "The region where GA resources will be created"
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

# Variable definitions for GA accelerator
variable "accelerator_name" {
  description = "The name of the GA accelerator"
  type        = string
}

variable "accelerator_description" {
  description = "The description of the GA accelerator"
  type        = string
  default     = "GA accelerator for health check"
}

variable "ip_area" {
  description = "The area of the IP address. Valid values: CM, CT, CU, EU, AP, AF, ME, GE"
  type        = string
  default     = "CM"
}

# Variable definitions for GA listener
variable "tags" {
  description = "The tags of the GA accelerator and listener"
  type        = map(string)
  default     = {}
}

variable "listener_name" {
  description = "The name of the GA listener"
  type        = string
}

variable "listener_protocol" {
  description = "The protocol of the GA listener. Valid values: TCP, UDP"
  type        = string
  default     = "TCP"
}

variable "listener_description" {
  description = "The description of the GA listener"
  type        = string
  default     = "GA listener for health check"
}

variable "port_from" {
  description = "The start port of the listener port range"
  type        = number
  default     = 4000
}

variable "port_to" {
  description = "The end port of the listener port range"
  type        = number
  default     = 4200
}

# Variable definitions for GA endpoint group
variable "endpoint_group_name" {
  description = "The name of the GA endpoint group"
  type        = string
}

variable "endpoint_group_description" {
  description = "The description of the GA endpoint group"
  type        = string
  default     = "GA endpoint group for health check"
}

variable "backend_region" {
  description = "The region where the backend resources are located"
  type        = string
  default     = "cn-south-1"
}

# Variable definitions for GA health check
variable "health_check_enabled" {
  description = "Whether to enable health check"
  type        = bool
  default     = true
}

variable "health_check_interval" {
  description = "The health check interval in seconds. Range: 1-60"
  type        = number
  default     = 10
}

variable "health_check_max_retries" {
  description = "The maximum number of retries. Range: 1-10"
  type        = number
  default     = 5
}

variable "health_check_port" {
  description = "The port used for health check. Range: 1-65535"
  type        = number
  default     = 8001
}

variable "health_check_timeout" {
  description = "The timeout duration of health check in seconds. Range: 1-60"
  type        = number
  default     = 10
}
