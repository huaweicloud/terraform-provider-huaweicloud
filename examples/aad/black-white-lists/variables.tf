# Variable definitions for authentication
variable "region_name" {
  description = "The region where the AAD instance is located"
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

# Variable definitions for AAD instance
variable "instance_id" {
  description = "The ID of an existing AAD instance"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_config" {
  description = "The configuration for creating a new AAD instance"

  type = object({
    instance_name                  = string
    ip_type                        = number
    resource_region                = string
    instance_access_type           = string
    duration                       = number
    amount                         = number
    period_type                    = number
    service_bandwidth              = number
    basic_bandwidth                = optional(number, null)
    elastic_bandwidth              = optional(number, null)
    basic_qps                      = optional(number, null)
    forwarding_rule                = optional(number, null)
    protected_domain               = optional(number, null)
    elastic_service_bandwidth_type = optional(number, null)
    elastic_service_bandwidth      = optional(number, null)
    protection_package             = optional(string, null)
    enterprise_project_id          = optional(string, null)
  })

  default = null

  validation {
    condition     = var.instance_id != "" || var.instance_config != null
    error_message = "The 'instance_config' is required when 'instance_id' is not specified."
  }
}

# Variable definitions for black/white lists
variable "blacklist_ips" {
  description = "The list of IP addresses to add to the blacklist"
  type        = list(string)
}

variable "whitelist_ips" {
  description = "The list of IP addresses to add to the whitelist"
  type        = list(string)
}
