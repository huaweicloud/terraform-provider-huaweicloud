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
variable "volume_image_id" {
  description = "The ID of the image used to create the volume, if not specified, the first available image matching the criteria will be used"
  type        = string
  default     = ""
}

variable "volume_image_visibility" {
  description = "The visibility of the volume image"
  type        = string
  default     = "public"
}

variable "volume_image_os" {
  description = "The OS of the volume image"
  type        = string
  default     = "Ubuntu"
}

variable "volume_type" {
  description = "The type of the volume"
  type        = string
  default     = "GPSSD"
}

variable "volume_availability_zone" {
  description = "The availability zone for the volume"
  type        = string
  default     = ""
  nullable    = false
}

variable "volume_description" {
  description = "The description of the volume"
  type        = string
  default     = ""
}

variable "volume_metadata" {
  description = "The metadata of the volume"
  type        = map(string)
  default     = {}
}

variable "volume_multiattach" {
  description = "The volume is shared volume or not"
  type        = bool
  default     = false
}

variable "volume_name" {
  description = "The name of the volume"
  type        = string
}

variable "volume_size" {
  description = "The size of the volume"
  type        = number
  default     = 40
}

variable "volume_tags" {
  description = "The tags of the volume"
  type        = map(string)
  default     = {}
}

variable "snapshot_name" {
  description = "The name of the snapshot"
  type        = string
}

variable "snapshot_metadata" {
  description = "The metadata information of the snapshot"
  type        = map(string)
  default     = {}
}

variable "snapshot_description" {
  description = "The description of the snapshot"
  type        = string
  default     = ""
}
