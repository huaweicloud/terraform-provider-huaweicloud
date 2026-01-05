# Variable definitions for authentication
variable "region_name" {
  description = "The region where the OMS sync task is located"
  type        = string
}

variable "access_key" {
  description = "The access key for the HuaweiCloud provider (can use source bucket access_key)"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "The secret key for the HuaweiCloud provider (can use source bucket secret_key)"
  type        = string
  sensitive   = true
}

# Variable definitions for resources
variable "buckets" {
  description = "The configurations of OBS buckets for sync task"
  type        = list(object({
    name          = string
    role          = string
    region        = string
    access_key    = string
    secret_key    = string
    storage_class = optional(string)
    acl           = optional(string)
    force_destroy = optional(bool)
    tags          = optional(map(string))
  }))

  validation {
    condition     = alltrue([for bucket in var.buckets : contains(["source", "destination"], bucket.role)])
    error_message = "The bucket role must be either 'source' or 'destination'."
  }

  validation {
    condition     = length([for bucket in var.buckets : bucket if bucket.role == "source"]) == 1
    error_message = "Exactly one bucket with role 'source' must be specified."
  }

  validation {
    condition     = length([for bucket in var.buckets : bucket if bucket.role == "destination"]) == 1
    error_message = "Exactly one bucket with role 'destination' must be specified."
  }
}

variable "source_object_configurations" {
  description = "The configurations of objects to be uploaded to the source bucket"
  type        = list(object({
    key          = string
    content      = string
    content_type = optional(string)
  }))
  default     = []
}

variable "sync_task_description" {
  description = "The description of the OMS migration sync task"
  type        = string
  default     = ""
}

variable "sync_task_enable_kms" {
  description = "Whether to enable KMS for the OMS migration sync task"
  type        = bool
  default     = false
}

variable "sync_task_enable_restore" {
  description = "Whether to enable restore for the OMS migration sync task"
  type        = bool
  default     = false
}

variable "sync_task_enable_metadata_migration" {
  description = "Whether to enable metadata migration for the OMS migration sync task"
  type        = bool
  default     = true
}

variable "sync_task_consistency_check" {
  description = "The consistency check method for the OMS migration sync task"
  type        = string
  default     = "crc64"

  validation {
    condition     = contains(["crc64", "no_check"], var.sync_task_consistency_check)
    error_message = "The sync_task_consistency_check must be one of: crc64, no_check."
  }
}

variable "sync_task_action" {
  description = "The action of the OMS migration sync task"
  type        = string
  default     = "start"
}

variable "source_cdn_configuration" {
  description = "The CDN configuration for the source bucket"
  type        = object({
    domain              = string
    protocol            = string
    authentication_type = optional(string)
    authentication_key  = optional(string)
  })
  default     = null
}
