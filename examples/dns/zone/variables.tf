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

variable "dns_public_zone_name" {
  description = "The name of the zone"
  type        = string
}

variable "dns_public_zone_email" {
  description = "The email address of the administrator managing the zone"
  type        = string
  default     = ""
}

variable "dns_public_zone_type" {
  description = "The type of zone"
  type        = string
  default     = "public"
}

variable "dns_public_zone_description" {
  description = "The description of the zone"
  type        = string
}

variable "dns_public_zone_ttl" {
  description = "The time to live (TTL) of the zone"
  type        = number
  default     = 300
}

variable "dns_public_zone_enterprise_project_id" {
  description = "The enterprise project ID of the zone"
  type        = string
  default     = ""
}

variable "dns_public_zone_status" {
  description = "The status of the zone"
  type        = string
  default     = "ENABLE"
}

variable "dns_public_zone_dnssec" {
  description = "Whether to enable DNSSEC for a public zone"
  type        = string
  default     = "DISABLE"
}

variable "dns_public_zone_router" {
  description = "The list of the router of the zone"
  type = list(object({
    router_id     = string
    router_region = string
  }))
  default = []
}
