# Variable definitions for authentication
variable "region_name" {
  description = "The region where the SecMaster resources will be created"
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

# Variable definitions for SecMaster workspace
variable "workspace_name" {
  description = "The name of the SecMaster workspace"
  type        = string
}

variable "workspace_description" {
  description = "The description of the SecMaster workspace"
  type        = string
  default     = "Created by Terraform"
}

# Variable definitions for SecMaster dataspace
variable "dataspace_name" {
  description = "The name of the SecMaster dataspace. The name can only contain English letters, digits and hyphens (-), and cannot start or end with a hyphen (-), nor can they appear consecutively. Valid length: 5-63"
  type        = string
}

variable "dataspace_description" {
  description = "The description of the SecMaster dataspace"
  type        = string
  default     = "Created by Terraform"
}
