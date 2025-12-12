# Variable definitions for authentication
variable "region_name" {
  description = "The region where the workspace is located"
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
variable "workspace_name" {
  description = "The name of the workspace"
  type        = string
}

variable "workspace_project_name" {
  description = "The name of the project to in which to create the workspace"
  type        = string
}

variable "workspace_description" {
  description = "The description of the workspace"
  type        = string
  default     = ""
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project to which the workspace belongs"
  type        = string
  default     = null
}

variable "workspace_tags" {
  description = "The key/value pairs to associate with the workspace"
  type        = map(string)
  default     = {}
}
