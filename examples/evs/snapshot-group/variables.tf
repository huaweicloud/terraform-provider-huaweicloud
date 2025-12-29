# Authentication variables
variable "region_name" {
  description = "The region where the EVS snapshot group is located"
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
  default     = 0
}

variable "instance_flavor_memory_size" {
  description = "The memory size in GB for the ECS instance flavor. Using this field to query available flavors if `instance_flavor_id` is not specified."
  type        = number
  default     = 0
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
  description = "The size of the system disk in GB."
  type        = number
  default     = 40
}

# EVS volumes variables
variable "volume_configuration" {
  description = "List of volume configurations to attach to the ECS instance."
  type        = list(object({
    name        = string
    size        = number
    volume_type = string
    device_type = string
  }))

  default = []
}

# EVS snapshot group variables
variable "instant_access" {
  description = "Whether to enable instant access for the snapshot group."
  type        = bool
  default     = false
}

variable "snapshot_group_name" {
  description = "The name of the snapshot group."
  type        = string
  default     = ""
}

variable "snapshot_group_description" {
  description = "The description of the snapshot group."
  type        = string
  default     = "Created by Terraform"
}

variable "enterprise_project_id" {
  description = "The enterprise project ID for the snapshot group."
  type        = string
  default     = "0"
}

variable "incremental" {
  description = "Whether to create an incremental snapshot."
  type        = bool
  default     = false
}

variable "tags" {
  description = "The key/value pairs to associate with the snapshot group."
  type        = map(string)
  default     = {
    environment = "test"
    created_by  = "terraform"
  }
}
