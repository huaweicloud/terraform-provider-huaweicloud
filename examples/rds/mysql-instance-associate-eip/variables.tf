# Variable definitions for authentication
variable "region_name" {
  description = "The region where the MySQL RDS instance is located"
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

variable "availability_zone" {
  description = "The availability zone to which the RDS instance belongs"
  type        = string
  default     = ""
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

variable "instance_flavor_id" {
  description = "The flavor ID of the RDS instance"
  type        = string
  default     = ""
}

variable "instance_db_type" {
  description = "The database engine type"
  type        = string
  default     = "MySQL"
}

variable "instance_db_version" {
  description = "The database engine version"
  type        = string
  default     = "8.0"
}

variable "instance_mode" {
  description = "The instance mode for the RDS instance flavor"
  type        = string
  default     = "single"
}

variable "instance_flavor_group_type" {
  description = "The group type for the RDS instance flavor"
  type        = string
  default     = "general"
}

variable "instance_flavor_vcpus" {
  description = "The number of the RDS instance CPU cores for the RDS instance flavor"
  type        = number
  default     = 2
}

variable "security_group_name" {
  description = "The security group name"
  type        = string
}

variable "instance_db_port" {
  description = "The database port"
  type        = number
  default     = 3306
}

variable "instance_password" {
  description = "The password for the RDS instance"
  type        = string
  default     = ""
  sensitive   = true
}

variable "instance_name" {
  description = "The MySQL RDS instance name"
  type        = string
}

variable "instance_volume_type" {
  description = "The storage volume type"
  type        = string
  default     = "CLOUDSSD"
}

variable "instance_volume_size" {
  description = "The storage volume size in GB"
  type        = number
  default     = 40
}

variable "instance_backup_time_window" {
  description = "The backup time window in HH:MM-HH:MM format"
  type        = string
}

variable "instance_backup_keep_days" {
  description = "The number of days to retain backups"
  type        = number
}

variable "associate_eip_address" {
  description = "The EIP address to associate with the RDS instance"
  type        = string
  default     = ""
}

variable "eip_type" {
  description = "The type of the EIP"
  type        = string
  default     = "5_bgp"
}

variable "bandwidth_name" {
  description = "The name for the bandwidth"
  type        = string
  default     = ""

  validation {
    condition     = var.associate_eip_address != "" || var.bandwidth_name != ""
    error_message = "The bandwidth name must be a non-empty string if the EIP address is not provided."
  }
}

variable "bandwidth_size" {
  description = "The size of the bandwidth"
  type        = number
  default     = 5
}

variable "bandwidth_share_type" {
  description = "The share type of the bandwidth"
  type        = string
  default     = "PER"
}

variable "bandwidth_charge_mode" {
  description = "The charge mode of the bandwidth"
  type        = string
  default     = "traffic"
}
