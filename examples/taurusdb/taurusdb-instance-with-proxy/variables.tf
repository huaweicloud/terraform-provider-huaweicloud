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

variable "master_availability_zone" {
  description = "The master availability zone of the TaurusDB instance. If not specified, the first available AZ from flavors will be used"
  type        = string
  default     = ""
}

variable "security_group_name" {
  description = "The security group name"
  type        = string
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

variable "read_replicas" {
  description = "The number of read replicas"
  type        = number
  default     = 4
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

variable "instance_db_port" {
  description = "The database port"
  type        = number
  default     = 3306
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

# Proxy variables
variable "proxy_node_num" {
  description = "The number of proxy nodes"
  type        = number
  default     = 2
}

variable "proxy_name" {
  description = "The name of the database proxy"
  type        = string
  default     = ""
}

variable "proxy_mode" {
  description = "The type of the proxy. Valid values are readwrite, readonly. Defaults to readwrite"
  type        = string
  default     = "readwrite"
}

variable "route_mode" {
  description = "The routing policy of the proxy. Valid values are 0, 1, 2"
  type        = number
  default     = 1
}

variable "proxy_new_node_auto_add_status" {
  description = "Whether to automatically add new nodes to the proxy. Valid values are ON, OFF"
  type        = string
  default     = "OFF"
}

variable "proxy_new_node_weight" {
  description = "The weight of new nodes automatically added to the proxy"
  type        = number
  default     = 20
}

variable "proxy_port" {
  description = "The proxy port"
  type        = number
  default     = 3306
}

variable "proxy_transaction_split" {
  description = "Whether to enable transaction split. Valid values are ON, OFF. Defaults to OFF"
  type        = string
  default     = "OFF"
}

variable "proxy_consistence_mode" {
  description = "The consistency mode. Valid values are session, global, eventual. Defaults to eventual"
  type        = string
  default     = "eventual"
}

variable "proxy_connection_pool_type" {
  description = "The connection pool type. Valid values are SESSION, CLOSED. Defaults to CLOSED"
  type        = string
  default     = "CLOSED"
}

variable "proxy_open_access_control" {
  description = "Whether to enable access control"
  type        = bool
  default     = true
}

variable "access_control_type" {
  description = "The access control mode. Valid values are white, black"
  type        = string
  default     = "white"
}

variable "proxy_dns_name_prefix" {
  description = "The DNS name prefix for the proxy"
  type        = string
  default     = ""
}

variable "proxy_access_control_ip_list" {
  description = "The access control IP list for the proxy"
  type        = list(object({
    ip          = string
    description = string
  }))
  default     = []
}

variable "proxy_master_node_weight" {
  description = "The weight of the master node in the proxy"
  type        = number
  default     = 20
}

variable "proxy_readonly_node_weight" {
  description = "The weight of read-only nodes in the proxy"
  type        = number
  default     = 30
}

variable "proxy_parameters" {
  description = "The parameters for the proxy"
  type        = list(object({
    name      = string
    value     = string
    elem_type = string
  }))
  default     = []
}
