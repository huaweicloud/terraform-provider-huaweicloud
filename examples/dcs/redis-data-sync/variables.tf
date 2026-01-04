# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DCS instances are located"
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

# Variable definitions for VPC resources
variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
  default     = "dcs-sync-vpc"
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "192.168.0.0/16"
}

variable "subnet_name" {
  description = "The name of the subnet"
  type        = string
  default     = "dcs-sync-subnet"
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
  default     = "dcs-sync-sg"
}

# Variable definitions for DCS instances
variable "instance_cache_mode" {
  description = "The cache mode of the DCS instances"
  type        = string
  default     = "ha"
}

variable "instance_capacity" {
  description = "The capacity of the DCS instances (GB)"
  type        = number
  default     = 4
}

variable "instance_engine_version" {
  description = "The engine version of the DCS instances"
  type        = string
  default     = "5.0"
}

variable "instance_name" {
  description = "The base name of the DCS Redis instances (will be suffixed with -0 and -1)"
  type        = string
}

variable "instance_password" {
  description = "The password of the DCS instances"
  type        = string
  sensitive   = true
}

# Variable definitions for migration tasks
variable "full_migration_task_name" {
  description = "The name of the full migration task"
  type        = string
  default     = "full-migration-task"
}

variable "full_migration_task_description" {
  description = "The description of the full migration task"
  type        = string
  default     = "Full data migration from source to target DCS instance"
}

variable "full_migration_resume_mode" {
  description = "The reconnection mode for full migration"
  type        = string
  default     = "auto"
}

variable "full_migration_bandwidth_limit_mb" {
  description = "The bandwidth limit for full migration (MB/s)"
  type        = string
  default     = ""
}

variable "enable_incremental_migration" {
  description = "Whether to enable incremental migration after full migration"
  type        = bool
  default     = true
}

variable "incremental_migration_task_name" {
  description = "The name of the incremental migration task"
  type        = string
  default     = "incremental-migration-task"
}

variable "incremental_migration_task_description" {
  description = "The description of the incremental migration task"
  type        = string
  default     = "Incremental data migration from source to target DCS instance"
}

variable "incremental_migration_resume_mode" {
  description = "The reconnection mode for incremental migration"
  type        = string
  default     = "auto"
}

variable "incremental_migration_bandwidth_limit_mb" {
  description = "The bandwidth limit for incremental migration (MB/s)"
  type        = string
  default     = ""
}
