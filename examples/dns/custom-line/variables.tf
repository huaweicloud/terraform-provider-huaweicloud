# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DNS service is located"
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

variable "dns_custom_line_name" {
  description = "The custom line name"
  type        = string
}

variable "dns_custom_line_ip_segments" {
  description = "The IP address range"
  type        = list(string)
}
