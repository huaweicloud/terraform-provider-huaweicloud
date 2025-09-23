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

# Volume related variables
variable "availability_zone" {
  description = "The availability zone to which the volume belongs"
  type        = string
  default     = ""
}

variable "volume_type" {
  description = "The type of the volume"
  type        = string
}

variable "volume_name" {
  description = "The name of the volume"
  type        = string
  default     = ""
}

variable "volume_size" {
  description = "The size of the volume in GB"
  type        = number
}

variable "volume_device_type" {
  description = "The device type of the volume"
  type        = string
  default     = "VBD"
}

# Vault related variables
variable "name" {
  description = "The name of the CBR vault"
  type        = string
}

variable "type" {
  description = "The type of the CBR vault"
  type        = string
  default     = "disk"
}

variable "protection_type" {
  description = "The protection type of the vault"
  type        = string
  default     = "backup"
}

variable "size" {
  description = "The size of the CBR vault in GB"
  type        = number
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project to which the vault belongs"
  type        = string
  default     = "0"
}
