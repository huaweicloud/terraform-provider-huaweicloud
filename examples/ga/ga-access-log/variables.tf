# Variable definitions for authentication
variable "region_name" {
  description = "The region where GA resources will be created. Currently, GA flow log interconnects with LTS only in cn-north-4"
  type        = string
  default     = "cn-north-4"
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
  default     = ""
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
  default     = ""
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

# Variable definitions for LTS
variable "lts_group_name" {
  description = "The name of the LTS log group"
  type        = string
}

variable "lts_ttl_in_days" {
  description = "The TTL in days for the LTS log group"
  type        = number
  default     = 30
}

variable "lts_stream_name" {
  description = "The name of the LTS log stream"
  type        = string
}
