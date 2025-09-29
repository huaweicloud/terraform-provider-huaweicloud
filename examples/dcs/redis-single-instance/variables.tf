# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DCS Redis single instance is located"
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
  nullable    = false
}

variable "subnet_gateway_ip" {
  description = "The gateway IP address of the subnet"
  type        = string
  default     = ""
  nullable    = false
}

variable "availability_zone" {
  description = "The availability zone to which the Redis single instance belongs"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_flavor_id" {
  description = "The flavor ID of the Redis single instance"
  type        = string
  default     = ""
}

variable "instance_capacity" {
  description = "The capacity of the Redis instance (in GB)"
  type        = number
  default     = 1
}

variable "instance_engine_version" {
  description = "The engine version of the Redis single instance"
  type        = string
  default     = "7.0"
}

variable "instance_name" {
  description = "The name of the Redis single instance"
  type        = string
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project to which the Redis single instance belongs"
  type        = string
  default     = null
}

variable "instance_password" {
  description = "The password for the Redis instance"
  type        = string
  sensitive   = true
  default     = null
}

variable "charging_mode" {
  description = "The charging mode of the Redis instance"
  type        = string
  default     = "postPaid"
}

variable "period_unit" {
  description = "The unit of the period"
  type        = string
  default     = null
}

variable "period" {
  description = "The period of the Redis instance"
  type        = number
  default     = null
}

variable "auto_renew" {
  description = "Whether auto renew is enabled"
  type        = string
  default     = "false"
}
