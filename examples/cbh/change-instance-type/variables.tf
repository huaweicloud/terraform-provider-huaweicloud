# Variable definitions for authentication
variable "region_name" {
  description = "The region where the CBH instance is located"
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

variable "server_id" {
  description = "The ID of the single node CBH instance to change type"
  type        = string
}

variable "availability_zone" {
  description = "The availability zone of the single-node CBH instance to change type"
  type        = string
  default     = ""
}
