# Variables definitions for authorization
variable "region_name" {
  description = "The region where the scaling configuration is located"
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

# Variables definitions for resource/data source
variable "policy_type" {
  description = "The type of the policy"
  type        = string
  default     = "system"
}

variable "policy_names" {
  description = "The name list of policies to be associated with the user group"
  type        = list(string)
}

variable "group_name" {
  description = "The name of the user group"
  type        = string
}

variable "group_description" {
  description = "The description of the user group"
  type        = string
  default     = ""
}
