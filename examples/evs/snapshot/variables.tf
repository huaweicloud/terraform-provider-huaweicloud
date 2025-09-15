# Variable definitions for authentication
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
  default     = "SAS"
}

variable "volume_size" {
  description = "The size of the volume"
  type        = number
  default     = 20
}

variable "voluem_description" {
  description = "The description of the volume"
  type        = string
  default     = ""
}

variable "vouleme_multiattach" {
  description = "The volume is shared volume or not"
  type        = bool
  default     = false
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

variable "snapshot_metadata" {
  description = "The metadata information of the snapshot"
  type        = map(string)
  default     = {}
}

variable "snapshot_force" {
  description = "The flag for forcibly creating a snapshot"
  type        = bool
  default     = false
}
