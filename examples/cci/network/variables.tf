# Variable definitions for authentication
variable "region_name" {
  description = "The region where the CCI network is located"
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

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
  default     = "tf-test-secgroup"
}

variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
  default     = "tf-test-vpc"
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "192.168.0.0/16"
}

variable "subnet_name" {
  description = "The name of the subnet"
  type        = string
  default     = "tf-test-subnet"
}

variable "subnet_cidr" {
  description = "The CIDR block of the subnet"
  type        = string
  default     = "192.168.0.0/24"
}

variable "subnet_gateway_ip" {
  description = "The gateway IP of the subnet"
  type        = string
  default     = "192.168.0.1"
}

variable "namespace_name" {
  description = "The name of the CCI namespace"
  type        = string
}

variable "network_name" {
  description = "The name of the CCI network"
  type        = string
}

variable "warm_pool_size" {
  description = "The size of the warm pool for the network"
  type        = string
  default     = "10"
}

variable "warm_pool_recycle_interval" {
  description = "The recycle interval of the warm pool in hours"
  type        = string
  default     = "2"
}
