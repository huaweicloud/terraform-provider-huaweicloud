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

# Variable definitions for bandwidth package
variable "bandwidth_package_name" {
  description = "The name of the bandwidth package"
  type        = string
}

variable "local_area_id" {
  description = "The local area ID"
  type        = string
}

variable "remote_area_id" {
  description = "The remote area ID"
  type        = string
}

variable "charge_mode" {
  description = "Billing option of the bandwidth package"
  type        = string
}

variable "billing_mode" {
  description = "Billing mode of the bandwidth package"
  type        = string
}

variable "bandwidth" {
  description = "Bandwidth in the bandwidth package"
  type        = number
}

variable "bandwidth_package_description" {
  description = "The description about the bandwidth package"
  type        = string
  default     = "Created by Terraform"
}

variable "enterprise_project_id" {
  description = "ID of the enterprise project that the bandwidth package belongs to"
  type        = string
  default     = "0"
}

variable "bandwidth_package_tags" {
  description = "The tags of the bandwidth package"
  type        = map(string)
  default     = {
    "Owner" = "terraform"
    "Env"   = "test"
  }
}
