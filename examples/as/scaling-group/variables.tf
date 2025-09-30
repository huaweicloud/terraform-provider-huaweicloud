# Variables definitions for authorization
variable "region_name" {
  description = "The region where the scaling configuration is located"
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

# Variables definitions for resource/data source
variable "availability_zone" {
  description = "The availability zone to which the scaling configuration belongs"
  type        = string
  default     = ""
}

variable "configuration_flavor_id" {
  description = "The flavor ID of the scaling configuration"
  type        = string
  default     = ""
}

variable "configuration_flavor_performance_type" {
  description = "The performance type of the scaling configuration"
  type        = string
  default     = "normal"
}

variable "configuration_flavor_cpu_core_count" {
  description = "The CPU core count of the scaling configuration"
  type        = number
  default     = 2
}

variable "configuration_flavor_memory_size" {
  description = "The memory size of the scaling configuration"
  type        = number
  default     = 4
}

variable "configuration_image_id" {
  description = "The image ID of the scaling configuration"
  type        = string
  default     = ""
}

variable "configuration_image_visibility" {
  description = "The visibility of the image"
  type        = string
  default     = "public"
}

variable "configuration_image_os" {
  description = "The OS of the image"
  type        = string
  default     = "Ubuntu"
}

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "keypair_name" {
  description = "The name of the key pair"
  type        = string
}

variable "keypair_public_key" {
  description = "The public key for SSH access"
  type        = string
  default     = ""
}

variable "configuration_name" {
  description = "The name of the scaling configuration"
  type        = string
}

variable "configuration_metadata" {
  description = "The metadata for the scaling configuration instances"
  type        = map(string)
}

variable "configuration_user_data" {
  description = "The user data script for scaling configuration instances initialization"
  type        = string
}

variable "configuration_disks" {
  description = "The disk configurations for the scaling configuration instances"
  type        = list(object({
    size        = number
    volume_type = string
    disk_type   = string
  }))

  nullable = false
}

variable "configuration_public_eip_settings" {
  description = "The public IP settings for the scaling configuration instances"
  type        = list(object({
    ip_type = string
    bandwidth = object({
      size          = number
      share_type    = string
      charging_mode = string
    })
  }))

  nullable = false
  default  = []
}

variable "scaling_group_vpc_id" {
  description = "The ID of the VPC"
  type        = string
  default     = ""
}

variable "scaling_group_subnet_id" {
  description = "The ID of the subnet"
  type        = string
  default     = ""
}

variable "scaling_group_vpc_name" {
  description = "The name of the VPC"
  type        = string
  default     = ""
  nullable    = false

  validation {
    condition     = var.scaling_group_vpc_id == "" && var.scaling_group_vpc_name != ""
    error_message = "The 'scaling_group_vpc_id' and 'scaling_group_vpc_name' is not allowed for only one of them to be empty"
  }
}

variable "scaling_group_vpc_cidr" {
  description = "The CIDR of the VPC"
  type        = string
  default     = "192.168.0.0/16"
}

variable "scaling_group_subnet_name" {
  description = "The name of the subnet"
  type        = string
  default     = ""
  nullable    = false

  validation {
    condition     = var.scaling_group_subnet_id == "" && var.scaling_group_subnet_name != ""
    error_message = "The 'scaling_group_subnet_id' and 'scaling_group_subnet_name' is not allowed for only one of them to be empty"
  }
}

variable "scaling_group_subnet_cidr" {
  description = "The CIDR of the subnet"
  type        = string
  default     = ""
}

variable "scaling_group_subnet_gateway_ip" {
  description = "The gateway IP of the subnet"
  type        = string
  default     = ""
}

variable "scaling_group_name" {
  description = "The name of the scaling group"
  type        = string
}

variable "scaling_group_desire_instance_number" {
  description = "The desired number of instances"
  type        = number
  default     = 2
}

variable "scaling_group_min_instance_number" {
  description = "The minimum number of instances"
  type        = number
  default     = 0
}

variable "scaling_group_max_instance_number" {
  description = "The maximum number of instances"
  type        = number
  default     = 10
}

variable "is_delete_scaling_group_publicip" {
  description = "Whether to delete the public IP address of the scaling instances when the scaling group is deleted"
  type        = bool
  default     = true
}

variable "is_delete_scaling_group_instances" {
  description = "Whether to delete the scaling instances when the scaling group is deleted"
  type        = string
  default     = "yes"
}
