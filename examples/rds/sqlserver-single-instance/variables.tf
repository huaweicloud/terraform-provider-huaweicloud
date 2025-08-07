# Variable definitions for authentication
variable "region_name" {
  description = "The region where the SQLServer RDS instance is located"
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
  default     = "192.168.0.0/16"
}

variable "availability_zone" {
  description = "The availability zone to which the RDS instance belongs"
  type        = string
  default     = ""
}

variable "subnet_name" {
  description = "The subnet name"
  type        = string
}

variable "subnet_cidr" {
  description = "The CIDR block of the subnet"
  type        = string
  default     = ""
}

variable "gateway_ip" {
  description = "The gateway IP address of the subnet"
  type        = string
  default     = ""
}

variable "instance_flavor_id" {
  description = "The flavor ID of the RDS instance"
  type        = string
  default     = ""
}

variable "instance_db_type" {
  description = "The database engine type"
  type        = string
  default     = "SQLServer"
}

variable "instance_db_version" {
  description = "The database engine version"
  type        = string
  default     = "2019_SE"
}

variable "instance_mode" {
  description = "The instance mode for the RDS instance flavor"
  type        = string
  default     = "single"
}

variable "instance_flavor_group_type" {
  description = "The group type for the RDS instance flavor"
  type        = string
  default     = "general"
}

variable "instance_flavor_vcpus" {
  description = "The number of the RDS instance CPU cores for the RDS instance flavor"
  type        = number
  default     = 2
}

variable "instance_flavor_memory" {
  description = "The memory size for the RDS instance flavor"
  type        = number
  default     = 4
}

variable "security_group_name" {
  description = "The security group name"
  type        = string
}

variable "instance_password" {
  description = "The password for the RDS instance"
  type        = string
  default     = ""
}

variable "instance_name" {
  description = "The SQLServer RDS instance name"
  type        = string
}

variable "instance_db_port" {
  description = "The database port"
  type        = number
  default     = 1433
}

variable "instance_volume_type" {
  description = "The storage volume type"
  type        = string
  default     = "CLOUDSSD"
}

variable "instance_volume_size" {
  description = "The storage volume size in GB"
  type        = number
  default     = 40
}

variable "instance_backup_time_window" {
  description = "The backup time window in HH:MM-HH:MM format"
  type        = string
  default     = "03:00-04:00"
}

variable "instance_backup_keep_days" {
  description = "The number of days to retain backups"
  type        = number
  default     = 7
}
