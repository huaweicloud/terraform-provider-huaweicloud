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
variable "availability_zone" {
  description = "The availability zone where the CCE cluster will be created"
  type        = string
  default     = ""
  nullable    = false
}

variable "vpc_id" {
  description = "The ID of the VPC"
  type        = string
  default     = ""
  nullable    = false
}

variable "subnet_id" {
  description = "The ID of the subnet"
  type        = string
  default     = ""
  nullable    = false
}

variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
  default     = ""
  nullable    = false

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
  nullable    = false

  validation {
    condition     = var.subnet_id != "" || var.subnet_name != ""
    error_message = "subnet_name must be provided if subnet_id is not provided."
  }
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

variable "eni_ipv4_subnet_id" {
  description = "The ID of the ENI subnet"
  type        = string
  default     = ""
  nullable    = false
}

variable "eni_subnet_name" {
  description = "The name of the ENI subnet"
  type        = string
  default     = ""
  nullable    = false

  validation {
    condition     = var.eni_ipv4_subnet_id == "" || var.eni_subnet_name == ""
    error_message = "eni_subnet_name must be provided if eni_ipv4_subnet_id is not provided."
  }
}

variable "eni_subnet_cidr" {
  description = "The CIDR block of the ENI subnet"
  type        = string
  default     = ""
  nullable    = false
}

variable "cluster_name" {
  description = "The name of the cluster"
  type        = string
  default     = ""
}

variable "cluster_flavor_id" {
  description = "The flavor ID of the cluster"
  type        = string
  default     = "cce.s1.small"
}

variable "cluster_version" {
  description = "The version of the cluster"
  type        = string
  default     = null # Default to the latest version
  nullable    = true
}

variable "cluster_type" {
  description = "The type of the cluster"
  type        = string
  default     = "VirtualMachine"
}

variable "container_network_type" {
  description = "The type of container network"
  type        = string
  default     = "eni"
}

variable "cluster_description" {
  description = "The description of the cluster"
  type        = string
  default     = ""
}

variable "cluster_tags" {
  description = "The tags of the cluster"
  type        = map(string)
  default     = {}
}

variable "node_flavor_id" {
  description = "The flavor ID of the node"
  type        = string
  default     = ""
  nullable    = false
}

variable "node_flavor_performance_type" {
  description = "The performance type of the node"
  type        = string
  default     = "normal"
}

variable "node_flavor_cpu_core_count" {
  description = "The CPU core count of the node"
  type        = number
  default     = 2
}

variable "node_flavor_memory_size" {
  description = "The memory size of the node"
  type        = number
  default     = 4
}

variable "node_partition" {
  description = "The name of the partition"
  type        = string
  default     = ""
  nullable    = false

  validation {
    condition     = var.node_partition != "" || var.partition_name != ""
    error_message = "node_partition must be provided if partition_name is not provided."
  }
}

variable "partition_name" {
  description = "The name of the partition"
  type        = string
  default     = ""
  nullable    = false
}

variable "partition_category" {
  description = "The category  of the partition"
  type        = string
  default     = "IES"
}

variable "partition_public_border_group" {
  description = "The group of the partition"
  type        = string
  default     = ""
  nullable    = false

  validation {
    condition     = var.node_partition == "" || var.partition_public_border_group == ""
    error_message = "partition_public_border_group must be provided if node_partition is not provided."
  }
}

variable "node_name" {
  description = "The name of the node"
  type        = string
}

variable "node_password" {
  description = "The root password to login node"
  type        = string
  sensitive   = true
}

variable "root_volume_type" {
  description = "The type of the root volume"
  type        = string
  default     = "SSD"
}

variable "root_volume_size" {
  description = "The size of the root volume"
  type        = number
  default     = 40
}

variable "data_volumes_configuration" {
  description = "The configuration of the data volumes"
  type        = list(object({
    volumetype = string
    size       = number
  }))

  default  = []
  nullable = false
}

variable "node_pool_name" {
  description = "The name of the node pool"
  type        = string
  default     = ""
  nullable    = false
}

variable "node_pool_os_type" {
  description = "The OS type of the node pool"
  type        = string
  default     = "EulerOS 2.9"
}

variable "node_pool_initial_node_count" {
  description = "The initial number of nodes in the node pool"
  type        = number
  default     = 1
}

variable "node_pool_password" {
  description = "The root password to login node"
  type        = string
  sensitive   = true
  default     = ""

  validation {
    condition     = var.node_pool_name != "" || var.node_pool_password != ""
    error_message = "node_pool_password must be provided if node_pool_name is not provided."
  }
}
