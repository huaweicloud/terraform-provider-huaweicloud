# Variable definitions for authentication
variable "region_name" {
  description = "The region where the RAM resource share is located"
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

# Variable definitions for querying permissions
variable "query_resource_type" {
  description = "The resource type for querying available permissions"
  type        = string
  default     = ""
}

variable "query_permission_type" {
  description = "The type of the permission to query"
  type        = string
  default     = "ALL"
}

variable "query_permission_name" {
  description = "The name of the permission to query"
  type        = string
  default     = ""
}

# Variable definitions for resources
variable "resource_share_id" {
  description = "The ID of the RAM resource share"
  type        = string
}

variable "permission_replace" {
  description = "Whether to replace existing permissions when associating a new permission"
  type        = bool
  default     = false
}
