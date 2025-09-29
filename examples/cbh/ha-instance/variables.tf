# Variable definitions for authentication
variable "region_name" {
  description = "The region where the CBH HA instance is located"
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
variable "master_availability_zone" {
  description = "The availability zone name of the master instance"
  type        = string
  default     = ""
}

variable "slave_availability_zone" {
  description = "The availability zone name of the slave instance"
  type        = string
  default     = ""
}

variable "instance_flavor_id" {
  description = "The flavor ID of the CBH HA instance"
  type        = string
  default     = ""
}

variable "instance_flavor_type" {
  description = "The flavor type of the CBH HA instance"
  type        = string
  default     = "basic"
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
  description = "The name of the CBH HA instance"
  type        = string
}

variable "instance_password" {
  description = "The login password for the CBH HA instance"
  type        = string
  sensitive   = true
}

variable "charging_mode" {
  description = "The charging mode of the CBH HA instance"
  type        = string
  default     = "prePaid"
}

variable "period_unit" {
  description = "The charging period unit of the CBH HA instance"
  type        = string
  default     = "month"
}

variable "period" {
  description = "The charging period of the CBH HA instance"
  type        = number
  default     = 1
}

variable "auto_renew" {
  description = "Whether to enable auto-renew for the CBH HA instance"
  type        = string
  default     = "false"
}
