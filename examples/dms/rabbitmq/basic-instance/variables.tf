# Variable definitions for authentication
variable "region_name" {
  description = "The region where the RabbitMQ instance is located"
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
  description = "The availability zones to which the RabbitMQ instance belongs"
  type        = list(string)
  default     = []
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
  default     = "192.168.0.0/24"
  nullable    = false
}

variable "subnet_gateway_ip" {
  description = "The gateway IP of the subnet"
  type        = string
  default     = "192.168.0.1"
  nullable    = false
}

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "instance_flavor_id" {
  description = "The flavor ID of the RabbitMQ instance"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_flavor_type" {
  description = "The flavor type of the RabbitMQ instance"
  type        = string
  default     = "cluster"
}

variable "instance_storage_spec_code" {
  description = "The storage specification code of the RabbitMQ instance"
  type        = string
  default     = "dms.physical.storage.ultra.v2"
}

variable "availability_zone_number" {
  description = "The number of availability zones to which the RabbitMQ instance belongs"
  type        = number
  default     = 1
}

variable "instance_name" {
  description = "The name of the RabbitMQ instance"
  type        = string
}

variable "instance_engine_version" {
  description = "The engine version of the RabbitMQ instance"
  type        = string
  default     = "3.8.35"
}

variable "instance_broker_num" {
  description = "The number of brokers of the RabbitMQ instance"
  type        = number
  default     = 3
}

variable "instance_storage_space" {
  description = "The storage space of the RabbitMQ instance"
  type        = number
  default     = 600
}

variable "instance_ssl_enable" {
  description = "The SSL enable of the RabbitMQ instance"
  type        = bool
  default     = false
}

variable "instance_access_user_name" {
  description = "The access user of the RabbitMQ instance"
  type        = string
}

variable "instance_password" {
  description = "The access password of the RabbitMQ instance"
  sensitive   = true
  type        = string
}

variable "instance_description" {
  description = "The description of the RabbitMQ instance"
  type        = string
  default     = ""
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project to which the RabbitMQ instance belongs"
  type        = string
  default     = null
}

variable "instance_tags" {
  description = "The key/value pairs to associate with the instance"
  type        = map(string)
  default     = {}
}

variable "charging_mode" {
  description = "The charging mode of the RabbitMQ instance"
  type        = string
  default     = "postPaid"
}

variable "period_unit" {
  description = "The period unit of the RabbitMQ instance"
  type        = string
  default     = null
}

variable "period" {
  description = "The period of the RabbitMQ instance"
  type        = number
  default     = null
}

variable "auto_renew" {
  description = "The auto renew of the RabbitMQ instance"
  type        = string
  default     = "false"
}
