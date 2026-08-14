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
  nullable    = false
}

variable "gateway_ip" {
  description = "The gateway IP address of the subnet"
  type        = string
  default     = ""
  nullable    = false
}

variable "security_group_name" {
  description = "The security group name"
  type        = string
}

# Variable definitions for RDS instances
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
  default     = "ha"
}

variable "source_rds_name" {
  description = "The name of the source RDS instance"
  type        = string
}

variable "dest_rds_name" {
  description = "The name of the destination RDS instance"
  type        = string
}

variable "rds_flavor" {
  description = "The flavor of the RDS instances"
  type        = string
  default     = "rds.mysql.x1.large.2.ha"
}

variable "source_rds_fixed_ip" {
  description = "The fixed IP address of the source RDS instance"
  type        = string
}

variable "dest_rds_fixed_ip" {
  description = "The fixed IP address of the destination RDS instance"
  type        = string
}

variable "db_password" {
  description = "The password for the RDS root user and DRS database connections"
  type        = string
  sensitive   = true
}

# Variable definitions for DRS job
variable "job_name" {
  description = "The DRS job name"
  type        = string
}

variable "description" {
  description = "The description of the DRS job"
  type        = string
  default     = ""
}

# Variable definitions for LTS
variable "lts_group_name" {
  description = "The name of the LTS log group"
  type        = string
}

variable "lts_ttl_in_days" {
  description = "The log retention period in days"
  type        = number
  default     = 30
}

variable "lts_stream_name" {
  description = "The name of the LTS log stream"
  type        = string
}
