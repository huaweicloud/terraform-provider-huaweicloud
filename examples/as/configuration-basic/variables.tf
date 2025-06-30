variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "kps_key_pair_name" {
  description = "The name of the key pair"
  type        = string
}

variable "kps_public_key" {
  description = "The public key for SSH access"
  type        = string
}

variable "as_configuration_name" {
  description = "The name of the AS configuration"
  type        = string
}

variable "as_metadata" {
  description = "The metadata for the instance"
  type        = map(string)
  default     = {}
}

variable "as_user_data" {
  description = "The user data script for instance initialization"
  type        = string
  default     = ""
}

variable "as_disks" {
  description = "The disk configurations for the instance"
  type = list(object({
    size        = number
    volume_type = string
    disk_type   = string
  }))

  nullable = false
}

variable "as_public_ip" {
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

  default  = []
  nullable = false

  validation {
    condition     = length(var.as_public_ip) <= 1
    error_message = "The maximum length of as_public_ip is 1"
  }
}
