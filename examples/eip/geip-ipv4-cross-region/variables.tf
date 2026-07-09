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
variable "instance_flavor_id" {
  description = "The ID of the ECS instance flavor"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_flavor_performance_type" {
  description = "The performance type of the ECS instance flavor"
  type        = string
  default     = "normal"
}

variable "instance_flavor_cpu_core_count" {
  description = "The CPU core count of the ECS instance flavor"
  type        = number
  default     = 2
}

variable "instance_flavor_memory_size" {
  description = "The memory size of the ECS instance flavor"
  type        = number
  default     = 4
}

variable "instance_image_id" {
  description = "The ID of the ECS instance image"
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
  description = "The OS of the ECS instance image"
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

variable "enterprise_project_id" {
  description = "The ID of the enterprise project"
  type        = string
  default     = ""
  nullable    = false
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

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "security_group_rule_configurations" {
  description = "The list of security group rule configurations"

  type = list(object({
    direction        = optional(string, "ingress")
    ethertype        = optional(string, "IPv4")
    protocol         = optional(string, null)
    ports            = optional(string, null)
    remote_ip_prefix = optional(string, "0.0.0.0/0")
  }))

  nullable = false
}

variable "instance_name" {
  description = "The name of the ECS instance"
  type        = string
}

variable "instance_administrator_password" {
  description = "The administrator password of the ECS instance"
  type        = string
  sensitive   = true
}

variable "internet_gateway_name" {
  description = "The name of the VPC internet gateway"
  type        = string
}

variable "internet_gateway_add_route" {
  description = "Whether to add route to the internet gateway"
  type        = bool
  default     = true
}

variable "global_eip_access_site" {
  description = "The access site used to filter the global EIP pool"
  type        = string
  default     = "cn-north-beijing"
  nullable    = false
}

variable "global_eip_ip_version" {
  description = "The IP version of the global EIP"
  type        = string
  default     = "4"
}

variable "internet_bandwidth_charge_mode" {
  description = "The charge mode of the global internet bandwidth"
  type        = string
  default     = "95peak_guar"
}

variable "internet_bandwidth_size" {
  description = "The size of the global internet bandwidth in Mbit/s"
  type        = number
  default     = 300
}

variable "internet_bandwidth_name" {
  description = "The name of the global internet bandwidth"
  type        = string
  default     = null
}

variable "internet_bandwidth_ingress_size" {
  description = "The ingress size of the global internet bandwidth in Mbit/s"
  type        = number
  default     = null
}

variable "internet_bandwidth_tags" {
  description = "The tags of the internet bandwidth"
  type        = map(string)
  default     = null
}

variable "global_eip_name" {
  description = "The name of the global EIP"
  type        = string
}

variable "global_eip_description" {
  description = "The description of the global EIP"
  type        = string
  default     = ""
}

variable "global_eip_tags" {
  description = "The tags of the global EIP"
  type        = map(string)
  default     = null
}

variable "gc_bandwidth_name" {
  description = "The name of the global connection bandwidth"
  type        = string
}

variable "gc_bandwidth_charge_mode" {
  description = "The charge mode of the global connection bandwidth"
  type        = string
  default     = "95"
}

variable "gc_bandwidth_size" {
  description = "The size of the global connection bandwidth in Mbit/s"
  type        = number
  default     = 100
}
