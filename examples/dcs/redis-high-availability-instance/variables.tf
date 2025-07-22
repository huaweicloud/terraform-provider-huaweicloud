# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DCS Redis instance is located"
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
  description = "The gateway IP address of the subnet"
  type        = string
  default     = ""
  nullable    = false
}

variable "availability_zones" {
  description = "The availability zones to which the Redis instance belongs"
  type        = list(string)
  default     = []
  nullable    = false
}

variable "instance_flavor_id" {
  description = "The flavor ID of the Redis instance"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_cache_mode" {
  description = "The cache mode of the Redis instance"
  type        = string
  default     = "ha"
}

variable "instance_capacity" {
  description = "The capacity of the Redis instance (in GB)"
  type        = number
  default     = 4
}

variable "instance_engine_version" {
  description = "The engine version of the Redis instance"
  type        = string
  default     = "5.0"
}

variable "instance_name" {
  description = "The name of the Redis instance"
  type        = string
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project to which the Redis instance belongs"
  type        = string
  default     = null
}

variable "instance_password" {
  description = "The password for the Redis instance"
  type        = string
  sensitive   = true
  default     = null
}

variable "instance_backup_policy" {
  description = "The backup policy of the Redis instance"
  type = object({
    backup_type = optional(string, "auto")
    backup_at   = list(number)
    begin_at    = string
    save_days   = optional(number, null)
    period_type = optional(string, null)
  })
  default = null
}

variable "instance_whitelists" {
  description = "The whitelists of the Redis instance"
  type = list(object({
    group_name = string
    ip_address = list(string)
  }))
  default  = []
  nullable = false
}

variable "instance_parameters" {
  description = "The parameters of the Redis instance"
  type = list(object({
    id    = string
    name  = string
    value = string
  }))
  default  = []
  nullable = false
}

variable "instance_tags" {
  description = "The tags of the Redis instance"
  type        = map(string)
  default     = {}
}

variable "instance_rename_commands" {
  description = "The rename commands of the Redis instance"
  type        = map(string)
  default     = {}
}

variable "charging_mode" {
  description = "The charging mode of the Redis instance"
  type        = string
  default     = "postPaid"
}

variable "period_unit" {
  description = "The unit of the period"
  type        = string
  default     = null
}

variable "period" {
  description = "The period of the Redis instance"
  type        = number
  default     = null
}

variable "auto_renew" {
  description = "Whether auto renew is enabled"
  type        = string
  default     = "false"
}
