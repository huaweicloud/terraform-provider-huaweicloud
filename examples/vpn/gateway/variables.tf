# Variable definitions for authentication
variable "region_name" {
  description = "The region where the VPN service is located"
  type        = string
}

variable "access_key" {
  description = "The access key of the IAM user"
  type        = string
}

variable "vpn_gateway_flavor" {
  description = "The flavor of the VPN gateway"
  type        = string
  default     = "professional1"
}

variable "vpn_gateway_attachment_type" {
  description = "The attachment type of the VPC"
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

variable "eip_type" {
  description = "The EIP type"
  type        = string
  default     = "5_bgp"
}

variable "bandwidth_name" {
  description = "The bandwidth name"
  type        = string
}

variable "bandwidth_size" {
  description = "The bandwidth size"
  type        = number
  default     = 8
}

variable "bandwidth_share_type" {
  description = "The bandwidth share type"
  type        = string
  default     = "PER"
}

variable "bandwidth_charge_mode" {
  description = "The bandwidth charge mode"
  type        = string
  default     = "traffic"
}

variable "vpn_gateway_name" {
  description = "The name of the VPN gateway"
  type        = string
  default     = "traffic"
}

variable "vpn_gateway_delete_eip_on_termination" {
  description = "Whether to delete the EIP when the VPN gateway is deleted."
  type        = bool
  default     = false
}
