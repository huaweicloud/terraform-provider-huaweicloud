# Variable definitions for authentication
variable "region_name" {
  description = "The region where the resources are located"
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

# Variable definitions for resource aggregator
variable "aggregator_name" {
  description = "The name of the resource aggregator"
  type        = string
}

variable "aggregator_type" {
  description = "The type of the resource aggregator, which can be ACCOUNT or ORGANIZATION"
  type        = string
}

variable "account_ids" {
  description = "The list of source account IDs"
  type        = set(string)
  default     = []
}
