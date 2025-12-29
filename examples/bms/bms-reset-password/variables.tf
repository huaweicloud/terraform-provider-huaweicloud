# Variable definitions for authentication
variable "region_name" {
  description = "The region where the BMS instance is located"
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
variable "bms_instance_id" {
  description = "The ID of the BMS instance"
  type        = string
}

variable "bms_instance_new_password" {
  description = "The new password of the BMS instance"
  type        = string
  sensitive   = true
}
