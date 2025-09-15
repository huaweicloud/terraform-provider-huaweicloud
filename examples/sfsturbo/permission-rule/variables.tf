# Variable definitions for authentication
variable "region_name" {
  description = "The region where the SFS Turbo service is located"
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

variable "turbo_availability_zone" {
  description = "The availability zone to which the SFS Turbo belongs"
  type        = string
  default     = ""
  nullable    = false
}

variable "turbo_name" {
  description = "The name of the SFS Turbo file system"
  type        = string
}

variable "turbo_size" {
  description = "The capacity of the SFS Turbo file system"
  type        = number
  default     = 500
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

variable "rule_ip_cidr" {
  description = "The IP address or network segment of the authorized object"
  type        = string
  default     = "192.168.0.0/16"
}

variable "rule_rw_type" {
  description = "The read and write permissions for the authorized object"
  type        = string
  default     = "rw"
}

variable "rule_user_type" {
  description = "The access permissions of the system users of the authorized object to the file system"
  type        = string
  default     = "no_root_squash"
}
