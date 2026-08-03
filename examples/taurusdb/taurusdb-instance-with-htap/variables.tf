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
  description = "The security group name of the TaurusDB instance"
  type        = string
}

variable "taurusdb_availability_zone_mode" {
  description = "The availability zone mode of the TaurusDB instance. Valid values are single, multi"
  type        = string
  default     = "multi"
}

variable "az_code" {
  description = "The AZ code of the HTAP StarRocks instance. If not specified, the first available AZ from flavors will be used"
  type        = string
  default     = ""
}

variable "fe_flavor_id" {
  description = "The specification ID of the frontend node. If not specified, the first available flavor will be used"
  type        = string
  default     = ""
}

variable "be_flavor_id" {
  description = "The specification ID of the backend node. If not specified, the first available flavor will be used"
  type        = string
  default     = ""
}

variable "engine_version" {
  description = "The major version number of the engine. If not specified, the latest version from the data source will be used"
  type        = string
  default     = ""
}

variable "taurusdb_instance_name" {
  description = "The TaurusDB instance name"
  type        = string
}

variable "taurusdb_flavor_ref" {
  description = "The flavor code of the TaurusDB instance. If not specified, the first available flavor will be used"
  type        = string
  default     = ""
}

variable "taurusdb_root_pwd" {
  description = "The database password of the TaurusDB instance. If not specified, a random password will be generated"
  type        = string
  default     = ""
  sensitive   = true
}

variable "enterprise_project_id" {
  description = "The enterprise project ID"
  type        = string
  default     = "0"
}

variable "taurusdb_read_replicas" {
  description = "The number of read replicas of the TaurusDB instance"
  type        = number
  default     = 2
}

variable "htap_instance_name" {
  description = "The HTAP StarRocks instance name. The name must start with a letter and consist of 4 to 64 characters"
  type        = string
}

variable "fe_count" {
  description = "The number of frontend nodes. For a single-node instance, the value is fixed to 1. For a cluster instance, the value ranges from 3 to 10"
  type        = number
  default     = 1
}

variable "be_count" {
  description = "The number of backend nodes. For a single-node instance, the value is fixed to 1. For a cluster instance, the value ranges from 3 to 10"
  type        = number
  default     = 1
}

variable "htap_db_root_pwd" {
  description = "The database password of the HTAP StarRocks instance. If not specified, a random password will be generated"
  type        = string
  default     = ""
  sensitive   = true
}

variable "time_zone" {
  description = "The time zone of the HTAP StarRocks instance"
  type        = string
  default     = "UTC+08:00"
}

variable "enable_users_sync" {
  description = "Whether to enable users synchronization. Valid values are true, false"
  type        = string
  default     = "true"
}

variable "open_slow_log_switch" {
  description = "Whether to enable the slow query log original text switch. Valid values are true, false"
  type        = string
  default     = "true"
}

variable "ha_mode" {
  description = "The deployment mode of the HTAP StarRocks instance. Valid values are Single, Cluster"
  type        = string
  default     = "Single"
}

variable "volume_io_type" {
  description = "The storage type of the frontend and backend nodes. Valid values are SSD, ESSD"
  type        = string
  default     = "SSD"
}

variable "fe_volume_capacity" {
  description = "The disk capacity in GB of the frontend node. The value ranges from 50 to 1000 and the increment is 10 GB"
  type        = number
  default     = 50
}

variable "be_volume_capacity" {
  description = "The disk capacity in GB of the backend node. The value ranges from 50 to 1000 and the increment is 10 GB"
  type        = number
  default     = 50
}

variable "be_parameter_values" {
  description = "A map contains mappings of parameter name and value to modify for the backend nodes"
  type        = map(string)
  default     = {}
}

variable "fe_parameter_values" {
  description = "A map contains mappings of parameter name and value to modify for the frontend nodes"
  type        = map(string)
  default     = {}
}
