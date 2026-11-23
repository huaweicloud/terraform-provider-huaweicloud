# Variable definitions for authentication
variable "region_name" {
  description = "The region where the volume is located"
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
  default     = 0
}

variable "instance_flavor_memory_size" {
  description = "The memory size in GB for the ECS instance flavor. Using this field to query available flavors if `instance_flavor_id` is not specified."
  type        = number
  default     = 0
}

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
  description = "The size of the system disk in GB."
  type        = number
  default     = 40
}

variable "volume_availability_zone" {
  description = "The availability zone for the volume"
  type        = string
  default     = ""
  nullable    = false
}

variable "volume_name" {
  description = "The name of the volume"
  type        = string
}

variable "volume_type" {
  description = "The type of the volume"
  type        = string
  default     = "SSD"
}

variable "volume_size" {
  description = "The size of the volume"
  type        = number
  default     = 40
}

variable "volume_description" {
  description = "The description of the volume"
  type        = string
  default     = ""
}

variable "volume_multiattach" {
  description = "The volume is shared volume or not"
  type        = bool
  default     = false
}

variable "volume_iops" {
  description = "The IOPS for the volume"
  type        = number
  default     = null
}

variable "volume_throughput" {
  description = "The throughput for the volume"
  type        = number
  default     = null
}

variable "volume_device_type" {
  description = "The device type of disk to create"
  type        = string
  default     = "VBD"
}

variable "enterprise_project_id" {
  description = "The enterprise project ID"
  type        = string
  default     = null
}

variable "volume_tags" {
  description = "The tags of the volume"
  type        = map(string)
  default     = {}
}

variable "charging_mode" {
  description = "The charging mode of the volume"
  type        = string
  default     = "postPaid"
}

variable "period_unit" {
  description = "The period unit of the volume"
  type        = string
  default     = null
}

variable "period" {
  description = "The period of the volume"
  type        = number
  default     = null
}

variable "auto_renew" {
  description = "The auto renew of the volume"
  type        = string
  default     = "false"
}

variable "volume_metadata" {
  description = "The volume metadata"
  type        = map(string)
  default     = {}
}

variable "snapshot_name" {
  description = "The name of the snapshot"
  type        = string
}

variable "snapshot_description" {
  description = "The description of the snapshot"
  type        = string
  default     = ""
}

variable "snapshot_tags" {
  description = "The tags of the snapshot"
  type        = map(string)
  default     = {}
}

variable "snapshot_instant_access" {
  description = "Whether to enable instant access for the snapshot"
  type        = bool
  default     = false
}

variable "snapshot_incremental" {
  description = "Whether to create an incremental snapshot"
  type        = bool
  default     = false
}

variable "snapshot_metadata" {
  description = "The snapshot metadata"
  type        = map(string)
  default     = {}
}
