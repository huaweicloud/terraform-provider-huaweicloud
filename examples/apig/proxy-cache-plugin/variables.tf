# Variable definitions for authentication
variable "region_name" {
  description = "The region where resources will be created"
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

# Variable definitions for resources/data sources
variable "availability_zones" {
  description = "The availability zones to which the instance belongs"
  type        = list(string)
  default     = []
  nullable    = false
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
  description = "The gateway IP address of the subnet"
  type        = string
  default     = ""
}

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "instance_name" {
  description = "The name of the APIG instance"
  type        = string
}

variable "instance_edition" {
  description = "The edition of the APIG instance"
  type        = string
  default     = "BASIC"
}

variable "availability_zones_count" {
  description = "The number of availability zones to which the instance belongs"
  type        = number
  default     = 1
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project"
  type        = string
  default     = null
}

variable "plugin_name" {
  description = "The name of the APIG plugin"
  type        = string
}

variable "plugin_description" {
  description = "The description of the APIG plugin"
  type        = string
  default     = null
}
