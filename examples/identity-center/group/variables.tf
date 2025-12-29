# Variable definitions for authentication
variable "region_name" {
  description = "The region where the LTS service is located"
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
variable "is_instance_create" {
  description = "Whether to create the identity center instance"
  type        = bool
  default     = true
}

variable "is_region_need_register" {
  description = "Whether to register the region"
  type        = bool
  default     = true
}

variable "instance_store_id_alias" {
  description = "The alias of the identity center instance"
  type        = string
  default     = ""
}

variable "group_name" {
  description = "The name of the identity center group"
  type        = string
}

variable "group_description" {
  description = "The description of the identity center group"
  type        = string
  default     = ""
}
