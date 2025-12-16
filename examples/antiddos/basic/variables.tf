# Variable definitions for authentication
variable "region_name" {
  description = "The region where the Anti-DDoS cloud domain is located"
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

# Variable definitions for module
variable "vpc_eip_publicip_type" {
  description = "The EIP type"
  type        = string
}

variable "vpc_eip_bandwidth_share_type" {
  description = "The bandwidth share type"
  type        = string
}

variable "vpc_eip_bandwidth_name" {
  description = "The bandwidth name"
  type        = string
  default     = null
}

variable "vpc_eip_bandwidth_size" {
  description = "The bandwidth size"
  type        = number
  default     = null
}

variable "vpc_eip_bandwidth_charge_mode" {
  description = "The bandwidth charge mode"
  type        = string
  default     = null
}

variable "smn_topic_name" {
  description = "The name of the topic to be created"
  type        = string
}

variable "smn_topic_display_name" {
  description = "The topic display name"
  type        = string
  default     = null
}

variable "smn_subscription_endpoint" {
  description = "The message endpoint"
  type        = string
}

variable "smn_subscription_protocol" {
  description = "The protocol of the message endpoint"
  type        = string
}

variable "smn_subscription_remark" {
  description = "The remark information"
  type        = string
  default     = null
}

variable "antiddos_traffic_threshold" {
  description = "The traffic cleaning threshold in Mbps"
  type        = number
}
