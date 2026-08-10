# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DBSS instance is located"
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
  description = "The availability zone to which the DBSS instance belongs"
  type        = string
  default     = ""
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
  nullable    = false
}

variable "subnet_gateway_ip" {
  description = "The gateway IP of the subnet"
  type        = string
  default     = ""
  nullable    = false
}

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "rds_instance_flavor" {
  description = "The flavor of the RDS instance"
  type        = string
  default     = ""
  nullable    = false
}

variable "database_type" {
  description = "The database type of the RDS instance"
  type        = string
  default     = "MySQL"
}

variable "database_version" {
  description = "The database version of the RDS instance"
  type        = string
  default     = "8.0"
}

variable "instance_mode" {
  description = "The mode of the RDS instance"
  type        = string
  default     = "single"
}

variable "instance_group_type" {
  description = "The performance specification"
  type        = string
  default     = "dedicated"
}

variable "instance_flavor_vcpus" {
  description = "The number of vCPUs of the RDS instance flavor"
  type        = number
  default     = 4
}

variable "rds_instance_name" {
  description = "The name of the RDS instance"
  type        = string
}

variable "volume_type" {
  description = "The type of the volume"
  type        = string
  default     = "CLOUDSSD"
}

variable "volume_size" {
  description = "The size of the volume in GB"
  type        = number
  default     = 100
}

variable "dbss_instance_flavor" {
  description = "The flavor ID of the DBSS instance"
  type        = string
  default     = ""
  nullable    = false
}

variable "dbss_instance_name" {
  description = "The name of the DBSS instance"
  type        = string
}

variable "instance_spec_code" {
  description = "The spec code of the DBSS instance"
  type        = string
  default     = "dbss.bypassaudit.low"
}

variable "instance_description" {
  description = "The description of the DBSS instance"
  type        = string
  default     = ""
}

variable "instance_tags" {
  description = "The tags of the DBSS instance"
  type        = map(string)
  default     = {}
}

variable "enterprise_project_id" {
  description = "The enterprise project ID"
  type        = string
  default     = null
}

variable "charging_mode" {
  description = "The charging mode of the DBSS instance"
  type        = string
  default     = "prePaid"
}

variable "period_unit" {
  description = "The period unit of the DBSS instance"
  type        = string
  default     = "month"
}

variable "period" {
  description = "The period of the DBSS instance"
  type        = number
  default     = 1
}

variable "auto_renew" {
  description = "The auto renew of the DBSS instance"
  type        = string
  default     = "false"
}
