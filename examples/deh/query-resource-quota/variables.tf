# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DEH service is located"
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

# Variable definitions for data sources
variable "host_type" {
  description = "The type of the dedicated host"
  type        = string
  default     = ""
}
