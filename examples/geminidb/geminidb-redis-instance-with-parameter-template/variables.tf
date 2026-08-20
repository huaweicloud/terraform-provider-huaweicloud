# Variable definitions for authentication
variable "region_name" {
  description = "The region where the GeminiDB Redis instance is located"
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

# Variable definitions for VPC and subnet
variable "vpc_name" {
  description = "The VPC name"
  type        = string
  nullable    = false
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  nullable    = false
  default     = "192.168.0.0/16"
}

variable "subnet_name" {
  description = "The subnet name"
  type        = string
  nullable    = false
}

variable "subnet_cidr" {
  description = "The CIDR block of the subnet"
  type        = string
  nullable    = false
  default     = ""
}

variable "gateway_ip" {
  description = "The gateway IP address of the subnet"
  type        = string
  nullable    = false
  default     = ""
}

# Variable definitions for security group
variable "security_group_name" {
  description = "The security group name"
  type        = string
  nullable    = false
}

variable "instance_db_port" {
  description = "The database port"
  type        = number
  nullable    = false
  default     = 8635
}

# Variable definitions for flavor query
variable "availability_zone" {
  description = "The availability zone to which the GeminiDB Redis instance belongs"
  type        = string
  nullable    = false
  default     = ""
}

# Variable definitions for parameter template
variable "parameter_template_name" {
  description = "The name of the GeminiDB Redis parameter template"
  type        = string
  nullable    = false
}

variable "parameter_template_description" {
  description = "The description of the GeminiDB Redis parameter template"
  type        = string
  nullable    = false
  default     = ""
}

# Variable definitions for datastore
variable "datastore_type" {
  description = "The database type"
  type        = string
  nullable    = false
  default     = "redis"
}

variable "datastore_version" {
  description = "The database version"
  type        = string
  nullable    = false
  default     = "5.0"
}

variable "parameter_template_values" {
  description = "The parameter key-value pairs for the GeminiDB Redis parameter template"
  type        = map(string)
  nullable    = false
  default     = {}
}

# Variable definitions for GeminiDB Redis instance
variable "instance_name" {
  description = "The GeminiDB Redis instance name"
  type        = string
  nullable    = false
}

variable "instance_password" {
  description = "The password for the GeminiDB Redis instance"
  type        = string
  sensitive   = true
  nullable    = false
  default     = ""
}

variable "instance_mode" {
  description = "The instance mode"
  type        = string
  nullable    = false
  default     = "Cluster"
}

variable "instance_ssl_option" {
  description = "The SSL option"
  type        = string
  nullable    = false
  default     = "on"
}

variable "datastore_storage_engine" {
  description = "The storage engine"
  type        = string
  nullable    = false
  default     = "rocksDB"
}

variable "instance_flavor_num" {
  description = "The number of nodes in the instance"
  type        = number
  nullable    = false
  default     = 3
}

variable "instance_flavor_size" {
  description = "The storage size in GB per node"
  type        = number
  nullable    = false
  default     = 16
}

variable "instance_flavor_storage" {
  description = "The storage type"
  type        = string
  nullable    = false
  default     = "ULTRAHIGH"
}

variable "instance_flavor_spec_code" {
  description = "The resource specification code"
  type        = string
  nullable    = false
  default     = ""
}

variable "instance_backup_time_window" {
  description = "The backup time window in HH:MM-HH:MM format"
  type        = string
  nullable    = false
  default     = "00:00-01:00"
}

variable "instance_backup_keep_days" {
  description = "The number of days to retain backups"
  type        = number
  nullable    = false
  default     = 7
}

variable "tags" {
  description = "The key/value pairs to associate with the GeminiDB Redis instance"
  type        = map(string)
  nullable    = false
  default     = {}
}
