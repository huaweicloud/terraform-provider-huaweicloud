# Variable definitions for authentication
variable "region_name" {
  description = "The region where the ServiceStage service is located"
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

# Variable definitions for VPC resources
variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
}

variable "eip_bandwidth_name" {
  description = "The name of the dedicated bandwidth of the EIP"
  type        = string
}

# Variable definitions for ServiceStage resources
variable "environment_name" {
  description = "The name of the environment"
  type        = string
}

variable "environment_description" {
  description = "The description of the environment"
  type        = string
  default     = ""
}
