# Variable definitions for authentication
variable "region_name" {
  description = "The region where the organizational unit will be created"
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
variable "organizational_unit_name" {
  description = "The name of the organizational unit"
  type        = string
}

variable "tags" {
  description = "The tags of the organizational unit"
  type        = map(string)
}
