# Variable definitions for authentication
variable "region_name" {
  description = "The region where CCI service is located"
  type        = string
}

variable "access_key" {
  description = "The access key of IAM user"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "The secret key of IAM user"
  type        = string
  sensitive   = true
}

variable "security_group_name" {
  description = "The name of security group"
  type        = string
  default     = "tf-test-secgroup"
}

variable "vpc_name" {
  description = "The name of VPC"
  type        = string
  default     = "tf-test-vpc"
}

variable "vpc_cidr" {
  description = "The CIDR block of VPC"
  type        = string
  default     = "192.168.0.0/16"
}

variable "subnet_name" {
  description = "The name of subnet"
  type        = string
  default     = "tf-test-subnet"
}

variable "subnet_cidr" {
  description = "The CIDR block of subnet"
  type        = string
  default     = "192.168.0.0/24"
}

variable "subnet_gateway_ip" {
  description = "The gateway IP of subnet"
  type        = string
  default     = "192.168.0.1"
}

variable "namespace_name" {
  description = "The name of CCI namespace"
  type        = string
}

variable "elb_name" {
  description = "The name of ELB load balancer"
  type        = string
}

variable "service_name" {
  description = "The name of CCI service"
  type        = string
}

variable "selector_app" {
  description = "The app label of selector"
  type        = string
  default     = "test1"
}

variable "service_type" {
  description = "The type of service"
  type        = string
  default     = "LoadBalancer"
}
