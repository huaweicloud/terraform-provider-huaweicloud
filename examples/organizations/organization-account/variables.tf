# Variable definitions for authentication
variable "region_name" {
  description = "The region where the organization account is located"
  type        = string
}

variable "access_key" {
  description = "The access key of the IAM user"
  type        = string
}

variable "secret_key" {
  description = "The secret key of the IAM user"
  type        = string
}

# Variable definitions for resources
variable "name" {
  description = "The name of the account"
  type        = string
}

variable "email" {
  description = "The email address of the account"
  type        = string
}

variable "phone" {
  description = "The mobile number of the account"
  type        = string
}

variable "agency_name" {
  description = "The agency name of the account"
  type        = string
}

variable "parent_id" {
  description = "The ID of the root or organization unit in which you want to create a new account"
  type        = string
}

variable "description" {
  description = "The description of the account"
  type        = string
}

variable "tags" {
  description = "The tags of the account"
  type        = map(string)
}
