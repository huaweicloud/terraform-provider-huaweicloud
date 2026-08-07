# Variable definitions for authentication
variable "region_name" {
  description = "The region where resources will be created"
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

# Variable definitions for DRS connection
variable "connection_name" {
  description = "The DRS connection name"
  type        = string
}

variable "description" {
  description = "The description of the DRS connection"
  type        = string
  default     = ""
}

variable "endpoint_ip" {
  description = "The IP address and port of the primary MongoDB database, e.g. 192.168.0.1:8080"
  type        = string
}

variable "db_user" {
  description = "The database username"
  type        = string
  default     = "mog"
}

variable "db_password" {
  description = "The password for the MongoDB database user"
  type        = string
  sensitive   = true
}

variable "db_name" {
  description = "The database name"
  type        = string
  default     = "root"
}

variable "shard1_ip" {
  description = "The IP address and port of the first MongoDB shard, e.g. 192.168.0.1:8000"
  type        = string
}

variable "shard2_ip" {
  description = "The IP address and port of the second MongoDB shard, e.g. 192.168.0.2:8000"
  type        = string
}

variable "driver_name" {
  description = "The driver name of the connection configuration"
  type        = string
  default     = "mongodb"
}
