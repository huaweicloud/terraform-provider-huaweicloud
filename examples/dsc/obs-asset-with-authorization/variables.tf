# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DSC OBS asset is located"
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

# Variable definitions for resources
variable "bucket_name" {
  description = "The name of the OBS bucket to be added as a DSC asset"
  type        = string
}

variable "asset_name" {
  description = "The name of the DSC OBS asset"
  type        = string
}
