# Variable definitions for authentication
variable "region_name" {
  description = "The name of the VPC"
  type        = string
  default     = "cn-north-4"
}

variable "access_key" {
  description = "The access key of the IAM user"
  type        = string
}

variable "secret_key" {
  description = "The secret key of the IAM user"
  type        = string
}
