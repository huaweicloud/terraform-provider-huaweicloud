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
  default     = 2
}

variable "availability_zone" {
  description = "The availability zone of the GeminiDB Redis instance. If not specified, the first AZ from data source will be used"
  type        = string
  default     = ""
}

variable "security_group_name" {
  description = "The security group name"
  type        = string
}

variable "instance_password" {
  description = "The password for the GeminiDB Redis instance"
  type        = string
  default     = ""
  sensitive   = true
}

variable "instance_name" {
  description = "The GeminiDB Redis instance name"
  type        = string
}

variable "instance_mode" {
  description = "The instance mode. Valid value is Cluster"
  type        = string
  default     = "Cluster"
}

variable "enterprise_project_id" {
  description = "The enterprise project ID"
  type        = string
  default     = "0"
}

variable "instance_db_port" {
  description = "The database port of the GeminiDB Redis instance"
  type        = number
  default     = 8888
}

variable "instance_ssl_option" {
  description = "Whether to enable SSL. Valid values are on, off"
  type        = string
  default     = "on"
}

variable "instance_flavor_num" {
  description = "The node quantity of the GeminiDB Redis instance"
  type        = number
  default     = 3
}

variable "instance_flavor_size" {
  description = "The disk size in GB of the GeminiDB Redis instance"
  type        = number
  default     = 16
}

variable "instance_flavor_storage" {
  description = "The disk type of the GeminiDB Redis instance"
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

variable "tags" {
  description = "The tags of the GeminiDB Redis instance"
  type        = map(string)
  default     = {
    foo = "bar"
    key = "value"
  }
}

# Account variables
variable "account_name" {
  description = "The username of the GeminiDB account"
  type        = string
}

variable "account_password" {
  description = "The password of the GeminiDB account. If not specified, a random password will be generated"
  type        = string
  default     = ""
  sensitive   = true
}

variable "account_privilege" {
  description = "The privilege of the GeminiDB account. Valid values are ReadWrite, ReadOnly"
  type        = string
  default     = "ReadWrite"
}

variable "account_databases" {
  description = "The list of databases for the GeminiDB account"
  type        = list(string)
  default     = ["1", "2"]
}
