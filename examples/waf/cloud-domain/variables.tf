# Variable definitions for authentication
variable "region_name" {
  description = "The region where the WAF cloud domain is located"
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

# Variable definitions for WAF
variable "cloud_instance_resource_spec_code" {
  description = "The resource specification code of the WAF cloud instance"
  type        = string
}

variable "cloud_instance_bandwidth_expack_product" {
  description = "The configuration of the bandwidth extended packages"
  type        = list(object({
    resource_size = number
  }))
  default     = []
}

variable "cloud_instance_domain_expack_product" {
  description = "The configuration of the domain extended packages"
  type        = list(object({
    resource_size = number
  }))
  default     = []
}

variable "cloud_instance_rule_expack_product" {
  description = "The configuration of the rule extended packages"
  type        = list(object({
    resource_size = number
  }))
  default     = []
}

variable "cloud_instance_charging_mode" {
  description = "The charging mode of the WAF cloud instance"
  type        = string
}

variable "cloud_instance_period_unit" {
  description = "The unit of the subscription period"
  type        = string
}

variable "cloud_instance_period" {
  description = "The subscription period"
  type        = number
}

variable "cloud_instance_auto_renew" {
  description = "Whether to enable auto-renewal for the WAF cloud instance"
  type        = string
  default     = "false"
}

variable "enterprise_project_id" {
  description = "The enterprise project ID"
  type        = string
  default     = "0"
}

variable "cloud_domain" {
  description = "The domain name to be protected by WAF"
  type        = string
}

variable "cloud_certificate_id" {
  description = "The ID of the SSL certificate for the domain"
  type        = string
  default     = ""
}

variable "cloud_certificate_name" {
  description = "The name of the SSL certificate for the domain"
  type        = string
  default     = ""
}

variable "cloud_proxy" {
  description = "Whether to enable proxy for the WAF domain"
  type        = bool
  default     = false
}

variable "cloud_description" {
  description = "The description of the WAF domain"
  type        = string
  default     = ""
}

variable "cloud_website_name" {
  description = "The website name for the WAF domain"
  type        = string
  default     = ""
}

variable "cloud_protect_status" {
  description = "The protection status of the WAF domain"
  type        = number
  default     = 0
}

variable "cloud_forward_header_map" {
  description = "The field forwarding configuration"
  type        = map(string)
  default     = {}
}

variable "cloud_custom_page" {
  description = "Configuration for custom error pages"
  type        = list(object({
    http_return_code = string
    block_page_type  = string
    page_content     = string
  }))
  default     = []
}

variable "cloud_timeout_settings" {
  description = "Timeout settings for the WAF domain"
  type        = list(object({
    connection_timeout = number
    read_timeout       = number
    write_timeout      = number
  }))
  default     = []
}

variable "cloud_traffic_mark" {
  description = "Traffic marking configuration for the WAF domain"
  type        = list(object({
    ip_tags     = list(string)
    session_tag = string
    user_tag    = string
  }))
  default     = []
}

variable "cloud_server" {
  description = "List of origin server configurations for the WAF domain"
  type        = list(object({
    client_protocol = string
    server_protocol = string
    address         = string
    port            = number
    type            = string
    weight          = number
  }))
}
