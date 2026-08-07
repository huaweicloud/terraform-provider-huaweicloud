# Variable definitions for authentication
variable "region_name" {
  description = "The primary region where the GaussDB instance is located"
  type        = string
  default     = "cn-north-4"
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

# Primary Region - Network Resources
variable "primary_vpc_name" {
  description = "The VPC name in the primary region"
  type        = string
}

variable "primary_vpc_cidr" {
  description = "The CIDR block of the VPC in the primary region"
  type        = string
  default     = "172.16.0.0/16"
}

variable "primary_subnet_names" {
  description = "The subnet names for each AZ in the primary region"
  type        = list(string)
}

variable "primary_availability_zones" {
  description = "The availability zones in the primary region"
  type        = list(string)
}

variable "primary_security_group_name" {
  description = "The security group name in the primary region"
  type        = string
}

# Security Group Rules
variable "secgroup_rules" {
  description = "The security group ingress rules"
  type        = list(object({
    ports  = string
    source = string
  }))
  default     = [
    { ports = "5432-5532", source = "local" },
    { ports = "5432-5532", source = "remote" },
    { ports = "20050",     source = "local" },
    { ports = "20050",     source = "remote" },
    { ports = "5000-5001", source = "local" },
    { ports = "5000-5001", source = "remote" },
    { ports = "2379-2380", source = "local" },
    { ports = "2379-2380", source = "remote" },
    { ports = "6000",      source = "local" },
    { ports = "6000",      source = "remote" },
    { ports = "6500",      source = "local" },
    { ports = "6500",      source = "remote" },
    { ports = "12016",     source = "local" },
    { ports = "12016",     source = "remote" },
  ]
}

# DR Region - Network Resources
variable "dr_vpc_name" {
  description = "The VPC name in the DR region"
  type        = string
}

variable "dr_vpc_cidr" {
  description = "The CIDR block of the VPC in the DR region"
  type        = string
  default     = "172.17.0.0/16"
}

variable "dr_subnet_name" {
  description = "The subnet name in the DR region"
  type        = string
}

variable "dr_availability_zone" {
  description = "The availability zone in the DR region"
  type        = string
}

variable "dr_security_group_name" {
  description = "The security group name in the DR region"
  type        = string
}

variable "dr_region_name" {
  description = "The standby region where the DR GaussDB instance is located"
  type        = string
  default     = "cn-east-3"
}

# Cloud Connection
variable "cc_bandwidth" {
  description = "The inter-region bandwidth (Mbit/s) for the Cloud Connection"
  type        = number
  default     = 10
}

# Primary Region - GaussDB Instance
variable "instance_passwords" {
  description = "The passwords for GaussDB instances (empty string to auto-generate), index 0=primary, 1=dr"
  type        = list(string)
  default     = ["", ""]
  sensitive   = true
}

variable "primary_instance_name" {
  description = "The name of the primary GaussDB instance"
  type        = string
}

variable "instance_flavor" {
  description = "The spec_code of the GaussDB instance flavor"
  type        = string
}

variable "primary_instance_availability_zones" {
  description = "The comma-separated AZ string for the primary GaussDB instance"
  type        = string
}

variable "instance_db_port" {
  description = "The database port"
  type        = number
  default     = 5432
}

variable "enterprise_project_id" {
  description = "The enterprise project ID to which the GaussDB instances belong"
  type        = string
}

variable "primary_instance_volume_type" {
  description = "The storage volume type of the primary GaussDB instance"
  type        = string
  default     = "ULTRAHIGH"
}

variable "primary_instance_volume_size" {
  description = "The storage volume size (GB) of the primary GaussDB instance"
  type        = number
  default     = 40
}

# DR Region - GaussDB Instance
variable "dr_instance_name" {
  description = "The name of the DR GaussDB instance"
  type        = string
}

variable "dr_instance_availability_zones" {
  description = "The comma-separated AZ string for the DR GaussDB instance"
  type        = string
}

variable "dr_instance_volume_type" {
  description = "The storage volume type of the DR GaussDB instance"
  type        = string
  default     = "ULTRAHIGH"
}

variable "dr_instance_volume_size" {
  description = "The storage volume size (GB) of the DR GaussDB instance"
  type        = number
  default     = 40
}

# Disaster Recovery Relationship
variable "dr_disaster_type" {
  description = "The disaster recovery type"
  type        = string
  default     = "stream"
}

variable "dr_user_name" {
  description = "The database username for disaster recovery"
  type        = string
}

variable "dr_user_password" {
  description = "The database password for disaster recovery"
  type        = string
  sensitive   = true
}
