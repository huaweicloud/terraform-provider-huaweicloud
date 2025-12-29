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
variable "associated_resources_configuration" {
  description = "The configuration of the associated resources"
  type        = list(object({
    type = string
    id   = string
  }))
}

variable "resource_tags" {
  description = "The tags of the resources"
  type        = map(string)
}
