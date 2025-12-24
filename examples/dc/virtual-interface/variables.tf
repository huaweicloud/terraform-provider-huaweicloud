# Variable definitions for authentication
variable "region_name" {
  description = "The region where the resources are located"
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

# Variable definitions for network resources
variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "192.168.0.0/16"
}

# Variable definitions for virtual gateway
variable "virtual_gateway_name" {
  description = "The name of the virtual gateway"
  type        = string
}

# Variable definitions for virtual interface
variable "direct_connect_id" {
  description = "The ID of the direct connection associated with the virtual interface"
  type        = string
}

variable "virtual_interface_name" {
  description = "The name of the virtual interface"
  type        = string
}

variable "virtual_interface_description" {
  description = "The description of the virtual interface"
  type        = string
  default     = "Created by Terraform"
}

variable "virtual_interface_type" {
  description = "The type of the virtual interface"
  type        = string
  default     = "private"
}

variable "route_mode" {
  description = "The route mode of the virtual interface"
  type        = string
  default     = "static"
}

variable "vlan" {
  description = "The VLAN for constom side"
  type        = number
}

variable "bandwidth" {
  description = "The ingress bandwidth size of the virtual interface"
  type        = number
}

variable "remote_ep_group" {
  description = "The CIDR list of remote subnets"
  type        = list(string)
}

variable "address_family" {
  description = "The address family type of the virtual interface"
  type        = string
  default     = "ipv4"
}

variable "local_gateway_v4_ip" {
  description = "The IPv4 address of the virtual interface in cloud side"
  type        = string
}

variable "remote_gateway_v4_ip" {
  description = "The IPv4 address of the virtual interface in client side"
  type        = string
}

variable "enable_bfd" {
  description = "Whether to enable the Bidirectional Forwarding Detection (BFD) function"
  type        = bool
  default     = false
}

variable "enable_nqa" {
  description = "Whether to enable the Network Quality Analysis (NQA) function"
  type        = bool
  default     = false
}

variable "virtual_interface_tags" {
  description = "The tags of the virtual interface"
  type        = map(string)
  default     = {
    "Owner" = "terraform"
    "Env"   = "test"
  }
}
