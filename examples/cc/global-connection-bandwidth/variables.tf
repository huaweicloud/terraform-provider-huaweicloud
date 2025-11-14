# Variable definitions for authentication
variable "region_name" {
  description = "The region where the resources are located"
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

# Variable definitions for global connection bandwidth
variable "global_connection_bandwidth_name" {
  description = "The name of the global connection bandwidth"
  type        = string
}

variable "bandwidth_type" {
  description = "The type of the global connection bandwidth"
  type        = string
}

variable "bordercross" {
  description = "Whether the global connection bandwidth crosses borders"
  type        = bool
}

variable "bandwidth_size" {
  description = "Bandwidth size of the global connection bandwidth"
  type        = number
}

variable "charge_mode" {
  description = "Billing option of the global connection bandwidth"
  type        = string
}

variable "global_connection_bandwidth_description" {
  description = "The description about the global connection bandwidth"
  type        = string
  default     = "Created by Terraform"
}

variable "global_connection_bandwidth_tags" {
  description = "The tags of the global connection bandwidth"
  type        = map(string)
  default     = {
    "Owner" = "terraform"
    "Env"   = "test"
  }
}
