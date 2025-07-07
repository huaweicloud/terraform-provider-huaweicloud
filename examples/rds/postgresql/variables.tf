variable "region_name" {
  description = "The region name"
  type        = string
}

variable "access_key" {
  description = "The access key of the IAM user"
  sensitive   = true
  type        = string
}

variable "secret_key" {
  description = "The secret key of the IAM user"
  sensitive   = true
  type        = string
}

variable "availability_zone" {
  description = "The availability zone"
  type        = string
  default     = ""
}

variable "subnet_cidr" {
  description = "The CIDR block of the subnet"
  type        = string
  default     = ""
}

variable "gateway" {
  description = "The gateway IP of the subnet"
  type        = string
  default     = ""
}

variable "flavor_id" {
  description = "The flavor ID for the instance"
  type        = string
  default     = ""
}

variable "db_type" {
  description = "The database engine type"
  type        = string
  default     = "PostgreSQL"
}

variable "db_version" {
  description = "The database engine version"
  type        = string
  default     = "16"
}

variable "instance_mode" {
  description = "The instance mode for the RDS instance flavor"
  type        = string
  default     = "single"
}

variable "group_type" {
  description = "The group type"
  type        = string
  default     = "general"
}

variable "vcpus" {
  description = "The CPU of flavor for the instance"
  type        = number
  default     = 2
}

variable "vpc_id" {
  description = "The ID of the existing VPC"
  type        = string
  default     = ""
}

variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "192.168.0.0/16"
}

variable "subnet_id" {
  description = "The ID of the existing subnet"
  type        = string
  default     = ""
}

variable "subnet_name" {
  description = "The name of the subnet"
  type        = string
}

variable "secgroup_id" {
  description = "The ID of the existing security group"
  type        = string
  default     = ""
}

variable "secgroup_name" {
  description = "The name of security group"
  type        = string
}

variable "db_port" {
  description = "The database port"
  type        = number
  default     = 5432
}

variable "instance_name" {
  description = "The name of the RDS instance"
  type        = string
}

variable "charging_mode" {
  description = "The billing method"
  type        = string
  default     = "postPaid"
}

variable "volume_type" {
  description = "The storage volume type"
  type        = string
  default     = "CLOUDSSD"
}

variable "volume_size" {
  description = "The storage volume size in GB"
  type        = number
  default     = 40
}

variable "backup_time_window" {
  description = "The backup time window in HH:MM-HH:MM format"
  type        = string
}

variable "backup_keep_days" {
  description = "The number of days to retain backups"
  type        = number
}

variable "account_name" {
  description = "Username with elevated privileges"
  type        = string
}

variable "db_name" {
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
