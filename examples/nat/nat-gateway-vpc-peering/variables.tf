# Provider Configuration
variable "region_name" {
  description = "The name of the region to deploy resources"
  type        = string
}

variable "access_key" {
  description = "The access key of HuaweiCloud"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "The secret key of HuaweiCloud"
  type        = string
  sensitive   = true
}

# ECS Instance Configuration
variable "image_name" {
  description = "The name of the image to use for ECS instances"
  type        = string
  default     = "Ubuntu 22.04 server 64bit"
}

variable "image_visibility" {
  description = "The visibility of the image"
  type        = string
  default     = "public"
}

variable "image_most_recent" {
  description = "If more than one result is returned, use the most recent image"
  type        = bool
  default     = true
}

variable "flavor_performance_type" {
  description = "The performance type of the flavor"
  type        = string
  default     = "normal"
}

variable "flavor_cpu_core_count" {
  description = "The number of CPU cores for the flavor"
  type        = number
  default     = 2
}

variable "flavor_memory_size" {
  description = "The memory size in GB for the flavor"
  type        = number
  default     = 4
}

# Security Group Configuration
variable "secgroup_name" {
  description = "The name of the security group"
  type        = string
  default     = "secgroup-example"
}

variable "security_group_rules" {
  description = "List of security group rules"
  type        = list(object({
    direction        = string
    description      = string
    ethertype        = string
    protocol         = string
    ports            = optional(string)
    remote_ip_prefix = string
  }))
  default     = [
    {
      direction        = "ingress"
      description      = "Allow SSH access"
      ethertype        = "IPv4"
      protocol         = "tcp"
      ports            = "22"
      remote_ip_prefix = "0.0.0.0/0"
    },
    {
      direction        = "ingress"
      description      = "Allow HTTP and HTTPS access"
      ethertype        = "IPv4"
      protocol         = "tcp"
      ports            = "80,443"
      remote_ip_prefix = "0.0.0.0/0"
    },
    {
      direction        = "ingress"
      description      = "Allow ICMP access"
      ethertype        = "IPv4"
      protocol         = "icmp"
      ports            = null
      remote_ip_prefix = "0.0.0.0/0"
    }
  ]
}

# VPC Configuration
variable "vpcs" {
  description = "List of VPC configurations"
  type        = list(object({
    vpc_name             = string
    vpc_cidr             = string
    subnet_name          = string
    subnet_cidr          = string
    subnet_gateway_ip    = optional(string)
    instance_name        = string
    ecs_flavor_id        = optional(string)
    ecs_image_id         = optional(string)
    ecs_system_disk_type = optional(string)
    ecs_system_disk_size = optional(number)
    ecs_admin_password   = optional(string)
    ecs_instance_tags    = optional(map(string))
  }))
  default     = [
    {
      vpc_name             = "vpc-one"
      vpc_cidr             = "192.168.0.0/16"
      subnet_name          = "subnet-one"
      subnet_cidr          = "192.168.1.0/24"
      subnet_gateway_ip    = "192.168.1.1"
      instance_name        = "instance-one"
      ecs_flavor_id        = null
      ecs_image_id         = null
      ecs_system_disk_type = null
      ecs_system_disk_size = null
      ecs_admin_password   = null
      ecs_instance_tags    = null
    },
    {
      vpc_name             = "vpc-other"
      vpc_cidr             = "10.0.0.0/16"
      subnet_name          = "subnet-other"
      subnet_cidr          = "10.0.2.0/24"
      subnet_gateway_ip    = "10.0.2.1"
      instance_name        = "instance-other"
      ecs_flavor_id        = null
      ecs_image_id         = null
      ecs_system_disk_type = null
      ecs_system_disk_size = null
      ecs_admin_password   = null
      ecs_instance_tags    = null
    }
  ]
}

# ECS Instance Configuration
variable "ecs_admin_password" {
  description = "The administrator password of ECS instances (used as fallback when not specified in vpcs)"
  type        = string
  sensitive   = true
}

# VPC Peering Configuration
variable "peering_connection_name" {
  description = "The name of the VPC peering connection"
  type        = string
  default     = "peering-connection-example"
}

variable "route_type" {
  description = "The type of the route"
  type        = string
  default     = "peering"
}

# EIP and Bandwidth Configuration
variable "eips" {
  description = "List of EIP configurations"
  type        = list(object({
    name_suffix = string
  }))
  default     = [
    {
      name_suffix = "snat"
    },
    {
      name_suffix = "dnat"
    }
  ]
}

# NAT Gateway Configuration
variable "nat_gateway_name" {
  description = "The name of the NAT gateway"
  type        = string
  default     = "nat-gateway-example"
}

variable "bandwidth_size" {
  description = "The bandwidth size in Mbps"
  type        = number
  default     = 5
}

variable "bandwidth_charge_mode" {
  description = "The charge mode of the bandwidth"
  type        = string
  default     = "traffic"
}

variable "nat_gateway_spec" {
  description = "The specification of the NAT gateway"
  type        = string
  default     = "1"
}

# NAT SNAT Rule Configuration
variable "snat_rules" {
  description = "List of SNAT rule configurations"
  type        = list(object({
    eip_index    = number
    subnet_index = optional(number)
    source_type  = optional(number)
    vpc_index    = optional(number)
  }))
  default     = [
    {
      eip_index    = 0
      subnet_index = 0
      source_type  = null
      vpc_index    = null
    },
    {
      eip_index    = 0
      subnet_index = null
      source_type  = 1
      vpc_index    = 1
    }
  ]
}

# NAT DNAT Rule Configuration
variable "dnat_rules" {
  description = "List of DNAT rule configurations"
  type        = list(object({
    eip_index             = number
    instance_index        = number
    external_service_port = string
  }))
  default     = [
    {
      eip_index             = 1
      instance_index        = 0
      external_service_port = "8022"
    },
    {
      eip_index             = 1
      instance_index        = 1
      external_service_port = "8023"
    }
  ]
}

variable "dnat_protocol" {
  description = "The protocol for DNAT rules"
  type        = string
  default     = "tcp"
}

variable "dnat_internal_port" {
  description = "The internal service port for DNAT rules"
  type        = string
  default     = "22"
}
