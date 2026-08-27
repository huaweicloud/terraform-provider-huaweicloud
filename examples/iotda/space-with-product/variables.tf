# Variable definitions for authentication
variable "region_name" {
  description = "The region where the IoTDA service is located"
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

# Variable definitions for IoTDA resources
variable "space_name" {
  description = "The name of the resource space"
  type        = string
}

variable "product_name" {
  description = "The name of the product"
  type        = string
}

variable "iotda_access_address" {
  description = "The HTTPS application access address of the IoTDA instance"
  type        = string
}
