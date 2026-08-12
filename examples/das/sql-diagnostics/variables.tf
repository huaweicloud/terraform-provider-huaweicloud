# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DAS resources are located"
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

# Variable definitions for huaweicloud_das_batch_set_sql_switch
variable "sql_diagnostics_engine_type" {
  description = "The engine type of the instances"
  type        = string
}

variable "batch_sql_switch_on" {
  description = "Whether to enable the SQL switch"
  type        = bool
}

variable "batch_sql_switch_type" {
  description = "The type of SQL switch to set"
  type        = string
}

variable "batch_sql_instance_ids" {
  description = "The list of instance IDs"
  type        = list(string)
}

variable "batch_sql_retention_hours" {
  description = "The retention hours of the SQL data"
  type        = number
  default     = null
}

# Variable definitions for huaweicloud_das_search_path_switch
variable "search_path_connection_id" {
  description = "The ID of the database connection (DB user ID)"
  type        = string
}

variable "search_path_switch_on" {
  description = "Whether to enable the search path switch"
  type        = bool
}
