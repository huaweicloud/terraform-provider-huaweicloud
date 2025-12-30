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
variable "vpn_gateway_az_flavor" {
  description = "The flavor name of vpn gateway az"
  type        = string
  default     = "professional1"
}

variable "vpn_gateway_az_attachment_type" {
  description = "The attachment type of vpn gateway az"
  type        = string
  default     = "vpc"
}

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

variable "vpc_eip_public_ip" {
  description = "The public IP of the EIP"
  type        = list(object({
    type = string
  }))
  default     = []
}

variable "vpc_eip_bandwidth" {
  description = "The bandwidth IP of the EIP"
  type        = list(object({
    name        = string
    size        = number
    share_type  = string
    charge_mode = string
  }))
  default     = []
}

variable "vpn_gateway_name" {
  description = "The name of the VPN gateway"
  type        = string
}

variable "vpn_customer_gateway_name" {
  description = "The customer gateway name"
  type        = string
}

variable "vpn_customer_gateway_id_value" {
  description = "The identifier of a customer gateway"
  type        = string
}

variable "vpn_connection_name" {
  description = "The name of the VPN connection"
  type        = string
}

variable "vpn_connection_peer_subnets" {
  description = "The list of customer subnets"
  type        = list(string)
}

variable "vpn_connection_vpn_type" {
  description = "The connection type"
  type        = string
}

variable "vpn_connection_psk" {
  description = "The pre-shared key"
  type        = string
}

variable "vpn_connection_enable_nqa" {
  description = "Whether to enable NQA check"
  type        = bool
}
