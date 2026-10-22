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
variable "rds1_name" {
  description = "The name of the first RDS instance"
  type        = string
}

variable "rds2_name" {
  description = "The name of the second RDS instance"
  type        = string
}

variable "rds_flavor" {
  description = "The flavor of the RDS instances"
  type        = string
  default     = "rds.mysql.x1.large.2.ha"
}

variable "rds1_fixed_ip" {
  description = "The fixed IP address of the first RDS instance"
  type        = string
  default     = "192.168.0.50"
}

variable "rds2_fixed_ip" {
  description = "The fixed IP address of the second RDS instance"
  type        = string
  default     = "192.168.0.51"
}

variable "db_password" {
  description = "The password for the RDS root user and DRS database connections"
  type        = string
  sensitive   = true
}

variable "db_name" {
  description = "The name of the database to be synchronized in both directions"
  type        = string
  default     = "test"
}

# Variable definitions for DRS jobs
variable "forward_job_name" {
  description = "The name of the forward DRS job (from the first RDS instance to the second one)"
  type        = string
}

variable "reverse_job_name" {
  description = "The name of the reverse DRS job (from the second RDS instance to the first one)"
  type        = string
}

variable "forward_job_description" {
  description = "The description of the forward DRS job"
  type        = string
  default     = ""
}

variable "reverse_job_description" {
  description = "The description of the reverse DRS job"
  type        = string
  default     = ""
}

variable "db_user" {
  description = "The database user name used by the RDS instances and DRS jobs"
  type        = string
  default     = "root"
}

variable "limit_speed" {
  description = "The speed limit of the DRS jobs. If not specified, the jobs are not limited"
  type        = object({
    speed      = string
    start_time = string
    end_time   = string
  })
  default     = null
}
