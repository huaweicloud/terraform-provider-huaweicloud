# Variable definitions for authentication
variable "region_name" {
  description = "The region where the VPC is located"
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

# Variable definitions for resources/data sources
variable "channel_name" {
  description = "The name of the custom event channel"
  type        = string
}

variable "source_name" {
  description = "The name of the custom event source"
  type        = string
}

variable "source_type" {
  description = "The type of the custom event source"
  type        = string
  default     = "APPLICATION"
}

variable "connection_name" {
  description = "The exact name of the connection to be queried"
  type        = string
  default     = "default"
}

variable "subscription_name" {
  description = "The name of the event subscription"
  type        = string
}

variable "sources_provider_type" {
  description = "The provider type of the event source"
  type        = string
  default     = "CUSTOM"
}

variable "source_op" {
  description = "The operation of the source"
  type        = string
  default     = "StringIn"
}

variable "targets_name" {
  description = "The name of the event target"
  type        = string
  default     = "HTTPS"
}

variable "targets_provider_type" {
  description = "The type of the event target"
  type        = string
  default     = "CUSTOM"
}

variable "transform" {
  description = "The transform configuration of the event target, in JSON format"
  type        = map(string)
  default     = {
    "type" : "ORIGINAL",
  }
}

variable "detail_name" {
  description = "The name(key) of the target detail configuration"
  type        = string
  default     = "detail"
}

variable "target_url" {
  description = "The target url of the event target"
  type        = string
}
