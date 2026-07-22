# Variable definitions for authentication
variable "region_name" {
  description = "The region where the GeminiDB instance is located"
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
  description = "The number of vCPUs"
  type        = string
  default     = "2"
}

variable "availability_zone" {
  description = "The availability zone to which the GeminiDB instance belongs"
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
  default     = 8635
}

variable "instance_password" {
  description = "The password for the GeminiDB instance"
  type        = string
  default     = ""
  sensitive   = true
}

variable "instance_name" {
  description = "The GeminiDB instance name"
  type        = string
}

variable "instance_mode" {
  description = "The instance mode. Valid values are Cluster, ReplicaSet, Single"
  type        = string
  default     = "Cluster"
}

variable "instance_ssl_option" {
  description = "The SSL option. Valid values are on, off"
  type        = string
  default     = "on"
}

variable "instance_flavor_num" {
  description = "The number of nodes in the instance"
  type        = number
  default     = 3
}

variable "instance_flavor_size" {
  description = "The storage size in GB per node"
  type        = number
  default     = 16
}

variable "instance_flavor_storage" {
  description = "The storage type. Valid values are ULTRAHIGH, ESSD, HIGH, NORMAL"
  type        = string
  default     = "ULTRAHIGH"
}

variable "instance_flavor_spec_code" {
  description = "The resource specification code. If empty, it will be queried from flavors data source"
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
  description = "The key/value pairs to associate with the GeminiDB instance"
  type        = map(string)
  default     = {}
}

variable "backup_name" {
  description = "The name for instance backups"
  type        = string
}

variable "backup_description" {
  description = "The description for instance backups"
  type        = string
  default     = "Terraform created backup"
}
