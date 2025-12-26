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

# Variable definitions for global gateway
variable "global_gateway_name" {
  description = "The name of the global gateway"
  type        = string
}

variable "global_gateway_description" {
  description = "The description of the global gateway"
  type        = string
  default     = "Created by Terraform"
}

variable "address_family" {
  description = "The IP address family of the global DC gateway"
  type        = string
  default     = "ipv4"
}

variable "bgp_asn" {
  description = "The BGP ASN of the global DC gateway"
  type        = number
}

variable "enterprise_project_id" {
  description = "The enterprise project ID that the global DC gateway belongs to"
  type        = string
  default     = "0"
}

variable "global_gateway_tags" {
  description = "The tags of the global gateway"
  type        = map(string)
  default     = {
    "Owner" = "terraform"
    "Env"   = "test"
  }
}
