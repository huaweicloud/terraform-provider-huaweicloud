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
variable "vpc_name" {
  description = "The VPC name"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "192.168.0.0/16"
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

variable "security_group_name" {
  description = "The security group name"
  type        = string
}

variable "instance_db_port" {
  description = "The database port"
  type        = number
  default     = 3306
}

variable "availability_zone_mode" {
  description = "The availability zone mode. Valid values are single, multi"
  type        = string
  default     = "multi"
}

variable "associate_eip_address" {
  description = "The existing EIP address to associate with the TaurusDB instance. If not specified, a new EIP will be created"
  type        = string
  default     = ""
}

variable "eip_type" {
  description = "The EIP type"
  type        = string
  default     = "5_bgp"
}

variable "bandwidth_name" {
  description = "The bandwidth name"
  type        = string
  default     = ""

  validation {
    condition     = var.associate_eip_address != "" || var.bandwidth_name != ""
    error_message = "The bandwidth name must be a non-empty string if the EIP address is not provided."
  }
}

variable "bandwidth_size" {
  description = "The bandwidth size in Mbit/s"
  type        = number
  default     = 5
}

variable "bandwidth_share_type" {
  description = "The share type of the bandwidth. Valid values are PER, WHOLE"
  type        = string
  default     = "PER"
}

variable "bandwidth_charge_mode" {
  description = "The charge mode of the bandwidth. Valid values are traffic, bandwidth"
  type        = string
  default     = "traffic"
}

variable "master_availability_zone" {
  description = "The master availability zone of the TaurusDB instance. If not specified, the first available AZ from flavors will be used"
  type        = string
  default     = ""
}

variable "instance_password" {
  description = "The password for the TaurusDB instance"
  type        = string
  default     = ""
  sensitive   = true
}

variable "instance_name" {
  description = "The TaurusDB instance name"
  type        = string
}

variable "instance_flavor_ref" {
  description = "The flavor code of the TaurusDB instance. If not specified, the first available flavor will be used"
  type        = string
  default     = ""
}

variable "instance_mode" {
  description = "The instance mode. Valid values are Cluster, StandSingle"
  type        = string
  default     = "Cluster"
}

variable "read_replicas" {
  description = "The number of read replicas"
  type        = number
  default     = 2
}

variable "enterprise_project_id" {
  description = "The enterprise project ID"
  type        = string
  default     = "0"
}

variable "volume_type" {
  description = "The storage type of the instance. Valid values are DL6, DL5"
  type        = string
  default     = "DL6"
}

variable "time_zone" {
  description = "The time zone of the instance"
  type        = string
  default     = "UTC+08:00"
}

variable "ssl_option" {
  description = "Whether to enable SSL. Valid values are true, false"
  type        = string
  default     = "true"
}

variable "sql_filter_enabled" {
  description = "Whether to enable SQL filter"
  type        = bool
  default     = true
}

variable "slow_log_show_original_switch" {
  description = "Whether to enable slow log show original switch"
  type        = bool
  default     = true
}

variable "table_name_case_sensitivity" {
  description = "Whether the kernel table name is case sensitive"
  type        = bool
  default     = true
}

variable "multi_tenant_switch" {
  description = "Whether to enable multi-tenancy switch. Valid values are true, false"
  type        = string
  default     = "true"
}

variable "maintain_begin" {
  description = "The start time of the maintenance window in HH:MM format"
  type        = string
  default     = "02:00"
}

variable "maintain_end" {
  description = "The end time of the maintenance window in HH:MM format"
  type        = string
  default     = "06:00"
}

variable "description" {
  description = "The description of the TaurusDB instance"
  type        = string
  default     = ""
}

variable "seconds_level_monitoring_enabled" {
  description = "Whether to enable seconds level monitoring"
  type        = bool
  default     = true
}

variable "seconds_level_monitoring_period" {
  description = "The seconds level collection period. Valid values are 1, 5"
  type        = number
  default     = 5
}

variable "audit_log_enabled" {
  description = "Whether to enable audit log"
  type        = bool
  default     = true
}

variable "audit_log_keep_days" {
  description = "The number of days for storing audit logs"
  type        = number
  default     = 7
}

variable "reserve_audit_logs" {
  description = "Whether to reserve historical audit logs when SQL audit is disabled. Valid values are true, false"
  type        = string
  default     = "true"
}

variable "instance_backup_time_window" {
  description = "The backup time window in HH:MM-HH:MM format"
  type        = string
}

variable "instance_backup_keep_days" {
  description = "The number of days to retain backups"
  type        = number
}

variable "tags" {
  description = "The tags of the TaurusDB instance"
  type        = map(string)
  default     = {}
}
