# Provider Configuration
variable "region_name" {
  description = "The name of the region to deploy resources"
  type        = string
}

variable "access_key" {
  description = "The access key of HuaweiCloud"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "The secret key of HuaweiCloud"
  type        = string
  sensitive   = true
}

# ESW Instance Configuration
variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "192.168.0.0/16"
}

variable "subnet_name" {
  description = "The name of the subnet"
  type        = string
}

variable "subnet_cidr" {
  description = "The CIDR block of the subnet"
  type        = string
  default     = ""
}

variable "subnet_gateway_ip" {
  description = "The gateway IP of the subnet"
  type        = string
  default     = ""
}

variable "esw_instance_name" {
  description = "The name of the instance"
  type        = string
  default     = ""
}

variable "esw_instance_ha_mode" {
  description = "The HA mode of the instance"
  type        = string
  default     = ""
}

variable "esw_instance_description" {
  description = "The description of the instance"
  type        = string
  default     = ""
}

variable "esw_instance_tunnel_ip" {
  description = "The tunnel IP of the instance"
  type        = string
  default     = "192.168.0.192"
}

variable "esw_connection_name" {
  description = "The name of the connection"
  type        = string
  default     = "192.168.0.192"
}

variable "esw_connection_fixed_ips" {
  description = "The fixed IP IP of the connection"
  type        = list(string)
  default     = []
}

variable "esw_connection_segmentation_id" {
  description = "The segmentation ID of the connection"
  type        = number
  default     = 10000
}
