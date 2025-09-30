# Variable definitions for authentication
variable "region_name" {
  description = "The region where the OBS bucket is located"
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
variable "bucket_encryption" {
  description = "The encryption of the OBS bucket"
  type        = bool
  default     = true
}

variable "bucket_encryption_key_id" {
  description = "The encryption key ID of the OBS bucket"
  type        = string
  default     = ""
  nullable    = false
}

variable "key_alias" {
  description = "The alias of the KMS key"
  type        = string
  default     = ""
  nullable    = false

  validation {
    condition     = var.key_alias != "" && var.bucket_encryption && var.bucket_encryption_key_id == ""
    error_message = "The key_alias must be set when bucket_encryption is true and bucket_encryption_key_id is not set."
  }
}

variable "key_usage" {
  description = "The usage of the KMS key"
  type        = string
  default     = "ENCRYPT_DECRYPT"
}

variable "bucket_name" {
  description = "The name of the OBS bucket"
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

variable "bucket_sse_algorithm" {
  description = "The SSE algorithm of the OBS bucket"
  type        = string
  default     = "kms"
}

variable "bucket_force_destroy" {
  description = "The force destroy of the OBS bucket"
  type        = bool
  default     = true
}

variable "bucket_tags" {
  description = "The tags of the OBS bucket"
  type        = map(string)
  default     = {}
}
