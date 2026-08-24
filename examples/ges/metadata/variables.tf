# Variable definitions for authentication
variable "region_name" {
  description = "The region where the GES metadata is located"
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

# Variable definitions for huaweicloud_obs_bucket
variable "bucket_name" {
  description = "The OBS bucket name for storing GES metadata schema files"
  type        = string
}

# Variable definitions for huaweicloud_ges_metadata
variable "metadata_name" {
  description = "The GES metadata name"
  type        = string
}

variable "metadata_description" {
  description = "The description of the GES metadata"
  type        = string
}

variable "metadata_schema_file" {
  description = "The schema file name in the OBS bucket"
  type        = string
}

variable "metadata_properties" {
  description = "The properties of the GES metadata label"
  type        = list(any)
  nullable    = false
}
