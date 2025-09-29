# Variable definitions for authentication
variable "region_name" {
  description = "The region where the PostgreSQL RDS instance is located"
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
  description = "The VPC name"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "192.168.0.0/16"
}

variable "availability_zones" {
  description = "The list of availability zones to which the RDS instance belong"
  type        = list(string)
  default     = []
  nullable    = false

  validation {
    condition     = var.instance_mode == "ha" && (length(var.availability_zones) == 0 || length(var.availability_zones) > 1) || var.instance_mode != "ha" && length(var.availability_zones) <= 1
    error_message = "The availability zones must be a list of strings"
  }
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
  default     = "PostgreSQL"
}

variable "instance_db_version" {
  description = "The database engine version"
  type        = string
  default     = "16"
}

variable "instance_mode" {
  description = "The instance mode for the RDS instance flavor"
  type        = string
  default     = "ha"
}

variable "instance_flavor_group_type" {
  description = "The group type for the RDS instance flavor"
  type        = string
  default     = "general"
}

variable "instance_flavor_vcpus" {
  description = "The CPU core numbers for the RDS instance flavor"
  type        = number
  default     = 4
}

variable "instance_flavor_memory" {
  description = "The memory size for the RDS instance flavor"
  type        = number
  default     = 8
}

variable "security_group_name" {
  description = "The security group name"
  type        = string
}

variable "instance_db_port" {
  description = "The database port"
  type        = number
  default     = 5432
}

variable "instance_password" {
  description = "The password for the RDS instance"
  type        = string
  default     = ""
  sensitive   = true
}

variable "instance_name" {
  description = "The name of the RDS instance"
  type        = string
}

variable "ha_replication_mode" {
  description = "The HA replication mode of the RDS instance"
  type        = string
  default     = "async"
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
}

variable "instance_backup_keep_days" {
  description = "The number of days to retain backups"
  type        = number
}

variable "account_name" {
  description = "Username with elevated privileges"
  type        = string
}

variable "account_password" {
  description = "The password for the database account"
  type        = string
  default     = ""
  sensitive   = true
}

variable "database_name" {
  description = "The name of the initial database"
  type        = string
}

variable "schema_name" {
  description = "The name of the database schema"
  type        = string
}

variable "backup_name" {
  description = "The name for instance backups"
  type        = string
}
