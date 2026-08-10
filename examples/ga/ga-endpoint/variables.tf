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
  default     = ""
}

variable "ip_area" {
  description = "The area of the IP address. Valid values: CM, CT, CU, EU, AP, AF, ME, GE"
  type        = string
  default     = "CM"
}

variable "tags" {
  description = "The tags of the GA accelerator and listener"
  type        = map(string)
  default     = {}
}

# Variable definitions for GA listener
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
  default     = "GA listener for endpoint"
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
  default     = "GA endpoint group"
}

variable "backend_region" {
  description = "The region where the backend EIP resource is located"
  type        = string
  default     = "cn-south-1"
}

# Variable definitions for EIP
variable "eip_type" {
  description = "The type of the EIP. Valid values: 5_bgp, 5_sbgp"
  type        = string
  default     = "5_bgp"
}

variable "eip_name" {
  description = "The name of the EIP bandwidth"
  type        = string
}

variable "bandwidth_size" {
  description = "The size of the EIP bandwidth"
  type        = number
  default     = 8
}

# Variable definitions for GA endpoint
variable "endpoint_weight" {
  description = "The weight of the endpoint for traffic distribution. Range: 0-100"
  type        = number
  default     = 10
}
