variable "availability_zone" {
  type        = string
  default     = ""
  description = "The availability zone to which the cloud desktop flavor and network belong."
}

variable "desktop_flavor" {
  type        = string
  default     = ""
  description = "The flavor name of the cloud desktop."
}

variable "desktop_cpu_core_number" {
  type        = number
  default     = 4
  description = "The number of the cloud desktop CPU cores."
}

variable "desktop_memory" {
  type        = number
  default     = 8
  description = "The number of the cloud desktop memories."
}

variable "desktop_os_type" {
  type        = string
  default     = "Windows"
  description = "The OS type of the cloud desktop."
}

variable "desktop_image_type" {
  type        = string
  default     = "market"
  description = "The specified image type that the cloud desktop used."
}

variable "desktop_image_id" {
  type        = string
  description = "The specified image ID that the cloud desktop used."
}

variable "vpc_name" {
  type        = string
  description = "The VPC name."
}

variable "subnet_name" {
  type        = string
  description = "The subnet name."
}

variable "security_group_name" {
  type        = string
  description = "The security group name."
}

variable "desktop_user_name" {
  type        = string
  description = "The user name that the cloud desktop used."
}

variable "desktop_user_email" {
  type        = string
  description = "The email address that the user used."
}

variable "cloud_desktop_name" {
  type        = string
  description = "The cloud desktop name."
}

variable "desktop_user_group_name" {
  type        = string
  default     = "users"
  description = "The name of the user group that cloud desktop used."
}

variable "desktop_root_volume_type" {
  type        = string
  default     = "SSD"
  description = "The storage type of system disk."
}

variable "desktop_root_volume_size" {
  type        = number
  default     = 100
  description = "The storage capacity of system disk."
}

variable "desktop_data_volumes" {
  type = list(object({
    type = string
    size = number
  }))
  default = [
    {type="SSD", size=100},
  ]
  description = "The storage configuration of data disks."
}
