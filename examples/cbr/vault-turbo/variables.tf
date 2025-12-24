# Authentication variables
variable "region_name" {
  description = "The region where the CBR vault and SFS Turbo are located"
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

# Network variables
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
  description = "The CIDR block of the subnet, if not specified, calculating a subnet cidr within the existing CIDR address block"
  type        = string
  default     = ""
}

variable "subnet_gateway_ip" {
  description = "The gateway IP of the subnet, if not specified, calculating a gateway IP within the existing CIDR address block"
  type        = string
  default     = ""
}

variable "availability_zone" {
  description = "The availability zone where the SFS Turbo will be created, if not specified, the first available zone will be used"
  type        = string
  default     = ""
}

variable "secgroup_name" {
  description = "The name of the security group"
  type        = string
}

# SFS Turbo variables
variable "turbo_name" {
  description = "The name of the SFS Turbo file system"
  type        = string
}

variable "turbo_size" {
  description = "The size of the SFS Turbo file system in GB"
  type        = number
}

# CBR vault variables
variable "enable_policy" {
  description = "Whether to enable backup policy"
  type        = bool
  default     = false
}

variable "vault_name" {
  description = "The name of the CBR vault"
  type        = string
}

variable "protection_type" {
  description = "The protection type of the vault (backup or replication)"
  type        = string
  default     = "backup"

  validation {
    condition     = contains(["backup", "replication"], var.protection_type)
    error_message = "The protection type must be either 'backup' or 'replication'."
  }
}

variable "vault_size" {
  description = "The size of the CBR vault in GB"
  type        = number
}

variable "auto_expand" {
  description = "Whether to automatically expand the vault capacity when it's full"
  type        = bool
  default     = false
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project to which the vault belongs"
  type        = string
  default     = "0"
}

variable "backup_name_prefix" {
  description = "The prefix of the backup name"
  type        = string
  default     = ""
}

variable "is_multi_az" {
  description = "Whether the vault is deployed across multiple AZs"
  type        = bool
  default     = false
}

variable "tags" {
  description = "The tags of the vault, in the format of key-value pairs"
  type        = map(string)
  default     = {
    environment = "test"
    terraform   = "true"
    service     = "sfs-turbo"
  }
}
