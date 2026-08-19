# Variable definitions for authentication
variable "region_name" {
  description = "The region where the VOD service is located"
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

# Variable definitions for VOD resources
variable "media_category_name" {
  description = "The name of the media category"
  type        = string
}

variable "media_asset_name" {
  description = "The name of the media asset"
  type        = string
}

variable "media_asset_url" {
  description = "The HTTP or HTTPS URL of the media source file"
  type        = string
}

variable "media_asset_description" {
  description = "The description of the media asset"
  type        = string
  default     = ""
}

variable "media_asset_labels" {
  description = "The labels of the media asset, separated by commas"
  type        = string
  default     = "tf_label_1,tf_label_2"
}
