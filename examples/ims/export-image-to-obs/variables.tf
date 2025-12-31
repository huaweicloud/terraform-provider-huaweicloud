# Variable definitions for authentication
variable "region_name" {
  description = "The region where the CDN domain is located"
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

# Variable definitions for image filtering
variable "image_type" {
  description = "The type of images to filter (ECS, BMS, etc.)"
  type        = string
  default     = ""
}

variable "image_os" {
  description = "The OS of images to filter (Ubuntu, CentOS, etc.)"
  type        = string
  default     = ""
}

variable "image_name_regex" {
  description = "The regex pattern to filter images by name"
  type        = string
  default     = ""
}

# Variable definitions for OBS bucket
variable "obs_bucket_name" {
  description = "The name of the OBS bucket for storing exported images"
  type        = string
}

variable "obs_bucket_tags" {
  description = "The tags of the OBS bucket"
  type        = map(string)
  default     = {}
}

# Variable definitions for image export
variable "file_format" {
  description = "The file format of the exported image (vhd, zvhd, vmdk, raw, qcow2, zvhd2, vdi)"
  type        = string
  default     = "zvhd2"
}
