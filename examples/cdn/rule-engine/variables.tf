# Variable definitions for authentication
variable "region_name" {
  description = "The region where the CDN rule engine rule is located"
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
  description = "The accelerated domain name to which the rule engine rule belongs"
  type        = string
}

variable "rule_name" {
  description = "The name of the rule engine rule"
  type        = string

  validation {
    condition     = length(var.rule_name) >= 1 && length(var.rule_name) <= 50
    error_message = "The rule_name length must be between 1 and 50 characters."
  }
}

variable "rule_status" {
  description = "Whether to enable the rule engine rule"
  type        = string
  default     = "on"

  validation {
    condition     = contains(["on", "off"], var.rule_status)
    error_message = "The rule_status must be one of: on, off."
  }
}

variable "rule_priority" {
  description = "The priority of the rule engine rule"
  type        = number
  default     = 1

  validation {
    condition     = var.rule_priority >= 1 && var.rule_priority <= 100
    error_message = "The rule_priority must be between 1 and 100."
  }
}

variable "conditions" {
  description = "The trigger conditions of the rule engine rule, in JSON format"
  type        = string
  default     = ""
  nullable    = false
}

variable "cache_rule" {
  description = "The cache rule configuration"
  type        = object({
    ttl           = number
    ttl_unit      = string
    follow_origin = optional(string)
    force_cache   = optional(string)
  })
  default     = null
}

variable "access_control" {
  description = "The access control configuration"
  type        = object({
    type = string
  })
  default     = null
}

variable "http_response_headers" {
  description = "The list of HTTP response header configurations"
  type        = list(object({
    name   = string
    value  = string
    action = string
  }))
  default     = []
  nullable    = false
}

variable "browser_cache_rule" {
  description = "The browser cache rule configuration"
  type        = object({
    cache_type = string
  })
  default     = null
}

variable "request_url_rewrite" {
  description = "The access URL rewrite configuration"
  type        = object({
    execution_mode = string
    redirect_url   = string
  })
  default     = null
}

variable "flexible_origins" {
  description = "The list of flexible origin configurations"
  type        = list(object({
    sources_type      = string
    ip_or_domain      = string
    priority          = number
    weight            = number
    http_port         = optional(number)
    https_port        = optional(number)
    origin_protocol   = optional(string)
    host_name         = optional(string)
    obs_bucket_type   = optional(string)
    bucket_access_key = optional(string)
    bucket_secret_key = optional(string)
    bucket_region     = optional(string)
    bucket_name       = optional(string)
  }))
  default     = []
  nullable    = false
}

variable "origin_request_headers" {
  description = "The list of origin request header configurations"
  type        = list(object({
    action = string
    name   = string
    value  = optional(string)
  }))
  default     = []
  nullable    = true
}

variable "origin_request_url_rewrite" {
  description = "The origin request URL rewrite configuration"
  type        = object({
    rewrite_type = string
    target_url   = string
  })
  default     = null
}

variable "origin_range" {
  description = "The origin range configuration"
  type        = object({
    status = string
  })
  default     = null
}

variable "request_limit_rule" {
  description = "The request rate limit configuration"
  type        = object({
    limit_rate_after = number
    limit_rate_value = number
  })
  default     = null
}

variable "error_code_cache" {
  description = "The list of error code cache configurations"
  type        = list(object({
    code = number
    ttl  = number
  }))
  default     = []
  nullable    = false
}
