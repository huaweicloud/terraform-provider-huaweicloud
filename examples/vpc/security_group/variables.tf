# Variable definitions for authentication
variable "region_name" {
  description = "The region where the security group is located"
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

# Variable definitions for resources/data sources
variable "security_group_name" {
  description = "The name of the security group."
  type        = string
}

variable "security_group_rule_configurations" {
  description = "The list of security group rule configurations. Each item is a map with keys: direction, ethertype, protocol, ports, remote_ip_prefix."
  type        = list(object({
    direction        = optional(string, "ingress")
    ethertype        = optional(string, "IPv4")
    protocol         = optional(string, null)
    ports            = optional(string, null)
    remote_ip_prefix = optional(string, "0.0.0.0/0")
  }))

  nullable = false
}
