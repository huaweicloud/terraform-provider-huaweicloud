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

variable "target_bucket_name" {
  description = "The name of the OBS bucket"
  type        = string
  default     = ""
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
  default     = "HPC"
}

variable "turbo_hpc_bandwidth" {
  description = "The bandwidth specification of the SFS Turbo file system"
  type        = string
  default     = "40M"
}

variable "target_file_path" {
  description = "The linkage directory name of the OBS target"
  type        = string
  default     = "testDir"
}

variable "target_obs_endpoint" {
  description = "The domain name of the region where the OBS bucket located"
  type        = string
  default     = "obs.cn-north-4.myhuaweicloud.com"
}

variable "target_events" {
  description = "The type of the data automatically exported to the OBS bucket"
  type        = list(string)
  default     = []
}

variable "target_prefix" {
  description = "The prefix to be matched in the storage backend"
  type        = string
  default     = ""
}

variable "target_suffix" {
  description = "The suffix to be matched in the storage backend"
  type        = string
  default     = ""
}

variable "target_file_mode" {
  description = "The permissions on the imported file"
  type        = string
  default     = ""
}

variable "target_dir_mode" {
  description = "The permissions on the imported directory"
  type        = string
  default     = ""
}

variable "target_uid" {
  description = "The ID of the user who owns the imported object"
  type        = number
  default     = 0
}

variable "target_gid" {
  description = "The ID of the user group to which the imported object belongs"
  type        = number
  default     = 0
}
