# Variable definitions for authentication
variable "region_name" {
  description = "The region where the protection group is located"
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
variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "192.168.0.0/16"
}

variable "protection_group_name" {
  description = "The name of the protection group"
  type        = string
}

variable "source_availability_zone" {
  description = "The production site AZ of the protection group"
  type        = string
  default     = ""

  validation {
    condition     = (
      (var.source_availability_zone == "" && var.target_availability_zone == "") ||
      (var.source_availability_zone != "" && var.target_availability_zone != "")
    )
    error_message = "Both `source_availability_zone` and `target_availability_zone` must be set, or both must be empty"
  }
}

variable "target_availability_zone" {
  description = "The disaster recovery site AZ of the protection group"
  type        = string
  default     = ""
}

variable "protection_group_dr_type" {
  description = "The deployment model"
  type        = string
  default     = null
}

variable "protection_group_description" {
  description = "The description of the protection group"
  type        = string
  default     = null
}

variable "protection_group_enable" {
  description = "Whether enable the protection group start protecting"
  type        = bool
  default     = null
}
