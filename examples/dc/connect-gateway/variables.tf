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

# Variable definitions for connect gateway
variable "connect_gateway_name" {
  description = "The name of the connect gateway"
  type        = string
}

variable "connect_gateway_description" {
  description = "The description of the connect gateway"
  type        = string
  default     = "Created by Terraform"
}

variable "address_family" {
  description = "The address family type of the connect gateway"
  type        = string
  default     = "ipv4"
}
