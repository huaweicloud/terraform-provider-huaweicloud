# Provider Configuration
variable "region_name" {
  description = "The name of the region to deploy resources"
  type        = string
}

variable "access_key" {
  description = "The access key of HuaweiCloud"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "The secret key of HuaweiCloud"
  type        = string
  sensitive   = true
}

# Variable definitions for resources/data sources
variable "migration_project_name" {
  description = "The migration project name"
  type        = string
}

variable "migration_project_region" {
  description = "The region name"
  type        = string
}

variable "migration_project_use_public_ip" {
  description = "Whether to use a public IP address for migration"
  type        = bool
}

variable "migration_project_exist_server" {
  description = "Whether the server already exists"
  type        = bool
}

variable "migration_project_type" {
  description = "The migration project typ"
  type        = string
}

variable "migration_project_syncing" {
  description = "Whether to continue syncing after the first copy or sync"
  type        = bool
}
