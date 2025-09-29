# Variable definitions for authentication
variable "region_name" {
  description = "The region where the Organizations service is located"
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

variable "enabled_policy_types" {
  description = "The list of Organizations policy types to enable in the Organization Root"
  type        = list(string)
}

variable "root_tags" {
  description = "The key/value to attach to the root"
  type        = map(string)
}
