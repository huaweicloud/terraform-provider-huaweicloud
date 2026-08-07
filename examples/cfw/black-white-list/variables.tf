# Variable definitions for authentication
variable "region_name" {
  description = "The region where the CFW firewall is located"
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
variable "fw_instance_id" {
  description = "The firewall instance ID"
  type        = string
  default     = ""
  nullable    = false
}

variable "blacklist_list_type" {
  description = "The list type of the blacklist rule. 4: blacklist, 5: whitelist"
  type        = number
  default     = 4
}

variable "whitelist_list_type" {
  description = "The list type of the whitelist rule. 4: blacklist, 5: whitelist"
  type        = number
  default     = 5
}

variable "blacklist_direction" {
  description = "The direction of the blacklist rule. 0: inbound, 1: outbound"
  type        = number
  default     = 0
}

variable "whitelist_direction" {
  description = "The direction of the whitelist rule. 0: inbound, 1: outbound"
  type        = number
  default     = 0
}

variable "blacklist_protocol" {
  description = "The protocol type of the blacklist rule. 6: TCP, 17: UDP, -1: any"
  type        = number
  default     = 6
}

variable "whitelist_protocol" {
  description = "The protocol type of the whitelist rule. 6: TCP, 17: UDP, -1: any"
  type        = number
  default     = 6
}

variable "blacklist_port" {
  description = "The destination port of the blacklist rule"
  type        = string
  default     = "22"
}

variable "whitelist_port" {
  description = "The destination port of the whitelist rule"
  type        = string
  default     = "80"
}

variable "blacklist_address_type" {
  description = "The IP address type of the blacklist rule. 0: IPv4, 1: IPv6"
  type        = number
  default     = 0
}

variable "whitelist_address_type" {
  description = "The IP address type of the whitelist rule. 0: IPv4, 1: IPv6"
  type        = number
  default     = 0
}

variable "blacklist_address" {
  description = "The IP address of the blacklist rule"
  type        = string
}

variable "whitelist_address" {
  description = "The IP address of the whitelist rule"
  type        = string
}
