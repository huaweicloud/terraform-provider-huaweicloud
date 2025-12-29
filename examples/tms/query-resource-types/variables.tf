# Variable definitions for authentication
variable "region_name" {
  description = "The region where the LTS service is located"
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
variable "exact_service_name" {
  description = "The exact service name to filter"
  type        = string
  default     = ""
}

variable "fuzzy_service_name" {
  description = "The fuzzy service name to filter"
  type        = string
  default     = ""
}

variable "fuzzy_resource_type_name" {
  description = "The fuzzy resource type name to filter"
  type        = string
  default     = ""
}
