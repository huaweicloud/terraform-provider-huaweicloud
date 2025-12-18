# Variable definitions for authentication
variable "region_name" {
  description = "The region where the CSMS secret is located"
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
variable "secret_name" {
  description = "The name of the secret"
  type        = string
}

variable "secret_text" {
  description = "The plaintext of a text secret"
  type        = string
  sensitive   = true
}

variable "secret_type" {
  description = "The type of the secret"
  type        = string
  default     = "COMMON"
}

variable "kms_key_id" {
  description = "The ID of the KMS key used to encrypt the secret"
  type        = string
  default     = ""
}

variable "secret_description" {
  description = "The description of the secret"
  type        = string
  default     = ""
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project to which the secret belongs"
  type        = string
  default     = null
}

variable "secret_tags" {
  description = "The key/value pairs to associate with the secret"
  type        = map(string)
  default     = {}
}
