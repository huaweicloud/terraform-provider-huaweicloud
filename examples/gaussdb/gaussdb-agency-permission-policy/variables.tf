# Variable definitions for authentication
variable "region_name" {
  description = "The region where the GaussDB instance is located"
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
variable "agency_name" {
  description = "The agency name. Currently, only RDSAccessProjectResource is supported"
  type        = string
  default     = "RDSAccessProjectResource"
}

variable "bind_role_names" {
  description = "The permission policies to be bound to the agency"
  type        = set(string)
}

variable "unbind_role_names" {
  description = "The permission policies to be unbound from the agency"
  type        = set(string)
}
