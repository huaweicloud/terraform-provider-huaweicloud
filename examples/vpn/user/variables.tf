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

# Variable definitions for resources/data sources
variable "vpn_user_server_id" {
  description = "The VPN server ID"
  type        = string
}

variable "vpn_user_name" {
  description = "The user name"
  type        = string
}

variable "vpn_user_password" {
  description = "The user password"
  type        = string
  sensitive   = true
}
