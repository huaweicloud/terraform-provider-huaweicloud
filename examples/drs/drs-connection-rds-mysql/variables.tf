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

# Variable definitions for VPC and network
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

# Variable definitions for RDS instance
variable "rds_db_type" {
  description = "The database type for querying RDS flavors"
  type        = string
  default     = "MySQL"
}

variable "rds_db_version" {
  description = "The database version for querying RDS flavors"
  type        = string
  default     = "5.7"
}

variable "rds_instance_mode" {
  description = "The instance mode for querying RDS flavors"
  type        = string
  default     = "single"
}

variable "rds_name" {
  description = "The name of the RDS instance"
  type        = string
}

variable "rds_flavor" {
  description = "The flavor of the RDS instance. If not specified, it will be queried from data source"
  type        = string
  default     = ""
}

variable "rds_fixed_ip" {
  description = "The fixed IP address of the RDS instance"
  type        = string
  default     = "192.168.0.100"
}

variable "db_password" {
  description = "The password for the RDS root user and DRS connection"
  type        = string
  sensitive   = true
}

# Variable definitions for DRS connection
variable "connection_name" {
  description = "The DRS connection name"
  type        = string
}

variable "description" {
  description = "The description of the DRS connection"
  type        = string
  default     = ""
}

variable "db_port" {
  description = "The database port"
  type        = string
  default     = "3306"
}

variable "db_user" {
  description = "The database username"
  type        = string
  default     = "root"
}

variable "driver_name" {
  description = "The driver name of the connection configuration"
  type        = string
  default     = "mysql"
}
