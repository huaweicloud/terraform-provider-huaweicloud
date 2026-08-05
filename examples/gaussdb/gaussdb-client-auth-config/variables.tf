# Variable definitions for authentication
variable "region_name" {
  description = "The region where the GaussDB instance is located"
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

# Variable definitions for resources
variable "instance_id" {
  description = "The ID of the GaussDB instance"
  type        = string
  default     = ""
}

variable "config_type" {
  description = "The client connection type. Valid values are host, hostssl, and hostnossl"
  type        = string
  default     = "host"
}

variable "config_database" {
  description = "The database name that the record matches. The value can be all or an existing database name"
  type        = string
  default     = "all"
}

variable "config_user" {
  description = "The database user that the record matches. The value can be all or an existing username"
  type        = string
  default     = "root"
}

variable "config_address" {
  description = "The IP address range that the record matches. The value must be in CIDR format"
  type        = string
  default     = "10.10.0.0/16"
}

variable "config_method" {
  description = "The authentication method. Valid values include md5, sha256, sm3, reject, cert, etc"
  type        = string
  default     = "md5"
}
