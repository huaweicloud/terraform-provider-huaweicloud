# Variable definitions for authentication
variable "region_name" {
  description = "The region where the KMS key is located"
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
variable "key_name" {
  description = "The alias name of the KMS key"
  type        = string
}

variable "key_algorithm" {
  description = "The generation algorithm of the KMS key"
  type        = string
  default     = "AES_256"
}

variable "key_usage" {
  description = "The usage of the KMS key"
  type        = string
  default     = "ENCRYPT_DECRYPT"
}

variable "key_source" {
  description = "The source of the KMS key"
  type        = string
  default     = "kms"
}

variable "key_description" {
  description = "The description of the KMS key"
  type        = string
  default     = ""
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project to which the KMS key belongs"
  type        = string
  default     = null
}

variable "key_tags" {
  description = "The key/value pairs to associate with the KMS key"
  type        = map(string)
  default     = {}
}

variable "key_schedule_time" {
  description = "The number of days after which the KMS key is scheduled to be deleted"
  type        = string
  default     = "7"
}
