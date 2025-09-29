# Variable definitions for authentication
variable "region_name" {
  description = "The region where the VPN service is located"
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

variable "er_instance_name" {
  description = "The ER instance namee"
  type        = string
}

variable "er_instance_asn" {
  description = "The ER instance asn"
  type        = number
}

variable "er_vpc_attachment_name" {
  description = "The bandwidth size"
  type        = string
}

variable "er_vpc_attachment_auto_create_vpc_routes" {
  description = "Whether to enable auto create VPC routes"
  type        = bool
  default     = true
}
