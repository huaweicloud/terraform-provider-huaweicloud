# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DNS service is located"
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

variable "name" {
  description = "The name of the zone"
  type        = string
}

variable "description" {
  description = "The description of the zone"
  type        = string
}

variable "ttl" {
  description = "The time to live (TTL) of the zone"
  type        = number
}

variable "dnssec" {
  description = "Whether to enable DNSSEC for a public zone"
  type        = string
}
