# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DCS Redis single instance is located"
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
variable "instance_id" {
  description = "The ID of the DCS instance to diagnose"
  type        = string
}

variable "begin_time" {
  description = "The start time of the diagnosis task, in RFC3339 format"
  type        = string
}

variable "end_time" {
  description = "The end time of the diagnosis task, in RFC3339 format"
  type        = string
}
