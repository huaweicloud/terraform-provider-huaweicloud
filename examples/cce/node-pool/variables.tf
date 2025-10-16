# Variable definitions for authentication
variable "region_name" {
  description = "The region where the CBH instance is located"
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
variable "vpc_id" {
  description = "The ID of the VPC"
  type        = string
  default     = ""
}

variable "subnet_id" {
  description = "The ID of the subnet"
  type        = string
  default     = ""
}

variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
  default     = ""

  validation {
    condition     = var.vpc_id != "" || var.vpc_name != ""
    error_message = "vpc_name must be provided if vpc_id is not provided."
  }
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "192.168.0.0/16"
}

variable "subnet_name" {
  description = "The name of the subnet"
  type        = string
  default     = ""

  validation {
    condition     = var.subnet_id == "" || var.subnet_name == ""
    error_message = "subnet_name must be provided if subnet_id is not provided."
  }
}

variable "subnet_cidr" {
  description = "The CIDR block of the subnet"
  type        = string
  default     = ""
}

variable "subnet_gateway_ip" {
  description = "The gateway IP of the subnet"
  type        = string
  default     = ""
}

variable "availability_zone" {
  description = "The availability zone where the CCE cluster will be created"
  type        = string
  default     = ""
}

variable "eip_address" {
  description = "The EIP address of the CCE cluster"
  type        = string
  default     = ""
}

variable "eip_type" {
  description = "The type of the EIP"
  type        = string
  default     = "5_bgp"
}

variable "bandwidth_name" {
  description = "The name of the bandwidth"
  type        = string
  default     = ""
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

variable "cluster_name" {
  description = "The name of the CCE cluster"
  type        = string
  default     = ""
}

variable "cluster_flavor_id" {
  description = "The flavor ID of the CCE cluster"
  type        = string
  default     = "cce.s1.small"
}

variable "cluster_version" {
  description = "The version of the CCE cluster"
  type        = string
  default     = null # Default to the latest version
  nullable    = true
}

variable "cluster_type" {
  description = "The type of the CCE cluster"
  type        = string
  default     = "VirtualMachine"
}

variable "container_network_type" {
  description = "The type of container network"
  type        = string
  default     = "overlay_l2"
}

variable "node_performance_type" {
  description = "The performance type of the node"
  type        = string
  default     = "general"
}

variable "node_cpu_core_count" {
  description = "The CPU core count of the node"
  type        = number
  default     = 4
}

variable "node_memory_size" {
  description = "The memory size of the node"
  type        = number
  default     = 8
}

variable "keypair_name" {
  description = "The name of the keypair"
  type        = string
}

variable "node_pool_type" {
  description = "The type of the node pool"
  type        = string
  default     = "vm"
}

variable "node_pool_name" {
  description = "The name of the node pool"
  type        = string
}

variable "node_pool_os_type" {
  description = "The OS type of the node pool"
  type        = string
  default     = "EulerOS 2.9"
}

variable "node_pool_initial_node_count" {
  description = "The initial node count of the node pool"
  type        = number
  default     = 2
}

variable "node_pool_min_node_count" {
  description = "The minimum node count of the node pool"
  type        = number
  default     = 2
}

variable "node_pool_max_node_count" {
  description = "The maximum node count of the node pool"
  type        = number
  default     = 10
}

variable "node_pool_scale_down_cooldown_time" {
  description = "The scale down cooldown time of the node pool"
  type        = number
  default     = 10
}

variable "node_pool_priority" {
  description = "The priority of the node pool"
  type        = number
  default     = 1
}

variable "node_pool_tags" {
  description = "The tags of the node pool"
  type        = map(string)
  default     = {}
}

variable "root_volume_type" {
  description = "The type of the root volume"
  type        = string
  default     = "SATA"
}

variable "root_volume_size" {
  description = "The size of the root volume"
  type        = number
  default     = 40
}

variable "data_volumes_configuration" {
  description = "The configuration of the data volumes"
  type        = list(object({
    volumetype     = string
    size           = number
    count          = number
    kms_key_id     = optional(string, null)
    extend_params  = optional(map(string), null)
    virtual_spaces = optional(list(object({
      name            = string
      size            = string
      lvm_lv_type     = optional(string, null)
      lvm_path        = optional(string, null)
      runtime_lv_type = optional(string, null)
    })), [])
  }))

  default = [
    {
      volumetype = "SSD"
      size       = 100
      count      = 1
    }
  ]

  validation {
    condition     = length(var.data_volumes_configuration) > 0
    error_message = "At least one data volume must be provided."
  }
}
