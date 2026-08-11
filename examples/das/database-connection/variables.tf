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

# Variable definitions for huaweicloud_das_database_instance_connection
variable "connection_instance_id" {
  description = "The ID of the RDS instance to connect"
  type        = string
}

variable "connection_engine_type" {
  description = "The engine type of the database instance"
  type        = string
}

variable "connection_network_type" {
  description = "The network type of the database instance connection"
  type        = string
}

variable "connection_username" {
  description = "The username for the database instance connection"
  type        = string
}

variable "connection_password" {
  description = "The password for the database instance connection"
  type        = string
  sensitive   = true
}

variable "connection_is_save_password" {
  description = "Whether to save the password for the database instance connection"
  type        = bool
  default     = true
  sensitive   = true
}

variable "connection_port" {
  description = "The port of the database instance connection"
  type        = number
  default     = null
}

variable "connection_database_name" {
  description = "The database name of the database instance connection"
  type        = string
  default     = null
}

variable "connection_sql_record_flag" {
  description = "Whether SQL recording is enabled for the database instance connection"
  type        = bool
  default     = null
}

variable "connection_description" {
  description = "The description of the database instance connection"
  type        = string
  default     = null
}

variable "connection_node_ids" {
  description = "The unique identifiers of the instance nodes"
  type        = list(string)
  default     = []
  nullable    = false
}

# Variable definitions for huaweicloud_das_database_user
variable "db_user_name" {
  description = "The name of the database user"
  type        = string
}

variable "db_user_password" {
  description = "The password of the database user"
  type        = string
  sensitive   = true
}

# Variable definitions for huaweicloud_das_shared_connection
variable "shared_user_id" {
  description = "The IAM user ID to share the connection with"
  type        = string
}

variable "shared_user_name" {
  description = "The IAM user name to share the connection with"
  type        = string
}

variable "shared_expired_at" {
  description = "The expiration time of the shared connection, in RFC3339 format"
  type        = string
  default     = null
}
