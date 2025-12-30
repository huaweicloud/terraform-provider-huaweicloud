# Variable definitions for authentication
variable "region_name" {
  description = "The region where the BMS instance is located"
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
variable "availability_zone" {
  description = "The availability zone to which the BMS instance belongs"
  type        = string
  default     = ""
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

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "instance_flavor_id" {
  description = "The flavor ID of the BMS instance"
  type        = string
  default     = ""
}

variable "instance_image_id" {
  description = "The image ID of the BMS instance"
  type        = string
  default     = ""
}

variable "keypair_name" {
  description = "The name of the KPS keypair"
  type        = string
}

variable "instance_name" {
  description = "The name of the BMS instance"
  type        = string
}

variable "instance_user_id" {
  description = "The user ID of the BMS instance"
  type        = string
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project to which the BMS instance belongs"
  type        = string
  default     = null
}

variable "instance_tags" {
  description = "The key/value pairs to associate with the BMS instance"
  type        = map(string)
  default     = {}
}

variable "charging_mode" {
  description = "The charging mode of the BMS instance"
  type        = string
  default     = "prePaid"
}

variable "period_unit" {
  description = "The period unit of the BMS instance"
  type        = string
  default     = "month"
}

variable "period" {
  description = "The period of the BMS instance"
  type        = number
  default     = 1
}

variable "auto_renew" {
  description = "The auto renew of the BMS instance"
  type        = string
  default     = "false"
}
