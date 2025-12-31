# Variable definitions for authentication
variable "region_name" {
  description = "The region where RGC account is located"
  type        = string
}

variable "access_key" {
  description = "The access key of IAM user"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "The secret key of IAM user"
  type        = string
  sensitive   = true
}

# Variable definitions for account
variable "account_name" {
  description = "The name of RGC account"
  type        = string
}

variable "account_email" {
  description = "The email of RGC account"
  type        = string
}

variable "account_phone" {
  description = "The phone number of RGC account"
  type        = string
  sensitive   = true
  default     = ""
}

variable "identity_store_user_name" {
  description = "The identity store user name of RGC account"
  type        = string
}

variable "identity_store_email" {
  description = "The identity store email of RGC account"
  type        = string
}

variable "parent_organizational_unit_name" {
  description = "The parent organizational unit name of RGC account"
  type        = string
}

variable "parent_organizational_unit_id" {
  description = "The parent organizational unit ID of RGC account"
  type        = string
}
