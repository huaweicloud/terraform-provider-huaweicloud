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

variable "vcpus" {
  description = "The number of vCPUs of the flavor to query"
  type        = number
  default     = 4
}

variable "availability_zone" {
  description = "The availability zone of the GeminiDB InfluxDB instance. If not specified, the first AZ from data source will be used"
  type        = string
  default     = ""
}

variable "security_group_name" {
  description = "The security group name"
  type        = string
}

variable "instance_password" {
  description = "The password for the GeminiDB InfluxDB instance"
  type        = string
  default     = ""
  sensitive   = true
}

variable "instance_name" {
  description = "The GeminiDB InfluxDB instance name"
  type        = string
}

variable "instance_mode" {
  description = "The instance mode. Valid values are Cluster, EnhancedCluster, InfluxdbSingle"
  type        = string
  default     = "InfluxdbSingle"
}

variable "enterprise_project_id" {
  description = "The enterprise project ID"
  type        = string
  default     = "0"
}

variable "instance_ssl_option" {
  description = "Whether to enable SSL. Valid values are on, off"
  type        = string
  default     = "on"
}

variable "maintenance_start_time" {
  description = "The start time of the maintenance window in UTC HH:MM format"
  type        = string
  default     = "02:00"
}

variable "maintenance_end_time" {
  description = "The end time of the maintenance window in UTC HH:MM format, 4 hours after start time"
  type        = string
  default     = "06:00"
}

variable "instance_flavor_num" {
  description = "The node quantity of the GeminiDB InfluxDB instance"
  type        = number
  default     = 1
}

variable "instance_flavor_size" {
  description = "The disk size in GB of the GeminiDB InfluxDB instance"
  type        = number
  default     = 100
}

variable "instance_flavor_storage" {
  description = "The disk type of the GeminiDB InfluxDB instance"
  type        = string
  default     = "ULTRAHIGH"
}

variable "instance_flavor_spec_code" {
  description = "The resource specification code. If not specified, the first available flavor will be used"
  type        = string
  default     = ""
}

variable "instance_backup_time_window" {
  description = "The backup time window in HH:MM-HH:MM format"
  type        = string
}

variable "instance_backup_keep_days" {
  description = "The number of days to retain backups"
  type        = number
}

variable "cold_storage_size" {
  description = "The size of the cold storage in GB, 0 means no cold storage"
  type        = number
  default     = 0
}

variable "charging_mode" {
  description = "The charging mode. Valid values are prePaid, postPaid"
  type        = string
  default     = "prePaid"
}

variable "period_unit" {
  description = "The charging period unit. Valid values are month, year"
  type        = string
  default     = "month"
}

variable "auto_renew" {
  description = "Whether to enable auto-renew. Valid values are true, false"
  type        = string
  default     = "true"
}

variable "period" {
  description = "The charging period"
  type        = number
  default     = 1
}

variable "tags" {
  description = "The tags of the GeminiDB InfluxDB instance"
  type        = map(string)
  default     = {
    foo = "bar"
    key = "value"
  }
}

# Variable definitions for GeminiDB backup
variable "backup_name" {
  description = "The name of the GeminiDB backup"
  type        = string
  default     = "tf_test_backup"
}

variable "backup_description" {
  description = "The description of the GeminiDB backup"
  type        = string
  default     = "Created by Terraform"
}
