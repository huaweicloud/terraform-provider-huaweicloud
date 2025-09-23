# Variables definitions for authorization
variable "region_name" {
  description = "The region where the AS configuration is located"
  type        = string
}

variable "access_key" {
  description = "The access key of the IAM user"
  type        = string
}

variable "secret_key" {
  description = "The secret key of the IAM user"
  type        = string
}

# Variables definitions for resource/data source
variable "availability_zone" {
  description = "The availability zone to which the AS configuration belongs"
  type        = string
  default     = ""
}

variable "configuration_flavor_id" {
  description = "The flavor ID of the AS configuration"
  type        = string
  default     = ""
}

variable "configuration_flavor_performance_type" {
  description = "The performance type of the AS configuration"
  type        = string
  default     = "normal"
}

variable "configuration_flavor_cpu_core_count" {
  description = "The CPU core count of the AS configuration"
  type        = number
  default     = 2
}

variable "configuration_flavor_memory_size" {
  description = "The memory size of the AS configuration"
  type        = number
  default     = 4
}

variable "configuration_image_id" {
  description = "The image ID of the AS configuration"
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
  description = "The name of the AS configuration"
  type        = string
}

variable "configuration_metadata" {
  description = "The metadata for the instance"
  type        = map(string)
}

variable "configuration_user_data" {
  description = "The user data script for instance initialization"
  type        = string
}

variable "configuration_disks" {
  description = "The disk configurations for the instance"
  type = list(object({
    size        = number
    volume_type = string
    disk_type   = string
  }))

  nullable = false
}

variable "configuration_public_ip" {
  description = "The public IP configuration for the instance"

  type = list(object({
    eip = object({
      ip_type = string
      bandwidth = object({
        size          = number
        share_type    = string
        charging_mode = string
      })
    })
  }))

  nullable = false
}
