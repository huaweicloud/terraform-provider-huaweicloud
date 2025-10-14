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

variable "eni_ipv4_subnet_id" {
  description = "The ID of the ENI subnet"
  type        = string
  default     = ""
}

variable "eni_subnet_name" {
  description = "The name of the ENI subnet"
  type        = string
  default     = ""

  validation {
    condition     = var.eni_ipv4_subnet_id == "" || var.eni_subnet_name == ""
    error_message = "eni_subnet_name must be provided if eni_ipv4_subnet_id is not provided."
  }
}

variable "eni_subnet_cidr" {
  description = "The CIDR block of the ENI subnet"
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
  default     = "eni"
}

variable "cluster_description" {
  description = "The description of the CCE cluster"
  type        = string
  default     = ""
}

variable "cluster_tags" {
  description = "The tags of the CCE cluster"
  type        = map(string)
  default     = {}
}
