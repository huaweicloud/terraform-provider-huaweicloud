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

# Variable definitions for central network
variable "central_network_name" {
  description = "The name of the central network"
  type        = string
}

variable "central_network_description" {
  description = "The description about the central network"
  type        = string
  default     = "Created by Terraform"
}

variable "enterprise_project_id" {
  description = "ID of the enterprise project that the central network belongs to"
  type        = string
  default     = "0"
}

variable "central_network_tags" {
  description = "The tags of the central network"
  type        = map(string)
  default     = {
    "Owner" = "terraform"
    "Env"   = "test"
  }
}
