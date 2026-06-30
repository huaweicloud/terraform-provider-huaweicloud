# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DataArts Studio instance is located"
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

# Variable definitions for DataArts Studio instance and workspace
variable "workspace_id" {
  description = "The ID of the workspace"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_id" {
  description = "The ID of the DataArts Studio instance to which the workspace belongs"
  type        = string
  default     = ""
  nullable    = false
}

variable "workspace_name" {
  description = "The name of the workspace used to filter results"
  type        = string
  default     = ""
  nullable    = false
}

# Variable definitions for workspace user
variable "user_id" {
  description = "The ID of the IAM user to be added to the workspace"
  type        = string
  default     = ""
  nullable    = false
}

variable "user_name" {
  description = "The name of the IAM user to be added to the workspace"
  type        = string
  default     = ""
  nullable    = false
}

variable "role_ids" {
  description = "The role ID list of the workspace user"
  type        = list(string)
  default     = []
  nullable    = false
}
