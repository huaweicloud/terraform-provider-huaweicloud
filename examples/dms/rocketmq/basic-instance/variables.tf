# Variable definitions for authentication
variable "region_name" {
  description = "The region where the RocketMQ instance is located"
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
variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "192.168.0.0/16"
}

variable "subnet_name" {
  description = "The name of the subnet"
  type        = string
}

variable "subnet_cidr" {
  description = "The CIDR block of the subnet"
  type        = string
  default     = ""
  nullable    = false
}

variable "subnet_gateway_ip" {
  description = "The gateway IP of the subnet"
  type        = string
  default     = ""
  nullable    = false
}

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "availability_zones" {
  description = "The list of the availability zones to which the RocketMQ instance belongs"
  type        = list(string)
  default     = []

  nullable = false
}

variable "instance_flavor_id" {
  description = "The flavor ID of the RocketMQ instance"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_flavor_type" {
  description = "The type of the RocketMQ instance flavor"
  type        = string
  default     = "cluster.small"
}

variable "availability_zones_count" {
  description = "The number of availability zones"
  type        = number
  default     = 1
}

variable "charging_mode" {
  description = "The charging mode of the RocketMQ instance"
  type        = string
  default     = "postPaid"
}

variable "instance_enable_publicip" {
  description = "Whether to enable the public IP of the RocketMQ instance"
  type        = bool
  default     = false
}

variable "instance_publicip_id" {
  description = "The ID of the public IP of the RocketMQ instance, multiple IDs separated by commas (,)"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_eips_count" {
  description = "The number of the public IP of the RocketMQ instance"
  type        = number
  default     = 0

  validation {
    condition     = !var.instance_enable_publicip || var.instance_eips_count > 0 || var.instance_publicip_id != ""
    error_message = "When 'instance_enable_publicip' is true, 'instance_eips_count' or 'instance_publicip_id' must be set"
  }
}

variable "eip_type" {
  description = "The type of the EIP"
  type        = string
  default     = "5_bgp"
}

variable "bandwidth_name" {
  description = "The name of the EIP bandwidth"
  type        = string
  default     = ""
  nullable    = false

  validation {
    condition     = var.instance_eips_count == 0 || var.bandwidth_name != ""
    error_message = "When 'instance_eips_count' is greater than 0, 'bandwidth_name' must be set"
  }
}

variable "bandwidth_size" {
  description = "The size of the bandwidth"
  type        = number
  default     = 5
}

variable "bandwidth_share_type" {
  description = "The share type of the bandwidth"
  type        = string
  default     = "PER"
}

variable "bandwidth_charge_mode" {
  description = "The charge mode of the bandwidth"
  type        = string
  default     = "traffic"
}

variable "instance_name" {
  description = "The name of the RocketMQ instance"
  type        = string
}

variable "instance_engine_version" {
  description = "The engine version of the RocketMQ instance"
  type        = string
  default     = ""
  nullable    = false

  validation {
    condition     = var.instance_flavor_id == "" || var.instance_engine_version != ""
    error_message = "When 'instance_flavor_id' is not empty, 'instance_engine_version' is required"
  }
}

variable "instance_storage_spec_code" {
  description = "The storage spec code of the RocketMQ instance"
  type        = string
  default     = ""
  nullable    = false

  validation {
    condition     = var.instance_flavor_id == "" || var.instance_storage_spec_code != ""
    error_message = "When 'instance_flavor_id' is not empty, 'instance_storage_spec_code' is required"
  }
}

variable "instance_storage_space" {
  description = "The storage space of the RocketMQ instance"
  type        = number
  default     = 800
}

variable "instance_broker_num" {
  description = "The number of the broker of the RocketMQ instance"
  type        = number
  default     = 1
}

variable "instance_description" {
  description = "The description of the RocketMQ instance"
  type        = string
  default     = ""
}

variable "instance_tags" {
  description = "The tags of the RocketMQ instance"
  type        = map(string)
  default     = {}
}

variable "enterprise_project_id" {
  description = "The enterprise project ID of the RocketMQ instance"
  type        = string
  default     = null
}

variable "instance_enable_acl" {
  description = "Whether to enable the ACL of the RocketMQ instance"
  type        = bool
  default     = false
}

variable "instance_tls_mode" {
  description = "The TLS mode of the RocketMQ instance"
  type        = string
  default     = "SSL"
}

variable "instance_configs" {
  description = "The configs of the RocketMQ instance"
  type        = list(map(string))
  default     = []
}

variable "period_unit" {
  description = "The period type of the RocketMQ instance"
  type        = string
  default     = null
}

variable "period" {
  description = "The period of the RocketMQ instance"
  type        = number
  default     = null
}

variable "auto_renew" {
  description = "Whether to enable the auto-renew of the RocketMQ instance"
  type        = string
  default     = null
}
