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

# Variable definitions for huaweicloud_cfw_firewall
variable "firewall_name" {
  description = "The CFW firewall name"
  type        = string
}

variable "firewall_flavor" {
  description = "The flavor version of the firewall"
  type        = string
  default     = "Professional"
}

variable "firewall_charging_mode" {
  description = "The charging mode of the firewall"
  type        = string
  default     = "postPaid"
}

variable "firewall_tags" {
  description = "The key/value pairs to associate with the resources"
  type        = map(string)
  default     = {
    key = "value"
    foo = "bar"
  }
}

# Variable definitions for huaweicloud_cfw_eip_auto_protection
variable "eip_auto_protection_status" {
  description = "Whether to enable auto-protection for EIPs. 1: enable, 0: disable"
  type        = number
  default     = 1
}

# Variable definitions for huaweicloud_cfw_eip_protection
variable "eip_protection_enabled" {
  description = "Whether to enable manual EIP protection for specific existing EIPs"
  type        = bool
  default     = false
}

variable "eip_protection_eip_ids" {
  description = "The list of existing EIPs to protect, each with id and public_ipv4"
  type        = list(object({
    id          = string
    public_ipv4 = string
  }))
  default     = []
  nullable    = false
}
