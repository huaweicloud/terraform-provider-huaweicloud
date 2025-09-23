# Variable definitions for authentication
variable "region_name" {
  description = "The region where the VPC is located"
  type        = string
}

variable "access_key" {
  description = "The access key of the IAM user"
  type        = string
}

variable "secret_key" {
  description = "The secret key of the IAM user"
  type        = string
}

# Variable definitions for resources/data sources
variable "vpc_name" {
  description = "The VPC name"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "172.16.0.0/16"
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project to which the VPC belongs"
  type        = string
  default     = null
}

variable "subnet_name" {
  description = "The subnet name"
  type        = string
}

variable "subnet_cidr" {
  description = "The CIDR block of the subnet"
  type        = string
  default     = "172.16.10.0/24"
}

variable "subnet_gateway" {
  description = "The gateway IP address of the subnet"
  type        = string
  default     = "172.16.10.1"
}

variable "dns_list" {
  description = "The list of DNS server IP addresses"
  type        = list(string)
  default     = null
}
