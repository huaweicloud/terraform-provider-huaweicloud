# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DDS backup is located"
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
variable "template_name" {
  description = "The name of the parameter template"
  type        = string
}

variable "template_mapping" {
  description = "The mapping between parameter names and parameter values"
  type        = map(string)
}

variable "template_node_type" {
  description = "The node type of the parameter template"
  type        = string
  default     = "mongos"
}

variable "database_version" {
  description = "The database version"
  type        = string
  default     = "4.0"
}

variable "template_description" {
  description = "The description of the parameter template"
  type        = string
  default     = ""
}
