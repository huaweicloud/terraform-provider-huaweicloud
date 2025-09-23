# Variable definitions for authentication
variable "region_name" {
  description = "The region where the VPC peering connection is located"
  type        = string
}

variable "access_key" {
  description = "The access key of the IAM user"
  type        = string
}

variable "secret_key" {
  description = "The secret key of the IAM user"
  type        = string
}

# Variable definitions for resources/data sources
variable "vpc_configurations" {
  description = "The list of VPC configurations for peering connection"
  type        = list(object({
    vpc_name              = string
    vpc_cidr              = string
    subnet_name           = string
    enterprise_project_id = optional(string, null)
  }))

  validation {
    condition     = length(var.vpc_configurations) == 2
    error_message = "Exactly 2 VPC configurations are required for peering connection."
  }
}

variable "peering_connection_name" {
  description = "The name of the VPC peering connection"
  type        = string
}
