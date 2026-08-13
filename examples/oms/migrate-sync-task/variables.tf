# Variable definitions for authentication
variable "region_name" {
  description = "The region where the OMS migration synchronization task is located"
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
variable "source_bucket_name" {
  description = "The name of the source OBS bucket"
  type        = string
}

variable "bucket_storage_class" {
  description = "The storage class of the OBS bucket"
  type        = string
  default     = "STANDARD"
}

variable "bucket_acl" {
  description = "The ACL of the OBS bucket"
  type        = string
  default     = "private"
}

variable "bucket_force_destroy" {
  description = "Whether to force destroy the OBS bucket"
  type        = bool
  default     = true
}

variable "dest_bucket_name" {
  description = "The name of the destination OBS bucket"
  type        = string
}

variable "source_cloud_type" {
  description = "The source cloud service provider"
  type        = string
  default     = "HuaweiCloud"
}

variable "source_region" {
  description = "The region where the source bucket is located"
  type        = string
}

variable "source_access_key" {
  description = "The access key for accessing the source bucket"
  type        = string
  sensitive   = true
}

variable "source_secret_key" {
  description = "The secret key for accessing the source bucket"
  type        = string
  sensitive   = true
}

variable "dest_access_key" {
  description = "The access key for accessing the destination bucket"
  type        = string
  sensitive   = true
}

variable "dest_secret_key" {
  description = "The secret key for accessing the destination bucket"
  type        = string
  sensitive   = true
}

variable "task_description" {
  description = "The description of the migration synchronization task"
  type        = string
  default     = ""
}

variable "consistency_check" {
  description = "The consistency check method"
  type        = string
  default     = "size_last_modified"
}

variable "enable_metadata_migration" {
  description = "Whether to enable metadata migration"
  type        = bool
  default     = false
}
