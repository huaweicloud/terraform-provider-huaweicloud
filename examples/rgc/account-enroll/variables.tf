# Variable definitions for authentication
variable "region_name" {
  description = "The region where the RGC resources will be created"
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

# Variable definitions for organizational unit
variable "organizational_unit_name" {
  description = "The name of the organizational unit to be created (required if create_organizational_unit is true)"
  type        = string
  default     = ""
}

variable "parent_organizational_unit_id" {
  description = "The ID of the parent organizational unit. Required for account enrollment and OU creation"
  type        = string
}

# Variable definitions for blueprint example
variable "blueprint_managed_account_id" {
  description = "The ID of the account to be enrolled with blueprint configuration"
  type        = string
}

variable "create_organizational_unit" {
  description = "Whether to create a new organizational unit. If false, use existing parent_organizational_unit_id"
  type        = bool
  default     = true
}

variable "blueprint_product_id" {
  description = "The ID of the blueprint product"
  type        = string
}

variable "blueprint_product_version" {
  description = "The version of the blueprint product"
  type        = string
}

variable "blueprint_variables" {
  description = "The variables for the blueprint configuration (JSON string format)"
  type        = string
}

variable "is_blueprint_has_multi_account_resource" {
  description = "Whether the blueprint has multi-account resources"
  type        = bool
}
