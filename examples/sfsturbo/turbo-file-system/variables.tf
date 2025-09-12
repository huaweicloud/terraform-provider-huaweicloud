# Variable definitions for authentication
variable "region_name" {
  description = "The region where the SFS Turbo file system is located"
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
  description = "The name of the VPC"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "192.168.0.0/16"
}

variable "subnet_name" {
  description = "The name of the subnet"
  type        = string
}

variable "subnet_cidr" {
  description = "The CIDR block of the subnet"
  type        = string
  default     = ""
  nullable    = false
}

variable "subnet_gateway_ip" {
  description = "The gateway IP of the subnet"
  type        = string
  default     = ""
  nullable    = false
}

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "turbo_name" {
  description = "The name of the SFS Turbo file system"
  type        = string
  default     = ""
}

variable "turbo_size" {
  description = "The capacity of the SFS Turbo file system"
  type        = number
  default     = 1228
}

variable "turbo_share_proto" {
  description = "The protocol of the SFS Turbo file system"
  type        = string
  default     = "NFS"
}

variable "turbo_share_type" {
  description = "The type of the SFS Turbo file system"
  type        = string
  default     = "STANDARD"
}

variable "turbo_hpc_bandwidth" {
  description = "The bandwidth specification of the SFS Turbo file system"
  type        = string
  default     = ""
}

variable "turbo_tags" {
  description = "The tags of the SFS Turbo file system"
  type        = map(string)
  default     = {}
}

variable "turbo_backup_id" {
  description = "The ID of the backup"
  type        = string
  default     = ""
}

variable "enterprise_project_id" {
  description = "The enterprise project ID of the SFS Turbo file system"
  type        = string
  default     = null
}

variable "charging_mode" {
  description = "The charging mode of the SFS Turbo file system"
  type        = string
  default     = "postPaid"
}

variable "period_unit" {
  description = "The period unit of the SFS Turbo file system"
  type        = string
  default     = null
}

variable "period" {
  description = "The period of the SFS Turbo file system"
  type        = number
  default     = null
}

variable "auto_renew" {
  description = "The auto renew of the SFS Turbo file system"
  type        = string
  default     = "false"
}
