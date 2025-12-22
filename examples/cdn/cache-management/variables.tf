# Variable definitions for authentication
variable "region_name" {
  description = "The region where the CTS service is located"
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
variable "refresh_file_urls" {
  description = "The list of file URLs that need to be refreshed"
  type        = list(string)
  default     = []
  nullable    = false

  validation {
    condition     = length(var.refresh_file_urls) <= 1000
    error_message = "The refresh_file_urls list can contain up to 1000 URLs."
  }
}

variable "zh_url_encode" {
  description = "Whether to encode Chinese characters in URLs before cache refresh/preheat"
  type        = bool
  default     = false
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project to which the resource belongs"
  type        = string
  default     = "0"
}

variable "preheat_urls" {
  description = "The list of URLs that need to be preheated"
  type        = list(string)
  default     = []
  nullable    = false

  validation {
    condition     = length(var.preheat_urls) <= 1000
    error_message = "The preheat_urls list can contain up to 1000 URLs."
  }
}
