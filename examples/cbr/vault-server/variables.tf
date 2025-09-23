# Authentication variables
variable "region_name" {
  description = "The region where the CBR vault is located"
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

# ECS Flavor variables
variable "instance_flavor_id" {
  description = "The flavor ID of the ECS instance, if not specified, the first available flavor matching the criteria will be used."
  type        = string
  default     = ""
}

variable "availability_zone" {
  description = "The availability zone where the ECS instance will be created, if not specified, the first available zone will be used"
  type        = string
  default     = ""
}

variable "instance_flavor_performance_type" {
  description = "The performance type of the ECS instance flavor. Using this field to query available flavors if `instance_flavor_id` is not specified."
  type        = string
  default     = "normal"
}

variable "instance_flavor_cpu_core_count" {
  description = "The number of CPU cores for the ECS instance flavor. Using this field to query available flavors if `instance_flavor_id` is not specified."
  type        = number
  default     = 2
}

variable "instance_flavor_memory_size" {
  description = "The memory size in GB for the ECS instance flavor. Using this field to query available flavors if `instance_flavor_id` is not specified."
  type        = number
  default     = 4
}

# ECS Image variables
variable "instance_image_id" {
  description = "The ID of the image used to create the ECS instance, if not specified, the first available image matching the criteria will be used."
  type        = string
  default     = ""
}

variable "instance_image_os_type" {
  description = "The OS type of the ECS instance image. Using this field to query available images if `instance_image_id` is not specified."
  type        = string
  default     = "Ubuntu"
}

variable "instance_image_visibility" {
  description = "The visibility of the ECS instance image. Using this field to query available images if `instance_image_id` is not specified."
  type        = string
  default     = "public"
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

variable "secgroup_name" {
  description = "The name of the security group"
  type        = string
}

# ECS variables
variable "ecs_instance_name" {
  description = "The name of the ECS instance"
  type        = string
}

variable "key_pair_name" {
  description = "The name of the key pair for ECS login"
  type        = string
  default     = ""
}

variable "system_disk_type" {
  description = "The type of the system disk"
  type        = string
  default     = "SAS"
}

variable "system_disk_size" {
  description = "The size of the system disk in GB"
  type        = number
  default     = 40
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

variable "consistent_level" {
  description = "The consistent level of the vault (crash_consistent or app_consistent)"
  type        = string
  default     = "crash_consistent"

  validation {
    condition     = contains(["crash_consistent", "app_consistent"], var.consistent_level)
    error_message = "The consistent level must be either 'crash_consistent' or 'app_consistent'."
  }
}

variable "vault_size" {
  description = "The size of the CBR vault in GB"
  type        = number
}

variable "auto_bind" {
  description = "Whether to automatically bind the vault to a policy"
  type        = bool
  default     = false
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

variable "exclude_volumes" {
  description = "Whether to exclude volumes from backup"
  type        = bool
  default     = false
}

variable "tags" {
  description = "The tags of the vault"
  type        = map(string)
  default     = {
    environment = "test"
    terraform   = "true"
  }
}
