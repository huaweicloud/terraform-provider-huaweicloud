# Variable definitions for authentication
variable "region_name" {
  description = "The region where the Kafka instance is located"
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

variable "target_account_access_key" {
  description = "The access key of the IAM user in the target account"
  type        = string
  sensitive   = true
}

variable "target_account_secret_key" {
  description = "The secret key of the IAM user in the target account"
  type        = string
  sensitive   = true
}

# Variable definitions for resources/data sources
variable "main_domain_name" {
  description = "The name of the main domain"
  type        = string
}

variable "sub_domain_prefix" {
  description = "The prefix of the sub-domain"
  type        = string
}

variable "recordset_type" {
  description = "The type of the recordset"
  type        = string
  default     = "TXT"
}

variable "recordset_ttl" {
  description = "The time to live (TTL) of the recordset"
  type        = number
  default     = 300
}
