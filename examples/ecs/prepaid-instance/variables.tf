# Variable definitions for authentication
variable "region_name" {
  description = "The region where the ECS instance is located"
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
variable "availability_zone" {
  description = "The availability zone to which the ECS instance belongs"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_flavor_id" {
  description = "The flavor ID of the ECS instance"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_performance_type" {
  description = "The performance type of the ECS instance flavor"
  type        = string
  default     = "normal"
}

variable "instance_cpu_core_count" {
  description = "The number of CPU cores of the ECS instance"
  type        = number
  default     = 2
}

variable "instance_memory_size" {
  description = "The memory size in GB of the ECS instance"
  type        = number
  default     = 4
}

variable "instance_image_id" {
  description = "The image ID of the ECS instance"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_image_visibility" {
  description = "The visibility of the ECS instance image"
  type        = string
  default     = "public"
}

variable "instance_image_os" {
  description = "The operating system of the ECS instance image"
  type        = string
  default     = "Ubuntu"
}

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

variable "security_group_ids" {
  description = "The list of security group IDs of the ECS instance"
  type        = list(string)
  default     = []
  nullable    = false
}

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
  default     = ""

  validation {
    condition     = !(length(var.security_group_ids) < 1 && var.security_group_name == "")
    error_message = "The security group name cannot be empty if the security group ID list is not set"
  }
}

variable "instance_name" {
  description = "The name of the ECS instance"
  type        = string
}

variable "instance_admin_password" {
  description = "The login password of the ECS instance"
  type        = string
  sensitive   = true
}

variable "instance_description" {
  description = "The description of the ECS instance"
  type        = string
  default     = null
}

variable "instance_system_disk_type" {
  description = "The type of the system disk of the ECS instance"
  type        = string
  default     = null
}

variable "instance_system_disk_size" {
  description = "The size of the system disk in GB of the ECS instance"
  type        = number
  default     = null
}

variable "instance_system_disk_iops" {
  description = "The IOPS of the system disk of the ECS instance"
  type        = number
  default     = null
}

variable "instance_system_disk_throughput" {
  description = "The throughput of the system disk of the ECS instance"
  type        = number
  default     = null
}

variable "instance_system_disk_dss_pool_id" {
  description = "The DSS pool ID of the system disk of the ECS instance"
  type        = string
  default     = null
}

variable "instance_metadata" {
  description = "The metadata key/value pairs of the ECS instance"
  type        = map(string)
  default     = null
}

variable "instance_tags" {
  description = "The key/value pairs to associate with the instance"
  type        = map(string)
  default     = null
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project to which the ECS instance belongs"
  type        = string
  default     = null
}

variable "instance_eip_id" {
  description = "The EIP ID to associate with the ECS instance"
  type        = string
  default     = null
}

variable "instance_eip_type" {
  description = "The EIP type to create and associate with the ECS instance"
  type        = string
  default     = null
}

variable "instance_bandwidth" {
  description = "The bandwidth configuration of the ECS instance"
  type        = object({
    share_type   = string
    id           = optional(string, null)
    size         = optional(number, null)
    charge_mode  = optional(string, null)
    extend_param = optional(map(string), null)
  })

  default = null
}

variable "period_unit" {
  description = "The unit of the period of the ECS instance"
  type        = string
  default     = "month"
}

variable "period" {
  description = "The charging period of the ECS instance"
  type        = number
  default     = 1
}

variable "auto_renew" {
  description = "Whether to enable auto-renewal of the ECS instance"
  type        = string
  default     = "false"
}
