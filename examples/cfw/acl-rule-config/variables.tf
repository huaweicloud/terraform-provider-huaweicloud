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

variable "address_group_name" {
  description = "The name of the IP address group"
  type        = string
}

variable "address_group_description" {
  description = "The description of the IP address group"
  type        = string
  default     = ""
}

variable "service_group_name" {
  description = "The name of the service group"
  type        = string
}

variable "service_group_description" {
  description = "The description of the service group"
  type        = string
  default     = ""
}

variable "domain_name_group_name" {
  description = "The name of the domain name group"
  type        = string
}

variable "domain_name_group_type" {
  description = "The type of the domain name group"
  type        = number
  default     = 0
}

variable "domain_name_group_description" {
  description = "The description of the domain name group"
  type        = string
  default     = ""
}

variable "domain_name_group_domains" {
  description = "The list of domain names in the domain name group"
  type        = list(object({
    domain_name = string
    description = string
  }))
  default     = [
    {
      domain_name = "*.example.com"
      description = ""
    }
  ]
  nullable    = false
}

variable "acl_rule_ip_name" {
  description = "The name of the IP-based ACL rule"
  type        = string
}

variable "acl_rule_ip_description" {
  description = "The description of the IP-based ACL rule"
  type        = string
  default     = ""
}

variable "acl_rule_type" {
  description = "The ACL rule type. 0: Internet rule, 1: VPC rule, 2: NAT rule"
  type        = number
  default     = 0
}

variable "acl_rule_address_type" {
  description = "The ACL rule address type. 0: IPv4, 1: IPv6"
  type        = number
  default     = 0
}

variable "acl_rule_action_type" {
  description = "The ACL rule action type. 0: Allow, 1: Deny"
  type        = number
  default     = 0
}

variable "acl_rule_long_connect_enable" {
  description = "Whether to enable persistent connections. 0: disable, 1: enable"
  type        = number
  default     = 0
}

variable "acl_rule_status" {
  description = "The ACL rule status. 0: disable, 1: enable"
  type        = number
  default     = 1
}

variable "acl_rule_applications" {
  description = "The application list of the ACL rule"
  type        = list(string)
  default     = ["HTTPS"]
  nullable    = false
}

variable "acl_rule_source_addresses" {
  description = "The source IP address list of the ACL rule"
  type        = list(string)
  default     = ["1.1.1.1"]
  nullable    = false
}

variable "acl_rule_destination_addresses" {
  description = "The destination IP address list of the ACL rule"
  type        = list(string)
  default     = ["1.1.1.2"]
  nullable    = false
}

variable "acl_rule_custom_service_protocol" {
  description = "The protocol type of the custom service. 6: TCP, 17: UDP"
  type        = number
  default     = 6
}

variable "acl_rule_custom_service_source_port" {
  description = "The source port of the custom service"
  type        = string
  default     = "81"
}

variable "acl_rule_custom_service_dest_port" {
  description = "The destination port of the custom service"
  type        = string
  default     = "82"
}

variable "tags" {
  description = "The key/value pairs to associate with the resources"
  type        = map(string)
  default     = {
    key = "value"
  }
  nullable    = false
}

variable "acl_rule_domain_name" {
  description = "The name of the domain-based ACL rule"
  type        = string
}

variable "acl_rule_domain_description" {
  description = "The description of the domain-based ACL rule"
  type        = string
  default     = ""
}

variable "acl_rule_domain_direction" {
  description = "The direction of the domain-based ACL rule. 0: inbound, 1: outbound"
  type        = number
  default     = 1
}

variable "acl_rule_destination_domain_address_name" {
  description = "The destination domain address name"
  type        = string
  default     = "*.baidu.com"
}

variable "acl_rule_group_name" {
  description = "The name of the group-based ACL rule"
  type        = string
}

variable "acl_rule_group_description" {
  description = "The description of the group-based ACL rule"
  type        = string
  default     = ""
}

variable "acl_rule_service_group_protocol" {
  description = "The protocol type used by the service group"
  type        = number
  default     = 6
}
