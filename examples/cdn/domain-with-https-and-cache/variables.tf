# Variable definitions for authentication
variable "region_name" {
  description = "The region where the CDN domain is located"
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
variable "domain_name" {
  description = "The name of the CDN domain to be accelerated"
  type        = string
}

variable "domain_type" {
  description = "The business type of the domain"
  type        = string
  default     = "web"

  validation {
    condition     = contains(["web", "download", "video", "wholeSite"], var.domain_type)
    error_message = "The domain_type must be one of: web, download, video, wholeSite."
  }
}

variable "service_area" {
  description = "The area covered by the acceleration service"
  type        = string
  default     = "mainland_china"

  validation {
    condition     = contains(["mainland_china", "outside_mainland_china", "global"], var.service_area)
    error_message = "The service_area must be one of: mainland_china, outside_mainland_china, global."
  }
}

variable "origin_server" {
  description = "The origin server address (IP address or domain name)"
  type        = string
}

variable "origin_type" {
  description = "The origin server type"
  type        = string
  default     = "ipaddr"

  validation {
    condition     = contains(["ipaddr", "domain", "obs_bucket"], var.origin_type)
    error_message = "The origin_type must be one of: ipaddr, domain, obs_bucket."
  }
}

variable "http_port" {
  description = "The HTTP port of the origin server"
  type        = number
  default     = 80
}

variable "https_port" {
  description = "The HTTPS port of the origin server"
  type        = number
  default     = 443
}

variable "origin_protocol" {
  description = "The protocol used to retrieve data from the origin server"
  type        = string
  default     = "http"

  validation {
    condition     = contains(["http", "https", "follow"], var.origin_protocol)
    error_message = "The origin_protocol must be one of: http, https, follow."
  }
}

variable "ipv6_enable" {
  description = "Whether to enable IPv6"
  type        = bool
  default     = false
}

variable "range_based_retrieval_enabled" {
  description = "Whether to enable range-based retrieval"
  type        = bool
  default     = false
}

variable "domain_description" {
  description = "The description of the CDN domain"
  type        = string
  default     = ""
}

variable "https_enabled" {
  description = "Whether to enable HTTPS"
  type        = bool
  default     = false
}

variable "certificate_name" {
  description = "The name of the SSL certificate (required when https_enabled is true)"
  type        = string
  default     = ""
  nullable    = false
}

variable "certificate_source" {
  description = "The source of the SSL certificate (required when https_enabled is true)"
  type        = string
  default     = "0"
  nullable    = false

  validation {
    condition     = contains(["0", "2"], var.certificate_source)
    error_message = "The certificate_source must be one of: 0, 2."
  }
}

variable "certificate_body_path" {
  description = "The file path to the SSL certificate (required when https_enabled is true and using custom certificate)."
  type        = string
  default     = ""
  sensitive   = false
  nullable    = false
}

variable "private_key_path" {
  description = "The file path to the private key (required when https_enabled is true and using custom certificate)."
  type        = string
  default     = ""
  sensitive   = false
  nullable    = false
}

variable "http2_enabled" {
  description = "Whether to enable HTTP/2 (only valid when https_enabled is true)"
  type        = bool
  default     = false
}

variable "ocsp_stapling_status" {
  description = "The OCSP stapling status (only valid when https_enabled is true)"
  type        = string
  default     = "off"

  validation {
    condition     = contains(["on", "off"], var.ocsp_stapling_status)
    error_message = "The ocsp_stapling_status must be one of: on, off."
  }
}

variable "cache_rules" {
  description = "The cache rules configuration"
  type        = list(object({
    rule_type           = string
    content             = string
    ttl                 = number
    ttl_type            = string
    priority            = number
    url_parameter_type  = optional(string)
    url_parameter_value = optional(string)
  }))
  default     = []
}

variable "domain_tags" {
  description = "The tags of the CDN domain"
  type        = map(string)
  default     = {}
}
