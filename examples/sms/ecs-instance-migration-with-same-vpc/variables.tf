# Variables definitions for authorization
variable "region_name" {
  description = "The region where the resources are located"
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
variable "availability_zone" {
  description = "The name of the availability zone to which the resources belong"
  type        = string
  default     = ""
  nullable    = false
}

variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "172.16.0.0/16"
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
  description = "The gateway IP address of the subnet"
  type        = string
  default     = ""
  nullable    = false
}

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "instance_flavor_id" {
  description = "The flavor ID of the ECS instance"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_flavor_performance_type" {
  description = "The performance type of the instance flavor"
  type        = string
  default     = "normal"
}

variable "instance_flavor_cpu_core_count" {
  description = "The CPU core count of the instance flavor"
  type        = number
  default     = 2
}

variable "instance_flavor_memory_size" {
  description = "The memory size of the instance flavor"
  type        = number
  default     = 4
}

variable "instance_image_id" {
  description = "The image ID of the instance"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_image_visibility" {
  description = "The visibility of the instance image"
  type        = string
  default     = "public"
}

variable "instance_image_os" {
  description = "The OS of the instance image"
  type        = string
  default     = "Ubuntu"
}

variable "instance_name" {
  description = "The name of the ECS instance"
  type        = string
}

variable "instance_admin_password" {
  description = "The login password of the ECS instance"
  type        = string
  sensitive   = true
}

variable "instance_system_disk_size" {
  description = "The size of the ECS instance system disk in GB"
  type        = number
  default     = 40
}

variable "instance_system_disk_type" {
  description = "The type of the ECS instance system disk"
  type        = string
  default     = "GPSSD"
}

variable "destination_instance_name" {
  description = "The name of the destination ECS instance"
  type        = string
}

variable "destination_instance_system_disk_size" {
  description = "The size of the destination ECS instance system disk in GB"
  type        = number
  default     = 80
}

variable "destination_instance_system_disk_type" {
  description = "The type of the destination ECS instance system disk"
  type        = string
  default     = "GPSSD"
}

variable "source_server_name" {
  description = "The name of the SMS source server"
  type        = string
}

variable "source_server_os_version" {
  description = "The OS version of the SMS source server"
  type        = string
}

variable "source_server_firmware" {
  description = "The firmware of the SMS source server"
  type        = string
  default     = "BIOS"
}

variable "source_server_boot_loader" {
  description = "The boot loader of the SMS source server"
  type        = string
  default     = "GRUB"
}

variable "source_server_has_rsync" {
  description = "Whether the SMS source server has rsync"
  type        = bool
  default     = true
}

variable "source_server_paravirtualization" {
  description = "Whether the SMS source server is paravirtualization"
  type        = bool
  default     = true
}

variable "source_server_cpu_quantity" {
  description = "The CPU quantity of the SMS source server"
  type        = number
  default     = 2
}

variable "source_server_memory" {
  description = "The memory of the SMS source server"
  type        = number
  default     = 4018196480
}

variable "source_server_agent_version" {
  description = "The agent version of the SMS source server"
  type        = string
}

variable "source_server_disks" {
  description = "The disks of the SMS source server"

  type = object({
    name       = string
    device_use = string
    size       = number
    used_size  = number

    physical_volumes = list(object({
      name        = optional(string)
      device_use  = optional(string)
      file_system = optional(string)
      mount_point = optional(string)
      size        = optional(number)
      used_size   = optional(number)
    }))

    partition_style = optional(string)
    relation_name   = optional(string)
    inode_size      = optional(number)
  })

  default = null
}

variable "migrate_task_type" {
  description = "The type of the SMS migration task"
  type        = string
}

variable "task_auto_start" {
  description = "Whether to automatically start the SMS task"
  type        = bool
  default     = false
}

variable "task_action" {
  description = "The action of the SMS task"
  type        = string
  default     = null
}

variable "task_target_server_disks" {
  description = "The disks of the SMS task target server"
  type        = object({
    name        = string
    size        = number
    device_type = string

    physical_volumes = optional(list(object({
      name         = string
      size         = number
      device_type  = string
      file_system  = string
      mount_point  = string
      volume_index = number
    })), [])
  })

  default = null
}

variable "task_configurations" {
  description = "The configurations of the SMS task"
  type        = list(object({
    config_key   = string
    config_value = string
  }))

  default  = []
  nullable = false
}
